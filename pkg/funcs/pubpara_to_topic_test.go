package funcs

import (
	"testing"

	"github.com/leonkaihao/msggw/pkg/model"
)

func TestPubParaToTopicFunc(t *testing.T) {
	fn := NewPubParaToTopicFunc()
	ctx := model.NewMsgContext("topic1", map[string]string{"pub_para": "sensor/command"}, []byte{})
	str, err := fn.Execute(ctx)
	if err != nil {
		t.Error(err)
	}
	expect := "sensor.command"
	if str != expect {
		t.Errorf("expect %v but got %v", expect, str)
	}
}
