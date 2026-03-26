package symbol

import (
	"fmt"

	"github.com/leonkaihao/msggw/pkg/model"
	log "github.com/sirupsen/logrus"
)

type symbolHub struct {
	cctx model.ConfigContext
}

func NewSymbolHub(cctx model.ConfigContext) *symbolHub {
	return &symbolHub{cctx}
}

func (sh *symbolHub) NewSymbol(tp, val string, syms []model.Symbol) (model.Symbol, error) {
	log.Debugf("NewSymbol: %v %v", tp, val)
	switch tp {
	case model.SYMTYPE_RAW:
		return NewSymbolRaw(val)
	case model.SYMTYPE_METADATA:
		return NewSymbolMetadata(val)
	case model.SYMTYPE_FUNC:
		return NewSymbolFunc(val)
	case model.SYMTYPE_TOPIC:
		return NewSymbolTopic()
	case model.SYMTYPE_PROP:
		return NewSymbolProp(val, sh.cctx)
	case model.SYMTYPE_OPERATOR:
		return NewSymbolOperator(val)
	case model.SYMTYPE_KEYWORD:
		return NewSymbolKeyword(val)
	case model.SYMTYPE_MIX:
		return NewSymbolMix(val, syms)
	default:
		return nil, fmt.Errorf("NewSymbol: no symbol type %v", tp)
	}
}
