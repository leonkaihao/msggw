package parser

import (
	"fmt"
	"strings"

	"github.com/leonkaihao/msggw/pkg/model"
	"github.com/leonkaihao/msggw/pkg/symbol"
	log "github.com/sirupsen/logrus"
)

type parser struct {
	cctx model.ConfigContext
}

func NewParser(cctx model.ConfigContext) model.Parser {
	return &parser{cctx}
}

func (ps *parser) ParseExpression(str string) ([]model.Symbol, error) {
	entries := strings.Fields(strings.TrimSpace(str))
	if len(entries) != 3 {
		return nil, fmt.Errorf("parser: the filter string `%v` expect 3 symbols, but got %v", str, len(entries))
	}
	lvar, err := ps.parseSymbol(entries[0])
	if err != nil {
		return nil, err
	}
	switch lvar.Type() {
	case model.SYMTYPE_METADATA:
	case model.SYMTYPE_TOPIC:
	default:
		return nil, fmt.Errorf("parser: parse `%v`, 1st filter symbol should not be of %v type", entries[0], lvar.Type())
	}

	op, err := ps.parseOperator(entries[1])
	if err != nil {
		return nil, err
	}
	switch op.Type() {
	case model.SYMTYPE_OPERATOR:
	default:
		return nil, fmt.Errorf("parser: parse `%v`, 2nd filter symbol should not be of %v type", entries[1], op.Type())
	}

	rval, err := ps.parseValue(entries[2])
	if err != nil {
		return nil, err
	}

	switch rval.Type() {
	case model.SYMTYPE_RAW:
	case model.SYMTYPE_TOPIC:
	case model.SYMTYPE_METADATA:
	case model.SYMTYPE_FUNC:
	case model.SYMTYPE_PROP:
	case model.SYMTYPE_KEYWORD:
	case model.SYMTYPE_MIX:
	default:
		return nil, fmt.Errorf("parser: parse `%v`, 3rd filter symbol should not be of %v type", entries[2], rval.Type())
	}
	return []model.Symbol{lvar, op, rval}, nil
}

func (ps *parser) ParseValue(str string) ([]model.Symbol, error) {
	sym, err := ps.parseValue(str)
	if err != nil {
		return nil, err
	}
	return []model.Symbol{sym}, nil
}

func (ps *parser) parseSymbol(str string) (model.Symbol, error) {
	log.Debugf("parseSymbol: %v", str)
	if str[0] != '{' || str[len(str)-1] != '}' {
		return nil, fmt.Errorf("parseSymbol: parse `%v`, expect complete brace", str)
	}
	fields := strings.Split(str[1:len(str)-1], ".")
	var tp, val string
	tp = fields[0]
	if len(fields) >= 2 {
		val = fields[1]
	}
	symhub := symbol.NewSymbolHub(ps.cctx)
	return symhub.NewSymbol(tp, val, nil)
}

func (ps *parser) parseOperator(str string) (model.Symbol, error) {
	log.Debugf("parseOperator: %v", str)
	symhub := symbol.NewSymbolHub(ps.cctx)
	return symhub.NewSymbol(model.SYMTYPE_OPERATOR, str, nil)
}

func (ps *parser) parseValue(str string) (model.Symbol, error) {
	log.Debugf("parseValue: %v", str)

	var (
		sym model.Symbol
		err error
	)
	symhub := symbol.NewSymbolHub(ps.cctx)
	sysms := []model.Symbol{}
	status := "out"
	var head, tail int
	head = 0
	for i, c := range str {
		switch status {
		case "out":
			if c == '{' {
				tail = i
				if head < tail {
					part := str[head:tail]
					switch part {
					case model.KEYWORD_NULL:
						sym, _ = symhub.NewSymbol(model.SYMTYPE_KEYWORD, part, nil)
					default:
						sym, err = symhub.NewSymbol(model.SYMTYPE_RAW, part, nil)
						if err != nil {
							return nil, err
						}
					}
					sysms = append(sysms, sym)
				}
				head = tail
				status = "in"
			} else if c == '}' {
				return nil, fmt.Errorf("parseValue: position %v of value '%v' has wrong brace", i, str)
			}
		case "in":
			if c == '{' {
				return nil, fmt.Errorf("parseValue: position %v of value '%v' has wrong brace", i, str)
			} else if c == '}' {
				tail = i + 1
				part := str[head:tail]
				sym, err := ps.parseSymbol(part)
				if err != nil {
					return nil, err
				}
				sysms = append(sysms, sym)
				head = tail
				status = "out"
			}
		}
	}
	if tail != len(str) {
		if status == "in" {
			return nil, fmt.Errorf("parseValue: '%v' has incompleted brace", str)
		}
		part := str[head:]
		switch part {
		case model.KEYWORD_NULL:
			sym, _ = symhub.NewSymbol(model.SYMTYPE_KEYWORD, part, nil)
		default:
			sym, err = symhub.NewSymbol(model.SYMTYPE_RAW, part, nil)
			if err != nil {
				return nil, err
			}
		}
		sysms = append(sysms, sym)
	}

	var result model.Symbol
	if len(sysms) > 1 {
		result, _ = symhub.NewSymbol(model.SYMTYPE_MIX, "", sysms)
	} else if len(sysms) == 1 {
		result = sysms[0]
	} else {
		result, _ = symhub.NewSymbol(model.SYMTYPE_RAW, "", nil)
	}
	return result, nil
}
