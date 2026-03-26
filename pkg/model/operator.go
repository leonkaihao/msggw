package model

type Operator interface {
	Name() string
	Do(ctx MsgContext, left, right Symbol) (bool, error)
}
