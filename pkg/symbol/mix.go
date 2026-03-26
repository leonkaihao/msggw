package symbol

import (
	"github.com/leonkaihao/msggw/pkg/model"
)

type symbolMix struct {
	symbolBase
}

func NewSymbolMix(val string, syms []model.Symbol) (model.Symbol, error) {
	return &symbolMix{
		symbolBase: symbolBase{
			tp:  model.SYMTYPE_MIX,
			val: val,
			chi: syms,
		}}, nil
}

// Format is of no use for operator
func (xsym *symbolMix) Format(ctx model.MsgContext) (string, error) {
	var join string
	for _, child := range xsym.chi {
		str, err := child.Format(ctx)
		if err != nil {
			return join, err
		}
		join += str
	}
	return join, nil
}
