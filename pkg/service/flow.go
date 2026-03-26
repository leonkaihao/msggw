package service

import (
	"fmt"

	mbmodel "github.com/leonkaihao/msgbus/pkg/model"
	"github.com/leonkaihao/msggw/pkg/model"
	log "github.com/sirupsen/logrus"
)

type flowCB struct {
	name     string
	branches []*Branch
}

func NewFlowCB(name string, branches []*Branch) *flowCB {
	return &flowCB{name, branches}
}

func (fcb *flowCB) OnReceive(csmr mbmodel.Consumer, msg mbmodel.Messager) {

	msgCtx := model.NewMsgContext(msg.Topic(), msg.Metadata(), msg.Data())
	var (
		i  int
		br *Branch
	)
	for i, br = range fcb.branches {
		matched, err := fcb.filter(msgCtx, br.filters)
		if err != nil {
			log.Errorf("[flow %v]: branch %v failed to filter: %v", fcb.name, i, err)
			continue
		}
		if !matched {
			continue
		}
		log.Debugf("[flow %v]: filtering pass, topic %v, metadata %v", fcb.name, msg.Topic(), msg.Metadata())
		log.Debugf("[Flow branch %v] Transforming... ", i)
		msgCtx, err = fcb.transform(msgCtx, br.transforms)
		if err != nil {
			log.Errorf("[flow %v]: branch %v failed to transform: %v", fcb.name, i, err)
			break
		}
		log.Debugf("[Flow branch %v] Transforming... ", i)
		err = fcb.sendTo(msgCtx, br.sendTo)
		if err != nil {
			log.Errorf("[flow %v]: branch %v failed to send: %v", fcb.name, i, err)
		}
		break
	}
	if i == len(fcb.branches) {
		log.Errorf("[flow %v]: no matched branch for topic %v, metadata %v", fcb.name, msg.Topic(), msg.Metadata())
	}
}

func (fcb *flowCB) filter(ctx model.MsgContext, filters []*Filter) (bool, error) {
	for _, ft := range filters {
		lv := ft.exp[0]
		op := ft.exp[1]
		rv := ft.exp[2]
		if op.Type() != model.SYMTYPE_OPERATOR {
			return false, fmt.Errorf("[flow %v]: filter expression expect %v but got %v", fcb.name, model.SYMTYPE_OPERATOR, op.Type())
		}
		switch tp := op.Value().(type) {
		case model.Operator:
		default:
			return false, fmt.Errorf("[flow %v]: filter expression got unanted value type %v", fcb.name, tp)
		}
		oprt := op.Value().(model.Operator)
		matched, err := oprt.Do(ctx, lv, rv)
		if err != nil {
			return false, err
		}
		if !matched {
			return false, nil
		}
	}
	return true, nil
}

func (fcb *flowCB) transform(ctx model.MsgContext, trans []*Transform) (model.MsgContext, error) {
	topic := ctx.Topic()
	metadata := ctx.Metadata()
	data := ctx.Data()
	var err error
	for i, tr := range trans {
		k := tr.key
		v := tr.val
		switch k.Type() {
		case model.SYMTYPE_TOPIC:
			topic, err = v.Format(ctx)
			if err != nil {
				return nil, err
			}
			log.Debugf("--[transform %v] transform topic to %v", i, topic)
		case model.SYMTYPE_METADATA:
			key, err := k.Format(ctx)
			if err != nil {
				return nil, err
			}
			val, err := v.Format(ctx)
			if err != nil {
				return nil, err
			}
			metadata[key] = val
			log.Debugf("--[transform %v] transform metadata.%v to %v", i, key, val)
		}
	}

	log.Debugf("--[transform] all done.")
	ctx = model.NewMsgContext(topic, metadata, data)
	return ctx, nil
}

func (fcb *flowCB) sendTo(ctx model.MsgContext, st *SendTo) error {
	brk := st.dest
	producer := brk.Producer(fcb.name, ctx.Topic())

	pl := st.payload
	if pl != "" {
		err := producer.SetOption(mbmodel.PRD_OPTION_PAYLOAD, pl)
		if err != nil {
			return err
		}
	}
	err := producer.Fire(ctx.Data(), ctx.Metadata())
	log.Debugf("--[sendTo]: send topic %v to %v", ctx.Topic(), brk.URL())
	return err
}

func (fcb *flowCB) OnDestroy(inst interface{}) {
	switch inst := inst.(type) {
	case mbmodel.Service:
		log.Infof("[flow %v]: service was removed", fcb.name)
	case mbmodel.Consumer:
		log.Infof("consumer(%v) was removed from flow", inst.ID())
	}
}

func (fcb *flowCB) OnError(inst interface{}, err error) {
	switch inst := inst.(type) {
	case mbmodel.Service:
		log.Infof("[flow %v]: service was closed, error: %v", fcb.name, err)
	case mbmodel.Consumer:
		log.Errorf("[flow %v]: consumer(%v) error: %v", fcb.name, inst.ID(), err)
	}
}

// --------------------------

type Branch struct {
	name       string
	filters    []*Filter
	transforms []*Transform
	sendTo     *SendTo
}

type Filter struct {
	exp []model.Symbol
}

type Transform struct {
	key model.Symbol
	val model.Symbol
}

type SendTo struct {
	dest    mbmodel.Broker
	payload string
}
