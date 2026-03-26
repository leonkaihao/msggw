package symbol

import "github.com/leonkaihao/msggw/pkg/model"

type symbolBase struct {
	tp  string
	val interface{}
	chi []model.Symbol
}

func (sym *symbolBase) Type() string {
	return sym.tp
}

func (sym *symbolBase) Value() interface{} {
	return sym.val
}

func (sym *symbolBase) Children() []model.Symbol {
	return sym.chi
}
