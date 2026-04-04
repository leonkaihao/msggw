package service

import (
	"fmt"

	mbclient "github.com/leonkaihao/msgbus/pkg/client"
	mbcommon "github.com/leonkaihao/msgbus/pkg/common"
	mbmodel "github.com/leonkaihao/msgbus/pkg/model"
	"github.com/leonkaihao/msggw/pkg/config"
	"github.com/leonkaihao/msggw/pkg/model"
	"github.com/leonkaihao/msggw/pkg/parser"
)

type Service interface {
	Start() error
	Close()
}

type service struct {
	cfg        *config.Config
	clients    map[string]mbmodel.Client
	brokers    map[string]mbmodel.Broker  // name: broker
	mbServices map[string]mbmodel.Service // name: flow
	cctx       model.ConfigContext
}

func NewService(cfg *config.Config) (Service, error) {
	svc := &service{
		cfg:     cfg,
		cctx:    model.NewConfigContext(cfg.Props),
		clients: make(map[string]mbmodel.Client),
	}
	if err := svc.buildAll(); err != nil {
		return nil, err
	}
	return svc, nil
}

func (svc *service) Start() error {
	for _, s := range svc.mbServices {
		go func(s mbmodel.Service) {
			_ = s.Serve()
		}(s)
	}
	return nil
}

func (svc *service) Close() {
	for _, s := range svc.mbServices {
		_ = s.Close()
	}
	for _, cli := range svc.clients {
		_ = cli.Close()
	}
	svc.brokers = map[string]mbmodel.Broker{}
	svc.mbServices = map[string]mbmodel.Service{}
}

//----------------------------------------------------------

func (svc *service) buildAll() error {
	err := svc.buildClient()
	if err != nil {
		return err
	}

	brks, err := svc.buildBrokers(svc.cfg.Brokers)
	if err != nil {
		return err
	}
	svc.brokers = brks

	mbServices, err := svc.buildFlows(svc.cfg.Flows)
	if err != nil {
		return err
	}
	svc.mbServices = mbServices
	return nil
}

func (svc *service) buildClient() error {
	natscli, err := mbclient.NewBuilder(mbclient.CLI_NATS).Build()
	if err != nil {
		return err
	}
	mqtt3cli, err := mbclient.NewBuilder(mbclient.CLI_MQTT3).Build()
	if err != nil {
		return err
	}
	inproccli, err := mbclient.NewBuilder(mbclient.CLI_INPROC).Build()
	if err != nil {
		return err
	}
	svc.clients[model.BROKERTYPE_NATS] = natscli
	svc.clients[model.BROKERTYPE_MQTT3] = mqtt3cli
	svc.clients[model.BROKERTYPE_INPROC] = inproccli

	return nil
}

func (svc *service) getClient(brokerType string) (mbmodel.Client, error) {
	cli, ok := svc.clients[brokerType]
	if !ok {
		return nil, fmt.Errorf("getClient: no client of type %v", brokerType)
	}
	return cli, nil
}

func (svc *service) buildBrokers(cfgbrokers []*config.Broker) (map[string]mbmodel.Broker, error) {
	brks := make(map[string]mbmodel.Broker)
	for _, cfgbrk := range cfgbrokers {
		cli, err := svc.getClient(cfgbrk.Type)
		if err != nil {
			return nil, err
		}

		brk := cli.Broker(cfgbrk.Url)
		brks[cfgbrk.Name] = brk
	}
	return brks, nil
}

func (svc *service) buildFlows(cfgflows []*config.Flow) (map[string]mbmodel.Service, error) {
	svcs := make(map[string]mbmodel.Service)
	for i, cfgflow := range cfgflows {
		if cfgflow.Name == "" {
			return nil, fmt.Errorf("buildFlows: the name of flow %v is empty", i)
		}
		mbsvc, err := svc.buildFlow(cfgflow)
		if err != nil {
			return nil, err
		}
		svcs[cfgflow.Name] = mbsvc
	}
	return svcs, nil
}

func (svc *service) buildFlow(cfgflow *config.Flow) (mbmodel.Service, error) {
	name := cfgflow.Name
	sourceBrk, ok := svc.brokers[cfgflow.Source]
	if !ok {
		return nil, fmt.Errorf("buildFlow: source broker name %v is not defined", cfgflow.Source)
	}

	csmrs := []mbmodel.Consumer{}
	for _, sub := range cfgflow.Subscribes {
		csmrs = append(csmrs, sourceBrk.Consumer(name, sub, "msggw"))
	}
	branches, err := svc.buildBranches(cfgflow.Branches)
	if err != nil {
		return nil, err
	}
	mbsvc := mbcommon.NewService(NewFlowCB(name, branches))
	mbsvc.AddConsumers(csmrs)
	return mbsvc, nil
}

func (svc *service) buildBranches(cfgBranches []*config.Branch) ([]*Branch, error) {
	branches := []*Branch{}
	for _, cfgBranch := range cfgBranches {
		branch, err := svc.buildBranch(cfgBranch)
		if err != nil {
			return nil, err
		}
		branches = append(branches, branch)
	}
	return branches, nil
}

func (svc *service) buildBranch(cfgBranch *config.Branch) (*Branch, error) {
	filters, err := svc.buildFilters(cfgBranch.Filters)
	if err != nil {
		return nil, err
	}
	trans, err := svc.buildTransforms(cfgBranch.Transforms)
	if err != nil {
		return nil, err
	}
	sendTo, err := svc.buildSendTo(cfgBranch.SendTo)
	if err != nil {
		return nil, err
	}
	branch := &Branch{
		name:       cfgBranch.Name,
		filters:    filters,
		transforms: trans,
		sendTo:     sendTo,
	}
	return branch, nil
}

func (svc *service) buildFilters(cfgFilters []string) ([]*Filter, error) {
	ps := parser.NewParser(svc.cctx)
	filters := []*Filter{}
	for _, cfgFilter := range cfgFilters {

		syms, err := ps.ParseExpression(cfgFilter)
		if err != nil {
			return nil, err
		}

		filter := new(Filter)
		filter.exp = syms
		filters = append(filters, filter)
	}
	return filters, nil
}

func (svc *service) buildTransforms(cfgTrans []map[string]string) ([]*Transform, error) {
	ps := parser.NewParser(svc.cctx)
	trans := []*Transform{}
	for _, cfgTran := range cfgTrans {
		for k, v := range cfgTran {
			symKey, err := ps.ParseValue(k)
			if err != nil {
				return nil, err
			}
			symVal, err := ps.ParseValue(v)
			if err != nil {
				return nil, err
			}
			if len(symKey) != 1 {
				return nil, fmt.Errorf("buildTransforms: %v is not a legal key", k)
			}
			if len(symVal) != 1 {
				return nil, fmt.Errorf("buildTransforms: %v is not a legal value", v)
			}
			switch symKey[0].Type() {
			case model.SYMTYPE_METADATA:
			case model.SYMTYPE_TOPIC:
			default:
				return nil, fmt.Errorf("buildTransforms: key %v should not be of %v type", k, symKey[0].Type())
			}
			tran := &Transform{symKey[0], symVal[0]}
			trans = append(trans, tran)
		}
	}
	return trans, nil
}

func (svc *service) buildSendTo(cfgSendTo *config.SendTo) (*SendTo, error) {
	broker, ok := svc.brokers[cfgSendTo.Dest]
	if !ok {
		return nil, fmt.Errorf("buildSendTo: cannot find broker %v", cfgSendTo.Dest)
	}
	if cfgSendTo.Payload != "" {
		_, err := mbmodel.FindPayloadByType(cfgSendTo.Payload)
		if err != nil {
			return nil, err
		}
	}
	sendto := &SendTo{broker, cfgSendTo.Payload}
	return sendto, nil
}
