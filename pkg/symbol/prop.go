package symbol

import (
	"fmt"

	"github.com/leonkaihao/msggw/pkg/model"
)

type symbolProp struct {
	symbolBase
	propVal string
}

func NewSymbolProp(val string, cctx model.ConfigContext) (model.Symbol, error) {
	if val == "" {
		return nil, fmt.Errorf("prop symbol has no value")
	}
	propVal, ok := cctx.Props()[val]
	if !ok {
		return nil, fmt.Errorf("prop has no key %v", val)
	}
	return &symbolProp{symbolBase: symbolBase{tp: model.SYMTYPE_PROP, val: val}, propVal: propVal}, nil
}

func (psym *symbolProp) Format(ctx model.MsgContext) (string, error) {
	return psym.propVal, nil
}
