package mirror

import "github.com/shurcooL/githubql"

type PageInfo struct {
	HasPreviousPage githubql.Boolean
	HasNextPage     githubql.Boolean
	StartCursor     githubql.String
	EndCursor       githubql.String
}

type IssueInRepo struct {
	Repository struct {
		Name        githubql.String
		Description githubql.String
		CreatedAt   githubql.DateTime
		Issues struct {
			TotalCount githubql.Int
			Nodes    []Issue
			PageInfo PageInfo
		} `graphql:"issues(first: 100, after: $after, states: $states)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

type IssueEdge struct {
	Cursor githubql.String
	Node   Issue
}

type Issue struct {
	CreatedAt githubql.DateTime
	Author struct {
		Login githubql.String
	}
	Body            githubql.String
	Closed          githubql.Boolean
	LastEditedAt    githubql.DateTime
	Locked          githubql.Boolean
	PublishedAt     githubql.DateTime
	State           githubql.IssueState
	Title           githubql.String
	UpdatedAt       githubql.DateTime
	ViewerDidAuthor githubql.Boolean

	Labels struct {
		TotalCount githubql.Int
		Nodes      []Label
	} `graphql:"labels(first: 100, after:null)"`
	Comments struct {
		TotalCount githubql.Int
	} `graphql:"comments(first: 100, after:null )"`
}

type Label struct {
	Color       githubql.String
	Description githubql.String
	IsDefault   githubql.Boolean
	Name        githubql.String
}

type IssueComment struct {
	Body            githubql.String
	ViewerDidAuthor githubql.Boolean
}
