package model

const (
	SYMTYPE_RAW      = "raw"
	SYMTYPE_METADATA = "metadata"
	SYMTYPE_FUNC     = "func"
	SYMTYPE_TOPIC    = "topic"
	SYMTYPE_PROP     = "prop"
	SYMTYPE_OPERATOR = "operator"
	SYMTYPE_KEYWORD  = "keyword"
	SYMTYPE_MIX      = "mix"
)

// operators type
const (
	OPTYPE_IS    = "IS"
	OPTYPE_NOT   = "NOT"
	OPTYPE_MATCH = "MATCH"
)

// keywords for configuration
const (
	KEYWORD_NULL = "NULL"
)

const (
	BROKERTYPE_NATS   = "nats"
	BROKERTYPE_MQTT3  = "mqtt3"
	BROKERTYPE_INPROC = "inproc"
)

const (
	FUNC_PUBADDR_TO_TOPIC = "PubAddrToTopic"
	FUNC_PUBPARA_TO_TOPIC = "PubParaToTopic"
	FUNC_TIMESTAMP_NOW    = "TimestampNow"
)
