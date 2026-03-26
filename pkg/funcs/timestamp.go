package funcs

import (
	"fmt"
	"time"

	"github.com/leonkaihao/msggw/pkg/model"
)

type TimestampNow struct {
}

func NewTimestampNow() model.Func {
	return &TimestampNow{}
}

func (tfunc *TimestampNow) Name() string {
	return model.FUNC_TIMESTAMP_NOW
}

func (tfunc *TimestampNow) Execute(ctx model.MsgContext) (string, error) {
	return fmt.Sprintf("%v", time.Now().UnixMilli()), nil
}
