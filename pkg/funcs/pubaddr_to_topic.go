package funcs

import (
	"fmt"
	"strings"

	"github.com/leonkaihao/msggw/pkg/model"
)

type PubAddrToTopicFunc struct {
}

func NewPubAddrToTopicFunc() model.Func {
	return &PubAddrToTopicFunc{}
}

func (pafunc *PubAddrToTopicFunc) Name() string {
	return model.FUNC_PUBADDR_TO_TOPIC
}

func (pafunc *PubAddrToTopicFunc) Execute(ctx model.MsgContext) (string, error) {
	pubaddr, ok := ctx.Metadata()["pub_addr"]
	if !ok {
		return "", fmt.Errorf("PubAddrToTopicFunc: no pub_addr for msg(%v)", ctx.Metadata())
	}
	params := strings.Split(strings.TrimSpace(pubaddr), ",")
	if len(params) != 2 {
		return "", fmt.Errorf("PubAddrToTopicFunc: pub_addr '%v' should have 2 params", pubaddr)
	}
	return fmt.Sprintf("org.%v.device.%v", params[0], params[1]), nil
}
