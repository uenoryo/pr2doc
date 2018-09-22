package main

import (
	"context"
	"fmt"
	"os"

	"github.com/uenoryo/pr2doc"
)

const (
	repoOwner = "owner"
	repoName  = "repo-name"
	token     = "github-access-token"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("pr2doc required commit hash.")
		return
	}

	ctx := context.Background()
	gs := pr2doc.NewGithubService(ctx, repoOwner, repoName, token)
	p2d := pr2doc.NewPr2Doc(gs)
	doc, err := p2d.Run(context.Background(), os.Args[1])
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(doc)
}
