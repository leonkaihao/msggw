package operator

import (
	"fmt"

	"github.com/leonkaihao/msggw/pkg/model"
)

func FindOperator(name string) (model.Operator, error) {
	switch name {
	case model.OPTYPE_IS:
		return NewOperatorIS(), nil
	case model.OPTYPE_NOT:
		return NewOperatorNOT(), nil
	case model.OPTYPE_MATCH:
		return NewOperatorMatch(), nil
	default:
		return nil, fmt.Errorf("FindOperator: unknown op type %v", name)
	}
}
