package model

type MsgContext interface {
	Topic() string
	Metadata() map[string]string
	Data() []byte
}

type ConfigContext interface {
	Props() map[string]string
}

type msgContext struct {
	topic    string
	metadata map[string]string
	data     []byte
}

func NewMsgContext(topic string, metadata map[string]string, data []byte) MsgContext {
	return &msgContext{
		topic, metadata, data,
	}
}

func (mctx *msgContext) Topic() string {
	return mctx.topic
}
func (mctx *msgContext) Metadata() map[string]string {
	return mctx.metadata
}
func (mctx *msgContext) Data() []byte {
	return mctx.data
}

type configContext struct {
	props map[string]string
}

func NewConfigContext(props map[string]string) ConfigContext {
	return &configContext{props}
}

func (cctx *configContext) Props() map[string]string {
	return cctx.props
}
