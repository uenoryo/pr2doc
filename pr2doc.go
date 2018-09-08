package pr2doc

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"text/template"

	"github.com/pkg/errors"
)

const (
	descriptionIdentifier = "share"
)

type Pr2Doc struct {
	gs GithubService
}

func NewPr2Doc(gs GithubService) *Pr2Doc {
	return &Pr2Doc{
		gs: gs,
	}
}

type Doc struct {
	Title       string
	Description string
}

// Run is XXX
func (pr2doc *Pr2Doc) Run(ctx context.Context, commitHash string) (string, error) {
	docs, err := pr2doc.collectDoc(ctx, commitHash)
	if err != nil {
		return "", errors.Wrap(err, "error collectDoc")
	}
	tmpl := template.Must(template.ParseFiles("doc.tmpl"))
	var res bytes.Buffer
	if err := tmpl.Execute(&res, docs); err != nil {
		return "", errors.Wrap(err, "error execute template")
	}

	return res.String(), nil
}

func (pr2doc *Pr2Doc) collectDoc(ctx context.Context, commitHash string) ([]*Doc, error) {
	var docs []*Doc
	cmt, err := pr2doc.gs.GetCommit(ctx, commitHash)
	if err != nil {
		return docs, errors.Wrap(err, "error GetCommit")
	}

	prNum, err := pr2doc.findPRNumber(*cmt.Commit.Message)
	if err != nil || prNum == 0 {
		return docs, errors.Wrapf(err, "error find PR number, commit:%s", commitHash)
	}

	commits, err := pr2doc.gs.GetPullRequestCommits(ctx, prNum)
	if err != nil {
		return docs, errors.Wrap(err, "error GetPullRequestCommits")
	}

	prNums := make([]int, 0, len(commits))
	for _, cmt := range commits {
		num, err := pr2doc.findPRNumber(*cmt.Commit.Message)
		if err != nil {
			log.Printf("failed to find PR number, error:%s", err.Error())
			continue
		}
		if num == 0 {
			continue
		}
		prNums = append(prNums, num)
	}

	for _, num := range prNums {
		var doc Doc
		pr, err := pr2doc.gs.GetPullRequest(ctx, num)
		if err != nil {
			// TODO: log
			doc.Title = fmt.Sprintf("[ERROR] 取得失敗:#%d", num)
		} else {
			doc.Title = *pr.Title
			doc.Description = pr2doc.findDescription(*pr.Body, descriptionIdentifier)
		}
		docs = append(docs, &doc)
	}
	return docs, nil
}

func (pr2doc *Pr2Doc) findPRNumber(text string) (int, error) {
	re := regexp.MustCompile(" #(?P<number>[0-9]{1,5}) ")
	res := re.FindAllStringSubmatch(text, -1)

	if len(res) == 0 || len(res[0]) < 2 {
		return 0, nil
	}

	prNum, err := strconv.Atoi(res[0][1])
	if err != nil {
		return 0, errors.Wrap(err, "error convert PR number string to integer")
	}

	return prNum, nil
}

func (pr2doc *Pr2Doc) findDescription(body, identifier string) string {
	format := fmt.Sprintf("(?s)```%s\n(?P<description>.*)\n```$", identifier)
	re := regexp.MustCompile(format)
	res := re.FindAllStringSubmatch(body, -1)

	if len(res) == 0 || len(res[0]) < 2 {
		return ""
	}
	return res[0][1]
}
