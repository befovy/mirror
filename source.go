package mirror


import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"time"
)

// The post source config, parsed from config file.
// Contains source type and config for the source
//
// 文章内容的来源配置，从配置文件中解析得到内容。
// 包括来源类型和每一种来源的自定义配置。
//
type SourceConfig struct {
	// 类型
	Type string `yaml:"type"`

	// 具体的配置
	Config map[string]interface{} `yaml:"config"`
}

// Parse `SourceConfig` from file
func newConfig(filename string) ([]SourceConfig, error) {

	buffer, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	var sc []SourceConfig

	err = yaml.Unmarshal(buffer, &sc)
	return sc, err
}

// An abstract blog article
//
// 抽象的一篇文章（博文）
type Post interface {
	Title() string
	Tags() []string
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Content() string
}

// An abstract interface for A collection of `Post`. Iterator Design
//
// 博文集合的抽象接口， 迭代器模式
type Source interface {
	// return next post, or nil if no more
	//
	// 返回下一篇文章，如果没有更多内容返回 nil
	Next() Post

	// full filename (with path) to save the content
	//
	// 包含路径的文件名，用于保存博文
	FileName(post Post) string
}
