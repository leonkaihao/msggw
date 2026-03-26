package funcs

import (
	"fmt"

	"github.com/leonkaihao/msggw/pkg/model"
)

func FindFunc(name string) (model.Func, error) {
	var fn model.Func
	switch name {
	case model.FUNC_PUBADDR_TO_TOPIC:
		fn = NewPubAddrToTopicFunc()
	case model.FUNC_PUBPARA_TO_TOPIC:
		fn = NewPubParaToTopicFunc()
	case model.FUNC_TIMESTAMP_NOW:
		fn = NewTimestampNow()
	default:
		return nil, fmt.Errorf("FindFunc: unknown func name %v", name)
	}
	return fn, nil
}
