package mirror

import (
	"os"
	"io"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Token  string
	Login  string
	Repo   string
	Output string
}

func NewConfig(filename string) (*Config, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return NewConfigWithReader(f)
}

func NewConfigWithReader(in io.Reader) (*Config, error) {
	decoder := yaml.NewDecoder(in)
	var config Config
	err := decoder.Decode(&config)
	return &config, err;
}
