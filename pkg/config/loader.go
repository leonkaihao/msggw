package config

import (
	"encoding/json"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type Loader interface {
	Load(configPath string) (*Config, error)
}

type jsonLoader struct {
}

func NewJsonConfigLoader() Loader {
	return &jsonLoader{}
}

func (ld *jsonLoader) Load(configPath string) (*Config, error) {

	fi, err := os.Open(configPath)
	if err != nil {
		log.Errorf("Error: %s", err)
		return nil, err
	}
	defer fi.Close()
	buf, err := io.ReadAll(fi)
	if err != nil {
		log.Errorf("Error: %s", err)
		return nil, err
	}

	cfg := &Config{}
	err = json.Unmarshal(buf, cfg)
	if err != nil {
		log.Errorf("Error: %s", err)
		return nil, err
	}
	for _, broker := range cfg.Brokers {
		log.Infof("load broker url %v", broker.Url)
		if broker.Url[0] == '$' {
			log.Infof("Parsing broker url %v...", broker.Url)
			url, ok := os.LookupEnv(broker.Url[1:])
			if ok {
				log.Infof("broker url %v from %v...", url, broker.Url)
				broker.Url = url
			}
		}
	}
	return cfg, err
}
