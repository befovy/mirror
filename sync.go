package mirror

import (
	"context"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/shurcooL/githubql"
)

func Sync() {
	config, err := NewConfig("./mirror.yaml")
	if err != nil {
		spew.Dump(err)
		fmt.Println(err.Error())
	} else {
		SyncWithConfig(config)
	}
}

func SyncWithConfig(config *Config) {

	httpClient := NewClient(config.Token)

	client := githubql.NewClient(httpClient)

	variables := map[string]interface{}{
		"owner":  githubql.String(config.Login),
		"name":   githubql.String(config.Repo),
		"states": []githubql.IssueState{githubql.IssueStateClosed},
		"after":  (*githubql.String)(nil),
	}

	var issueInRepo IssueInRepo
	succeed, failed := 0, 0
	for true {
		err := client.Query(context.Background(), &issueInRepo, variables)

		if err != nil {
			fmt.Println(err.Error())
		}
		s, f := HandleIsues(config.Output, issueInRepo.Issues())
		succeed += s
		failed += f
		if !issueInRepo.HasNextPage() {
			break
		}
		variables["after"] = issueInRepo.NextCursor()
	}

	fmt.Printf("Mirror finished, with %d sueeccd and %d failed\n", succeed, failed)
}
