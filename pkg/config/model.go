package config

type Config struct {
	LogLevel string            `json:"logLevel"`
	Props    map[string]string `json:"props"`
	Brokers  []*Broker         `json:"brokers"`
	Flows    []*Flow           `json:"flows"`
}

type Broker struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Url  string `json:"url"`
}

type Flow struct {
	Name       string    `json:"name"`
	Source     string    `json:"source"` // source broker
	Subscribes []string  `json:"subscribes"`
	Payload    string    `json:"payload"`
	Branches   []*Branch `json:"branches"`
}

type Branch struct {
	Name       string              `json:"name"`
	Filters    []string            `json:"filters"`
	Transforms []map[string]string `json:"transforms"`
	SendTo     *SendTo             `json:"sendTo"`
}

type SendTo struct {
	Dest    string `json:"dest"`
	Payload string `json:"payload"`
}
