package mirror

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
)

func Sync() {
	configs, err := NewConfig("./mirror.yaml")
	if err != nil {
		spew.Dump(err)
		fmt.Println(err.Error())
	} else {
		for _, config := range configs {
			SyncWithConfig(config)
		}
	}
}

var handlers map[string]func(config SourceConfig)

//var sources map[string]func(config SourceConfig) Source

func RegsiterSource(name string, handle func(config SourceConfig)) {
	if handlers == nil {
		handlers = make(map[string]func(config SourceConfig))
	}
	handlers[name] = handle
}

func SyncWithConfig(config SourceConfig) {

	if handle, found := handlers[config.Type]; found {
		handle(config)
	} else {
		fmt.Printf("Unknown source type %s", config.Type)
	}

}
