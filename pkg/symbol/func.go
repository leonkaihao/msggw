package symbol

import (
	"github.com/leonkaihao/msggw/pkg/funcs"
	"github.com/leonkaihao/msggw/pkg/model"
)

type symbolFunc struct {
	symbolBase
	fn model.Func
}

func (fsym *symbolFunc) Format(ctx model.MsgContext) (string, error) {
	return fsym.fn.Execute(ctx)
}

func NewSymbolFunc(name string) (model.Symbol, error) {

	fn, err := funcs.FindFunc(name)
	if err != nil {
		return nil, err
	}
	return &symbolFunc{symbolBase: symbolBase{tp: model.SYMTYPE_FUNC, val: fn}, fn: fn}, nil
}
