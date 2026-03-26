package operator

import (
	"github.com/leonkaihao/msggw/pkg/model"
)

type OperatorNOT struct {
}

func NewOperatorNOT() model.Operator {
	return &OperatorNOT{}

}

func (op *OperatorNOT) Name() string {
	return model.OPTYPE_NOT
}

func (op *OperatorNOT) Do(ctx model.MsgContext, lvalue, rvalue model.Symbol) (bool, error) {
	l, err := lvalue.Format(ctx)
	if err != nil {
		return false, err
	}
	r, err := rvalue.Format(ctx)
	if err != nil {
		return false, err
	}

	result := l != r
	return result, nil
}
