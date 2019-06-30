package issues

import (
	"context"
	"errors"
	"fmt"
	"github.com/befovy/mirror"
	"github.com/shurcooL/githubql"
	"io"
	"net/http"
	"path"
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

func (ew *errWriter) writeString(buf string) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.w.WriteString(buf)
}

type issueSource struct {
	login  string
	repo   string
	output string
	prefix string
	client *githubql.Client

	posts   []IssueEdge
	after   githubql.String
	hasNext bool

	lastPostCur githubql.String
}

func (is *issueSource) nextPage(now *IssueEdge) *IssueEdge {
	if is == nil {
		return nil
	}
	if now != nil {
		is.lastPostCur = now.Cursor
	}
	if is.hasNext && (is.posts == nil || len(is.posts) == 0) {
		var after interface{}
		after = (*githubql.String)(nil)
		if len(is.after) > 0 {
			after = is.after
		}
		variables := map[string]interface{}{
			"owner":   githubql.String(is.login),
			"first":   githubql.Int(100),
			"name":    githubql.String(is.repo),
			"states":  []githubql.IssueState{githubql.IssueStateClosed, githubql.IssueStateOpen},
			"after":   after,
			"csFirst": githubql.Int(100),
			"csAfter": (*githubql.String)(nil),
		}

		var issueInRepo IssueInRepo

		err := is.client.Query(context.Background(), &issueInRepo, variables)

		if err != nil {
			fmt.Println(err.Error())
			return nil
		}

		is.hasNext = issueInRepo.hasNextPage()
		is.after = issueInRepo.nextCursor()
		is.posts = issueInRepo.issueEdges()
	}

	if is.posts != nil && len(is.posts) > 0 {
		issue := is.posts[0]
		is.posts = is.posts[1:]
		return &issue
	}
	return nil
}

func (is *issueSource) Next() mirror.Post {
	var issue *IssueEdge = nil
	for true {
		pre := issue
		issue = is.nextPage(pre)
		if issue == nil {
			break
		}

		if !issue.Node.Closed {
			continue
		}
		if !issue.Node.ViewerDidAuthor {
			continue
		}

		is.getAllComment4IssueEdge(issue)
		break
	}
	return issue
}

func (is *issueSource) FileName(post mirror.Post) string {
	ie, ok := post.(*IssueEdge)
	if !ok {
		return ""
	}
	fname := fmt.Sprintf("%s_%s_%d.md", is.prefix, ie.Node.CreatedAt.Format("20060102"), ie.Node.Number)
	return path.Join(is.output, fname)
}

func (is *issueSource) getAllComment4IssueEdge(ie *IssueEdge) {
	if ie == nil || ie.Node.Comments.PageInfo.HasNextPage == false {
		return
	}

	for {
		if !ie.Node.Comments.PageInfo.HasNextPage {
			break
		}

		var after interface{}
		after = (*githubql.String)(nil)
		if len(is.lastPostCur) > 0 {
			after = is.lastPostCur
		}

		variables := map[string]interface{}{
			"owner":   githubql.String(is.login),
			"first":   githubql.Int(1),
			"name":    githubql.String(is.repo),
			"states":  []githubql.IssueState{githubql.IssueStateClosed, githubql.IssueStateOpen},
			"after":   after,
			"csFirst": githubql.Int(100),
			"csAfter": ie.Node.Comments.PageInfo.EndCursor,
		}

		var issueInRepo IssueInRepo
		err := is.client.Query(context.Background(), &issueInRepo, variables)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		if len(issueInRepo.Repository.Issues.Edges) != 1 {
			return
		}

		issueEdge := issueInRepo.Repository.Issues.Edges[0]

		ie.Node.Comments.PageInfo.HasNextPage = issueEdge.Node.Comments.PageInfo.HasNextPage
		ie.Node.Comments.PageInfo.EndCursor = issueEdge.Node.Comments.PageInfo.EndCursor
		ie.Node.Comments.Nodes = append(ie.Node.Comments.Nodes, issueEdge.Node.Comments.Nodes...)

	}
}

func issueHandler(config mirror.SourceConfig) mirror.Source {

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

	prefix, ok := data["prefix"].(string)
	if !ok || len(prefix) == 0 {
		return nil
	}

	client := githubql.NewClient(httpClient)

	return &issueSource{
		login:   login,
		repo:    repo,
		output:  output,
		client:  client,
		prefix:  prefix,
		hasNext: true,
	}
}

func (ie *IssueEdge) Title() string {
	return string(ie.Node.Title)
}

func (ie *IssueEdge) Tags() []string {
	return nil
}

func (ie *IssueEdge) CreatedAt() time.Time {
	return ie.Node.CreatedAt.Time
}

func (ie *IssueEdge) UpdatedAt() time.Time {
	return ie.Node.UpdatedAt.Time
}

func (ie *IssueEdge) Content() string {
	sb := new(strings.Builder)

	ew := &errWriter{w: sb}

	ew.writeString(string(ie.Node.Body))
	ew.writeString("\n")
	ew.writeString("\n")

	for _, comment := range ie.Node.Comments.Nodes {
		if !comment.ViewerDidAuthor {
			continue
		}

		body := strings.TrimLeftFunc(string(comment.Body), unicode.IsSpace)
		if !strings.HasPrefix(body, "<!-") {
			continue
		} else {
			ew.writeString(body)
			ew.writeString("\n")
			ew.writeString("\n")
		}

	}
	ew.writeString("\n")
	ew.writeString(fmt.Sprintf("> 本文通过 mirror 和 hugo 生成，原始地址 https://github.com%s", ie.Node.ResourcePath.String()))

	if ew.err != nil {
		fmt.Printf("Write err %s\n", ew.err.Error())
	}

	return sb.String()
}

func init() {
	fmt.Println("regist issue source")
	mirror.RegsiterSource("issues", issueHandler)
}
