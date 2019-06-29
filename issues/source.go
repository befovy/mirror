package issues

import (
	"context"
	"errors"
	"fmt"
	"github.com/befovy/mirror"
	"github.com/shurcooL/githubql"
	"io"
	"net/http"
	"strings"
	"time"
	"unicode"
)

type transport struct {
	token string
	base  http.RoundTripper
}

// RoundTrip authorizes and authenticates the request with an
// access token. If no token exists or token is expired,
// tries to refresh/fetch a new token.
func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	if len(t.token) == 0 {
		return nil, errors.New("transport's Token is empty")
	}
	req.Header.Set("Authorization", "Bearer"+" "+t.token)
	return t.base.RoundTrip(req)
}

type errWriter struct {
	err error
	w   io.StringWriter
}

func (ew *errWriter) WriteString(buf string) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.w.WriteString(buf)
}

type issues struct {
	login  string
	repo   string
	output string

	client *githubql.Client

	issues  []Issue
	after   githubql.String
	hasNext bool
}

func (is *issues) nextIsue() *Issue {
	if is == nil {
		return nil
	}

	if is.hasNext && (is.issues == nil || len(is.issues) == 0) {

		var after interface{}
		after = (*githubql.String)(nil)
		if len(is.after) > 0 {
			after = is.after
		}
		variables := map[string]interface{}{
			"owner":  githubql.String(is.login),
			"first":  githubql.Int(100),
			"name":   githubql.String(is.repo),
			"states": []githubql.IssueState{githubql.IssueStateClosed, githubql.IssueStateOpen},
			"after":  after,
		}

		var issueInRepo IssueInRepo

		err := is.client.Query(context.Background(), &issueInRepo, variables)

		if err != nil {
			return nil
		}

		is.hasNext = issueInRepo.HasNextPage()
		is.after = issueInRepo.NextCursor()
		is.issues = issueInRepo.Issues()
	}

	if is.issues != nil && len(is.issues) > 0 {
		issue := is.issues[0]
		is.issues = is.issues[1:]
		return &issue
	}
	return nil
}

func (is *issues) Next() mirror.Post {
	var issue *Issue
	for true {
		issue = is.nextIsue()
		if issue == nil {
			break
		}

		if !issue.Closed {
			continue
		}
		break
	}
	return issue
}

func issuehandler(config mirror.SourceConfig) mirror.Source {

	data := config.Config

	token, ok := data["token"].(string)
	if !ok {
		return nil
	}

	httpClient := &http.Client{
		Transport: &transport{
			token: token,
			base:  http.DefaultTransport,
		},
	}

	login, ok := data["login"].(string)
	if !ok || len(login) == 0 {
		return nil
	}
	repo, ok := data["repo"].(string)
	if !ok || len(repo) == 0 {
		return nil
	}
	output, ok := data["output"].(string)
	if !ok || len(output) == 0 {
		return nil
	}

	client := githubql.NewClient(httpClient)

	return &issues{
		login:   login,
		repo:    repo,
		output:  output,
		client:  client,
		hasNext: true,
	}
}

func (i *Issue) FileName() string {
	return fmt.Sprintf("issues_%s_%d.md", i.ICreatedAt.Format("20060102"), i.Number)
}

func (i *Issue) Title() string {
	return string(i.ITitle)
}

func (i *Issue) Tags() []string {
	return nil
}

func (i *Issue) CreatedAt() time.Time {
	return i.ICreatedAt.Time
}

func (i *Issue) UpdatedAt() time.Time {
	return i.IUpdatedAt.Time
}

func (i *Issue) Content() string {
	sb := new(strings.Builder)

	ew := &errWriter{w: sb}

	ew.WriteString(string(i.Body))
	ew.WriteString("\n")
	ew.WriteString("\n")

	for _, comment := range i.Comments.Nodes {
		if !comment.ViewerDidAuthor {
			continue
		}

		body := strings.TrimLeftFunc(string(comment.Body), unicode.IsSpace)
		if !strings.HasPrefix(body, "<!-") {
			continue
		} else {
			ew.WriteString(body)
			ew.WriteString("\n")
			ew.WriteString("\n")
		}
		ew.WriteString(fmt.Sprintf("> 本文通过 mirror 和 hugo 生成，原始地址 https://github.com%s", i.ResourcePath.String()))
	}

	if ew.err != nil {
		fmt.Printf("Write err %s\n", ew.err.Error())
	}

	return sb.String()
}

func init() {
	fmt.Println("regist issue source")
	mirror.RegsiterSource("issues", issuehandler)
}
