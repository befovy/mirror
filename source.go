package mirror

import "time"

type Post interface {
	FileName() string
	Title() string
	Tags() []string
	CreatedAt() time.Time
	UpdatedAt() time.Time
	Draft() bool
	Content() string
}

// An abstract interface for `Post` source
// Iterator design
type Source interface {
	Next() Post
}
