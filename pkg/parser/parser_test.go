package parser

import (
	"fmt"
	"testing"

	"github.com/leonkaihao/msggw/pkg/model"
)

func Test_parser_ParseExpression(t *testing.T) {
	parser := NewParser(model.NewConfigContext(map[string]string{"name": "test", "deviceId": "dev1"}))

	//negtive tests
	_, err := parser.ParseExpression("")
	if err == nil {
		t.Errorf("expect num error but got nil")
	}

	_, err = parser.ParseExpression("A B")
	if err == nil {
		t.Errorf("expect num error but got nil")
	}

	_, err = parser.ParseExpression("A B C")
	if err == nil {
		t.Errorf("expect error for 1st param, should have '{}'")
	}

	_, err = parser.ParseExpression("{A} B C")
	if err == nil {
		t.Errorf("expect error for 1st param, invalid type in '{}' ")
	}

	_, err = parser.ParseExpression("{metadata} B C")
	if err == nil {
		t.Errorf("expect error for 1st param, no value for type metadata ")
	}

	_, err = parser.ParseExpression("{prop.name} IS test")
	if err == nil {
		t.Errorf("expect error for 1st param, no support for type prop ")
	}

	_, err = parser.ParseExpression("{func.PubAddrToTopic} is test")
	if err == nil {
		t.Errorf("expect error for 1st param, no support for type func ")
	}

	_, err = parser.ParseExpression("{metadata.key1} B C")
	if err == nil {
		t.Errorf("expect error for 2nd param, invalid operator type")
	}

	_, err = parser.ParseExpression("{metadata.key1} IS {}")
	if err == nil {
		t.Errorf("expect error for 3rd param, invalid value type")
	}

	_, err = parser.ParseExpression("{metadata.key1} IS {unk:va}")
	if err == nil {
		t.Errorf("expect error for 3rd param, invalid value type")
	}

	// positive tests
	syms, err := parser.ParseExpression("{metadata.key1} IS 1")
	if err != nil {
		t.Error(err)
	}
	if syms[0].Type() != model.SYMTYPE_METADATA || syms[0].Value().(string) != "key1" {
		t.Errorf("wrong symbol %v:%v", syms[0].Type(), syms[0].Value())
	}
	if syms[1].Type() != model.SYMTYPE_OPERATOR || syms[1].Value().(model.Operator).Name() != model.OPTYPE_IS {
		t.Errorf("wrong symbol %v:%v", syms[1].Type(), syms[1].Value())
	}
	if syms[2].Type() != model.SYMTYPE_RAW || syms[2].Value().(string) != "1" {
		t.Errorf("wrong symbol %v:%v", syms[2].Type(), syms[2].Value())
	}

	syms, err = parser.ParseExpression("{topic} MATCH {prop.deviceId}")
	if err != nil {
		t.Error(err)
		return
	}
	if syms[0].Type() != model.SYMTYPE_TOPIC {
		t.Errorf("wrong symbol %v:%v", syms[0].Type(), syms[0].Value())
	}
	if syms[1].Type() != model.SYMTYPE_OPERATOR || syms[1].Value().(model.Operator).Name() != model.OPTYPE_MATCH {
		t.Errorf("wrong symbol %v:%v", syms[1].Type(), syms[1].Value())
	}
	if syms[2].Type() != model.SYMTYPE_PROP || syms[2].Value().(string) != "deviceId" {
		t.Errorf("wrong symbol %v:%v", syms[2].Type(), syms[2].Value())
	}
	if str, err := syms[2].Format(nil); err != nil || str != "dev1" {
		t.Error("prop format fail")
	}

	syms, err = parser.ParseExpression("{metadata.key2} NOT {func.TimestampNow}")
	if err != nil {
		t.Error(err)
		return
	}
	if syms[1].Type() != model.SYMTYPE_OPERATOR || syms[1].Value().(model.Operator).Name() != model.OPTYPE_NOT {
		t.Errorf("wrong symbol %v:%v", syms[1].Type(), syms[1].Value())
	}
	if syms[2].Type() != model.SYMTYPE_FUNC || syms[2].Value().(model.Func).Name() != model.FUNC_TIMESTAMP_NOW {
		t.Errorf("wrong symbol %v:%v", syms[2].Type(), syms[2].Value())
	}
	if str, err := syms[2].Format(nil); err != nil || str == "" {
		t.Error("func TimestampNow format fail")
	}
}

func Test_parser_ParseValue(t *testing.T) {

	parser := NewParser(model.NewConfigContext(map[string]string{"name": "test", "deviceId": "dev1"}))

	//negative tests
	_, err := parser.ParseValue("{")
	if err == nil {
		t.Error("expect incompleted brace error")
	}

	_, err = parser.ParseValue("}")
	if err == nil {
		t.Error("expect wrong brace error")
	}

	_, err = parser.ParseValue("{}")
	if err == nil {
		t.Error("expect invalid type")
	}

	_, err = parser.ParseValue("{hello}")
	if err == nil {
		t.Error("expect invalid type")
	}

	_, err = parser.ParseValue("{metadata}")
	if err == nil {
		t.Error("expect no value error")
	}

	_, err = parser.ParseValue("{prop}")
	if err == nil {
		t.Error("expect no value error")
	}

	_, err = parser.ParseValue("{func}")
	if err == nil {
		t.Error("expect no value error")
	}

	syms, err := parser.ParseValue("{prop.none}")
	if err == nil {
		t.Errorf("expect wrong symbol %v:%v", syms[0].Type(), syms[0].Value())
	}

	syms, err = parser.ParseValue("{func.None}")
	if err == nil {
		t.Errorf("expect wrong symbol %v:%v", syms[0].Type(), syms[0].Value())
	}

	// positive tests
	syms, err = parser.ParseValue("")
	if err != nil {
		t.Error(err)
	}
	if len(syms) != 1 || syms[0].Type() != model.SYMTYPE_RAW {
		t.Error("wrong symbol")
	}

	syms, err = parser.ParseValue("hello world")
	if err != nil {
		t.Error(err)
	}
	if len(syms) != 1 || syms[0].Type() != model.SYMTYPE_RAW {
		t.Error("wrong symbol")
	}

	notKeys := []string{"MATCH", "IS", "NOT"}
	for _, k := range notKeys {
		syms, err = parser.ParseValue(k)
		if err != nil {
			t.Error(err)
		}
		if len(syms) != 1 || syms[0].Type() != model.SYMTYPE_RAW {
			t.Error("wrong symbol")
		}
	}

	syms, err = parser.ParseValue("{topic}")
	if err != nil || len(syms) != 1 || syms[0].Type() != model.SYMTYPE_TOPIC {
		t.Errorf("wrong symbol %v:%v", syms[0].Type(), syms[0].Value())
	}

	syms, err = parser.ParseValue("{metadata.name}")
	if err != nil || len(syms) != 1 || syms[0].Type() != model.SYMTYPE_METADATA || syms[0].Value().(string) != "name" {
		t.Errorf("wrong symbol %v:%v", syms[0].Type(), syms[0].Value())
	}

	syms, err = parser.ParseValue("{prop.deviceId}")
	if err != nil || len(syms) != 1 || syms[0].Type() != model.SYMTYPE_PROP || syms[0].Value().(string) != "deviceId" {
		t.Errorf("wrong symbol %v:%v", syms[0].Type(), syms[0].Value())
	}
	if str, err := syms[0].Format(nil); err != nil || str != "dev1" {
		t.Errorf("wrong symbol %v:%v format %v", syms[0].Type(), syms[0].Value(), str)
	}

	funcCases := []string{model.FUNC_PUBADDR_TO_TOPIC, model.FUNC_PUBPARA_TO_TOPIC, model.FUNC_TIMESTAMP_NOW}
	for _, fcase := range funcCases {
		val := fmt.Sprintf("{func.%v}", fcase)
		syms, err = parser.ParseValue(val)
		if err != nil || len(syms) != 1 || syms[0].Type() != model.SYMTYPE_FUNC || syms[0].Value().(model.Func).Name() != fcase {
			t.Errorf("wrong symbol %v:%v", syms[0].Type(), syms[0].Value())
		}
	}

	syms, err = parser.ParseValue("hello/{topic}world/{prop.deviceId}{metadata.name}!")
	if err != nil || len(syms) != 1 || syms[0].Type() != model.SYMTYPE_MIX || len(syms[0].Children()) != 6 {
		t.Errorf("wrong symbol %v:%v", syms[0].Type(), syms[0].Value())
	}
	symMix := syms[0].Children()
	if symMix[0].Type() != model.SYMTYPE_RAW ||
		symMix[1].Type() != model.SYMTYPE_TOPIC ||
		symMix[2].Type() != model.SYMTYPE_RAW ||
		symMix[3].Type() != model.SYMTYPE_PROP ||
		symMix[4].Type() != model.SYMTYPE_METADATA ||
		symMix[5].Type() != model.SYMTYPE_RAW {
		t.Errorf("mix symbol has wrong child")
	}
}
