package operator

import (
	"regexp"

	"github.com/leonkaihao/msggw/pkg/model"
)

type OperatorMatch struct {
}

func NewOperatorMatch() model.Operator {
	return &OperatorMatch{}

}

func (op *OperatorMatch) Name() string {
	return model.OPTYPE_MATCH
}

func (op *OperatorMatch) Do(ctx model.MsgContext, lvalue, rvalue model.Symbol) (bool, error) {
	l, err := lvalue.Format(ctx)
	if err != nil {
		return false, err
	}
	r, err := rvalue.Format(ctx)
	if err != nil {
		return false, err
	}
	exp := regexp.MustCompile(r)
	result := exp.MatchString(l)

	return result, nil
}
