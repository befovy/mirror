package main

import (
	"fmt"
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/shurcooL/githubql"
	"github.com/baishuai/mirror"
)

// https://developer.github.com/v4/explorer/
var query struct {
	Viewer struct {
		Login     githubql.String
		CreatedAt githubql.DateTime
	}
}

var issueCount mirror.IssueSummary

func HandleOneIssue(number githubql.Int) {

}

func main() {
	fmt.Println("hello world")

	httpClient := mirror.NewCLient(mirror.Token)

	client := githubql.NewClient(httpClient)

	var err error
	err = client.Query(context.Background(), &query, nil)
	if err != nil {
		// Handle error.
		spew.Dump(err)
		return
	}
	fmt.Println("    Login:", query.Viewer.Login)
	fmt.Println("CreatedAt:", query.Viewer.CreatedAt)

	variables := map[string]interface{}{
		"owner":  githubql.String("baishuai"),
		"name":   githubql.String("mirror"),
		"states": []githubql.IssueState{githubql.IssueStateOpen},
	}

	err = client.Query(context.Background(), &issueCount, variables)
	if err != nil {
		// Handle error.
		spew.Dump(err)
	}
	fmt.Println("issueRepositoryName:", issueCount.Rspository.Name)
	fmt.Println("issueCount:", issueCount.Rspository.Issues.TotalCount)

	fmt.Println(issueCount.Rspository.CreatedAt.String())

	//spew.Dump(issueCount)
}
