package symbol

import "github.com/leonkaihao/msggw/pkg/model"

type symbolTopic struct {
	symbolBase
}

func NewSymbolTopic() (model.Symbol, error) {
	return &symbolTopic{symbolBase: symbolBase{tp: model.SYMTYPE_TOPIC}}, nil
}
func (tsym *symbolTopic) Format(ctx model.MsgContext) (string, error) {
	return ctx.Topic(), nil
}
