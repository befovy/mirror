package mirror

import (
	"github.com/shurcooL/githubql"
)

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
		Issues      struct {
			TotalCount githubql.Int
			Nodes      []Issue
			PageInfo   PageInfo
		} `graphql:"issues(first: 100, after: $after, states: $states)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

func (i *IssueInRepo) HasNextPage() bool {
	return (bool)(i.Repository.Issues.PageInfo.HasNextPage)
}

func (i *IssueInRepo) NextCursor() githubql.String {
	return i.Repository.Issues.PageInfo.EndCursor
}

func (i *IssueInRepo) Issues() []Issue {
	return i.Repository.Issues.Nodes
}

type IssueEdge struct {
	Cursor githubql.String
	Node   Issue
}

type Issue struct {
	CreatedAt githubql.DateTime
	Author    struct {
		Login githubql.String
	}
	Number          githubql.Int
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
	} `graphql:"labels(first: 100, after: null)"`
	Comments struct {
		TotalCount githubql.Int
		Nodes      []IssueComment
	} `graphql:"comments(first: 100, after: null)"`
}

func (i *Issue) IsPage() bool {
	return i.hasMirrorLabel("PAGE")
}

func (i *Issue) IsPost() bool {
	return i.hasMirrorLabel("POST")
}

const mirrorLabelTag = "__mirror__"

func (i *Issue) hasMirrorLabel(tag string) bool {
	post := false
	for _, label := range i.Labels.Nodes {
		if string(label.Name) == tag && string(label.Description) == mirrorLabelTag {
			post = true
			break
		}
	}
	return post
}

func (i *Issue) RealLabels() []Label {
	labels := make([]Label, 0)
	for _, l := range i.Labels.Nodes {
		if string(l.Description) != mirrorLabelTag {
			labels = append(labels, l)
		}
	}
	return labels
}

type Label struct {
	//Color       githubql.String
	//IsDefault   githubql.Boolean
	Name        githubql.String
	Description githubql.String
}

type IssueComment struct {
	Body            githubql.String
	ViewerDidAuthor githubql.Boolean
}
