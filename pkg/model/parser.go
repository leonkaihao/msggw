package model

type Parser interface {
	ParseExpression(str string) ([]Symbol, error)
	ParseValue(str string) ([]Symbol, error)
}
