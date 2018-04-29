package mirror

import (
	"bufio"
	"os"
	"fmt"
	"path"
	"strings"
	"unicode"
)

func HandleIsues(path string, issues []Issue) {
	for _, issue := range issues {
		if !issue.ViewerDidAuthor {
			continue
		}
		err := handleOne(path, issue)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func handleOne(filepath string, issue Issue) error {
	filename := fmt.Sprintf("%s_%d.md", issue.CreatedAt.Format("2006-Jan-_2"), issue.Number)

	if !issue.IsPage() {
		filepath = path.Join(filepath, "post")
	}
	filename = path.Join(filepath, filename)
	var f *os.File
	var err error
	if exist(filename) {
		f, err = os.OpenFile(filename, os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
	} else {
		f, err = os.Create(filename)
		if err != nil {
			return err
		}
	}

	defer f.Close()

	wf := bufio.NewWriter(f)
	wf.WriteString("---\n")
	wf.WriteString(fmt.Sprintf("title: \"%s\"\n", issue.Title))
	wf.WriteString(fmt.Sprintf("date: %s\n", issue.CreatedAt.String()))
	wf.WriteString("---\n")
	wf.WriteString(string(issue.Body))

	for _, comment := range issue.Comments.Nodes {
		if !comment.ViewerDidAuthor {
			continue
		}
		body := strings.TrimLeftFunc(string(comment.Body), unicode.IsSpace)
		if !strings.HasPrefix(body, "<!-") {
			continue
		}
		wf.WriteString(body)
	}

	wf.Flush()
	return nil
}
