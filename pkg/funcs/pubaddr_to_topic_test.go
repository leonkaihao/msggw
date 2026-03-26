package funcs

import (
	"testing"

	"github.com/leonkaihao/msggw/pkg/model"
)

func TestPubAddrToTopicFunc(t *testing.T) {
	fn := NewPubAddrToTopicFunc()
	ctx := model.NewMsgContext("topic1", map[string]string{"pub_addr": "org1,dev1"}, []byte{})
	str, err := fn.Execute(ctx)
	if err != nil {
		t.Error(err)
	}
	expect := "org.org1.device.dev1"
	if str != expect {
		t.Errorf("expect %v but got %v", expect, str)
	}
}
