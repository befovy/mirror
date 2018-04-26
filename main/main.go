package main

import (
	"fmt"
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/shurcooL/githubql"
	"github.com/baishuai/mirror"
)

func HandleOneIssue(number githubql.Int) {

}

func main() {
	fmt.Println("hello world")

	httpClient := mirror.NewCLient(mirror.Token)

	client := githubql.NewClient(httpClient)

	var err error

	variables := map[string]interface{}{
		"owner":  githubql.String("baishuai"),
		"name":   githubql.String("mirror"),
		"states": []githubql.IssueState{githubql.IssueStateClosed},
		"after":  (*githubql.String)(nil),
	}

	var issueInRepo mirror.IssueInRepo

	err = client.Query(context.Background(), &issueInRepo, variables)
	if err != nil {
		// Handle error.
		spew.Dump(err)
	}
	fmt.Println("issueRepositoryName:", issueInRepo.Repository.Name)
	fmt.Println("issueCount:", issueInRepo.Repository.Issues.TotalCount)

	fmt.Println(issueInRepo.Repository.CreatedAt.String())

	spew.Dump(issueInRepo)
}
