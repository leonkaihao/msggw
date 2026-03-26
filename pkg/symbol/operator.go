package symbol

import (
	"github.com/leonkaihao/msggw/pkg/model"
	"github.com/leonkaihao/msggw/pkg/operator"
)

type symbolOperator struct {
	symbolBase
}

func NewSymbolOperator(val string) (model.Symbol, error) {
	op, err := operator.FindOperator(val)
	if err != nil {
		return nil, err
	}
	return &symbolOperator{symbolBase: symbolBase{tp: model.SYMTYPE_OPERATOR, val: op}}, nil
}

// Format is of no use for operator
func (osym *symbolOperator) Format(ctx model.MsgContext) (string, error) {
	return osym.val.(string), nil
}
