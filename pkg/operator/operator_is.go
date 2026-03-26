package operator

import (
	"github.com/leonkaihao/msggw/pkg/model"
)

type OperatorIS struct {
}

func NewOperatorIS() model.Operator {
	return &OperatorIS{}
}

func (op *OperatorIS) Name() string {
	return model.OPTYPE_IS
}

func (op *OperatorIS) Do(ctx model.MsgContext, lvalue, rvalue model.Symbol) (bool, error) {
	l, err := lvalue.Format(ctx)
	if err != nil {
		return false, err
	}
	r, err := rvalue.Format(ctx)
	if err != nil {
		return false, err
	}
	result := l == r
	return result, nil
}
