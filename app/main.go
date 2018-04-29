package main

import (
	"fmt"
	"context"
	"github.com/davecgh/go-spew/spew"
	"github.com/shurcooL/githubql"
	"github.com/baishuai/mirror"
	//"github.com/gohugoio/hugo/commands"
)

func HandleOneIssue(number githubql.Int) {

}

func main() {
	fmt.Println("hello world")

	config, err := mirror.NewConfig("./config.yaml")
	if err != nil {
		spew.Dump(err)
		fmt.Println(err.Error())
	}

	httpClient := mirror.NewCLient(config.Token)

	client := githubql.NewClient(httpClient)

	variables := map[string]interface{}{
		"owner":  githubql.String(config.Login),
		"name":   githubql.String(config.Repo),
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

	mirror.HandleIsues(config.Output, issueInRepo.Issues())

	//commands.
}
