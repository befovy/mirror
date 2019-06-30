package mirror

import (
	"bufio"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"os"
	"reflect"
	"strings"
)

// Read Sync Config from file `./mirror.yaml`,
// Then run sync process to get content from the config source
func Sync() {
	SyncWithConfig("./mirror.yaml")
}


// Read Sync Config from `configfile`,
// Then run sync process to get content from the config source
func SyncWithConfig(configfile string) {
	configs, err := newConfig(configfile)
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


// A source handler return an instance of Source interface according the SourceConfig.
// The sync process find a handler by name firstly, then iterate the `Post` in the `Source`,
// output post to a specified location.
func RegsiterSource(name string, handle func(config SourceConfig) Source) {
	if handlers == nil {
		handlers = make(map[string]func(config SourceConfig) Source)
	}
	handlers[name] = handle
}

func syncWithConfig(config SourceConfig) {
	if handle, found := handlers[config.Type]; found {
		source := handle(config)
		err := syncSource(source)
		if err != nil {
			fmt.Println(err.Error())
		}
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

func (ew *errWriter) writeString(buf string) {
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
		fmt.Println(source.FileName(p))

		filename := source.FileName(p)
		var f *os.File
		var err error
		if exist(filename) {
			f, err = os.OpenFile(filename, os.O_WRONLY, 0644)
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}
		} else {
			f, err = os.Create(filename)
			if err != nil {
				fmt.Println(err.Error())
				return nil
			}
		}

		wf := bufio.NewWriter(f)
		ew := &errWriter{w: wf}

		ew.writeString("---\n")
		ew.writeString(fmt.Sprintf("title: %s\n", p.Title()))
		ew.writeString(fmt.Sprintf("date: %s\n", p.CreatedAt().String()))
		ew.writeString(fmt.Sprintf("lastmod: %s\n", p.UpdatedAt().String()))
		tags := p.Tags()
		if tags != nil && len(tags) > 0 {
			ew.writeString(fmt.Sprintf("tags: %s\n", strings.Join(tags, ",")))
		}
		ew.writeString("---\n")
		ew.writeString(string(p.Content()))
		ew.writeString("\n")
		ew.writeString("\n")

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
