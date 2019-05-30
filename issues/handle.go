package issues

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
	"unicode"
)

type errWriter struct {
	err error
	w   *bufio.Writer
}

func (ew *errWriter) WriteString(buf string) {
	if ew.err != nil {
		return
	}
	_, ew.err = ew.w.WriteString(buf)
}

func HandleIsues(path string, issues []Issue) (int, int) {
	succeed, failed := 0, 0
	for _, issue := range issues {
		if !issue.ViewerDidAuthor {
			continue
		}
		err := handleOne(path, issue)
		if err != nil {
			fmt.Println(err.Error())
			failed += 1
		} else {
			succeed += 1
		}
	}
	return succeed, failed
}

func exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func handleOne(filepath string, issue Issue) error {
	filename := fmt.Sprintf("%s_%d.md", issue.CreatedAt.Format("20060102"), issue.Number)

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

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("got error: ", r)
		}
	}()

	wf := bufio.NewWriter(f)
	ew := &errWriter{w: wf}

	ew.WriteString("---\n")
	ew.WriteString(fmt.Sprintf("title: \"%s\"\n", issue.Title))
	ew.WriteString(fmt.Sprintf("date: %s\n", issue.CreatedAt.String()))
	ew.WriteString(fmt.Sprintf("lastmod: %s\n", issue.LastEditedAt.String()))
	if !issue.Closed {
		ew.WriteString("draft: true\n")
	}
	ew.WriteString("---\n")
	ew.WriteString(string(issue.Body))
	ew.WriteString("\n")
	ew.WriteString("\n")
	for _, comment := range issue.Comments.Nodes {
		if !comment.ViewerDidAuthor {
			continue
		}
		body := strings.TrimLeftFunc(string(comment.Body), unicode.IsSpace)
		if !strings.HasPrefix(body, "<!-") {
			continue
		} else {
			bodyLines := strings.SplitN(body, "\n", 2)
			if len(bodyLines) == 2 {
				body = bodyLines[1]
			}
		}
		ew.WriteString(body)
		ew.WriteString("\n")
		ew.WriteString("\n")
	}

	ew.WriteString(fmt.Sprintf("> 本文通过 mirror 和 hugo 生成，原始地址 https://github.com%s", issue.ResourcePath.String()))

	if ew.err != nil {
		fmt.Printf("Write err %s\n", ew.err.Error())
		err = ew.err
	} else {
		err = wf.Flush()
	}

	if err == nil {
		err = f.Close()
	} else {
		_ = f.Close()
	}

	return err
}
