package mirror

import "github.com/shurcooL/githubql"

type PageInfo struct {
	HasPreviousPage githubql.Boolean
	HasNextPage     githubql.Boolean
	StartCursor     githubql.String
	EndCursor       githubql.String
}

type IssueSummary struct {
	Rspository struct {
		Name        githubql.String
		Description githubql.String
		CreatedAt   githubql.DateTime
		Issues struct {
			TotalCount githubql.Int

			PageInfo PageInfo
		} `graphql:"issues(first: 1, states: $states)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type Issue struct {
	CreateAt githubql.DateTime
}
