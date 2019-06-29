package mirror

import (
	"bufio"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os"
	"reflect"
	"strings"
)

func Sync() {
	configs, err := newConfig("./mirror.yaml")
	if err != nil {
		spew.Dump(err)
		fmt.Println(err.Error())
	} else {
		for _, config := range configs {
			syncWithConfig(config)
		}
	}
}

var handlers map[string]func(config SourceConfig) Source

//var sources map[string]func(config SourceConfig) Source

func RegsiterSource(name string, handle func(config SourceConfig) Source) {
	if handlers == nil {
		handlers = make(map[string]func(config SourceConfig) Source)
	}
	handlers[name] = handle
}

func syncWithConfig(config SourceConfig) {

	if handle, found := handlers[config.Type]; found {
		source := handle(config)
		syncSource(source)
	} else {
		fmt.Printf("Unknown source type %s", config.Type)
	}
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

type errWriter struct {
	err error
	w   *bufio.Writer
}

func (ew *errWriter) WriteString(buf string) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.w.WriteString(buf)
}

func syncSource(source Source) error {
	for true {
		p := source.Next()

		if p == nil || (reflect.ValueOf(p).Kind() == reflect.Ptr && reflect.ValueOf(p).IsNil()) {
			break
		}

		fmt.Println(p.Title())
		fmt.Println(p.FileName())

		filename := p.FileName()
		var f *os.File
		var err error
		if exist(filename) {
			f, err = os.OpenFile(filename, os.O_WRONLY, 0644)
			if err != nil {
				return nil
			}
		} else {
			f, err = os.Create(filename)
			if err != nil {
				return nil
			}
		}

		wf := bufio.NewWriter(f)
		ew := &errWriter{w: wf}

		ew.WriteString("---\n")
		ew.WriteString(fmt.Sprintf("title: %s\n", p.Title()))
		ew.WriteString(fmt.Sprintf("date: %s\n", p.CreatedAt().String()))
		ew.WriteString(fmt.Sprintf("lastmod: %s\n", p.UpdatedAt().String()))
		tags := p.Tags()
		if tags != nil && len(tags) > 0 {
			ew.WriteString(fmt.Sprintf("tags: %s\n", strings.Join(tags, ",")))
		}
		ew.WriteString("---\n")
		ew.WriteString(string(p.Content()))
		ew.WriteString("\n")
		ew.WriteString("\n")

		if ew.err != nil {
			fmt.Printf("Write err %s\n", ew.err.Error())
			err = ew.err
		} else {
			err = wf.Flush()
		}

		if err == nil {
			err = f.Close()
		} else {
			_ = f.Close()
		}

		if err != nil {
			return err
		}
	}
	return nil
}
