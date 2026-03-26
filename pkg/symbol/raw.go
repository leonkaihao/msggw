package symbol

import "github.com/leonkaihao/msggw/pkg/model"

type symbolRaw struct {
	symbolBase
}

func NewSymbolRaw(val string) (model.Symbol, error) {
	return &symbolRaw{symbolBase: symbolBase{tp: model.SYMTYPE_RAW, val: val}}, nil
}

func (rsym *symbolRaw) Format(ctx model.MsgContext) (string, error) {
	return rsym.val.(string), nil
}
