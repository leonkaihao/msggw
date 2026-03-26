package symbol

import (
	"fmt"

	"github.com/leonkaihao/msggw/pkg/model"
)

type symbolMeta struct {
	symbolBase
	key string
}

func NewSymbolMetadata(val string) (model.Symbol, error) {
	if val == "" {
		return nil, fmt.Errorf("metadata symbol has no value")
	}
	return &symbolMeta{symbolBase: symbolBase{tp: model.SYMTYPE_METADATA, val: val}, key: val}, nil
}
func (msym *symbolMeta) Format(ctx model.MsgContext) (string, error) {
	val := ctx.Metadata()[msym.key]
	return val, nil
}
