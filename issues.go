package mirror

import (
	"context"
	"fmt"
	"github.com/befovy/mirror/issues"
	"github.com/shurcooL/githubql"
)

func init() {
	RegsiterSource("issues", issueHanlder)
}

func issueHanlder(config SourceConfig) {

	data := config.Config

	token, ok := data["token"].(string)
	if !ok {
		return
	}
	httpClient := issues.NewClient(token)

	login := data["login"].(string)
	repo := data["repo"].(string)
	output := data["output"].(string)

	client := githubql.NewClient(httpClient)

	variables := map[string]interface{}{
		"owner":  githubql.String(login),
		"name":   githubql.String(repo),
		"states": []githubql.IssueState{githubql.IssueStateClosed, githubql.IssueStateOpen},
		"after":  (*githubql.String)(nil),
	}

	var issueInRepo issues.IssueInRepo
	succeed, failed := 0, 0
	for true {
		err := client.Query(context.Background(), &issueInRepo, variables)

		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("Get %d issues in page.\n", issueInRepo.Repository.Issues.TotalCount)
		fmt.Println(issueInRepo.Repository.CreatedAt)
		s, f := issues.HandleIsues(output, issueInRepo.Issues())
		succeed += s
		failed += f
		if !issueInRepo.HasNextPage() {
			break
		}
		variables["after"] = issueInRepo.NextCursor()
	}

	fmt.Printf("Mirror finished, with %d sueeccd and %d failed\n", succeed, failed)
}
