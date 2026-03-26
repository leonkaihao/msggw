package funcs

import (
	"fmt"
	"strings"

	"github.com/leonkaihao/msggw/pkg/model"
)

type PubParaToTopicFunc struct {
}

func NewPubParaToTopicFunc() model.Func {
	return &PubParaToTopicFunc{}
}

func (ppfunc *PubParaToTopicFunc) Name() string {
	return model.FUNC_PUBPARA_TO_TOPIC
}

func (ppfunc *PubParaToTopicFunc) Execute(ctx model.MsgContext) (string, error) {
	pubpara, ok := ctx.Metadata()["pub_para"]
	if !ok {
		return "", fmt.Errorf("PubParaToTopicFunc: no pub_para for msg(%v)", ctx.Metadata())
	}
	params := strings.Replace(strings.TrimSpace(pubpara), "/", ".", -1)
	return params, nil
}
