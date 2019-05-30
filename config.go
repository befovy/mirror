package mirror

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type SourceConfig struct {
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

func NewConfig(filename string) ([]SourceConfig, error) {

	buffer, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	var sc []SourceConfig

	err = yaml.Unmarshal(buffer, &sc)
	return sc, err
}
