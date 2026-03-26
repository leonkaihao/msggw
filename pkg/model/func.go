package model

type Func interface {
	Name() string
	Execute(ctx MsgContext) (string, error)
}
