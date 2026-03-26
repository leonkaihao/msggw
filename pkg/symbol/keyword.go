package symbol

import (
	"fmt"

	"github.com/leonkaihao/msggw/pkg/model"
)

type symbolKeyword struct {
	symbolBase
	kval string
}

func NewSymbolKeyword(val string) (model.Symbol, error) {
	var kval string
	switch val {
	case model.KEYWORD_NULL:
		kval = ""
	default:
		return nil, fmt.Errorf("NewSymbolKeyword: unknown keyword %v", val)
	}
	return &symbolKeyword{
		symbolBase: symbolBase{
			tp:  model.SYMTYPE_KEYWORD,
			val: val,
		}, kval: kval}, nil
}

func (ksym *symbolKeyword) Format(ctx model.MsgContext) (string, error) {
	return ksym.kval, nil
}
