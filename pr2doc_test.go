package pr2doc

import "testing"

func Test_findDescription(t *testing.T) {
    type test struct {
        Title  string
        Input  string
        Output string
    }

    identifier := "test"
    tests := []test{
        {
            Title:  "success",
            Input:  "```test\ndescription\n```",
            Output: "description",
        },
        {
            Title:  "success (includes new line)",
            Input:  "```test\ndescription\ndescription\ndescription\n```",
            Output: "description\ndescription\ndescription",
        },
        {
            Title:  "error (mismatch identifier)",
            Input:  "```go\ntest\n```",
            Output: "",
        },
        {
            Title:  "error (missing new line)",
            Input:  "```test\ndescription```",
            Output: "",
        },
    }

    for _, test := range tests {
        t.Run(test.Title, func(t *testing.T) {
            p2d := &Pr2Doc{}
            if g, w := p2d.findDescription(test.Input, identifier), test.Output; g != w {
                t.Errorf("error find description %s, want %s", g, w)
            }
        })
    }
}

func Test_findPRNumber(t *testing.T) {
    type test struct {
        Title  string
        Input  string
        Output int
    }

    tests := []test{
        {
            Title:  "success",
            Input:  "Merge pull request #1234 from pr2doc/develop",
            Output: 1234,
        },
        {
            Title:  "success (too many number)",
            Input:  "Merge pull request #123 #567 #789 from pr2doc/develop",
            Output: 123,
        },
        {
            Title:  "error (missing PR number requires # )",
            Input:  "Merge pull request 12345 from pr2doc/develop",
            Output: 0,
        },
        {
            Title:  "error (no spacing)",
            Input:  "Merge pull request#987 from pr2doc/develop",
            Output: 0,
        },
        {
            Title:  "error (too large number)",
            Input:  "Merge pull request #456789 from pr2doc/develop",
            Output: 0,
        },
    }

    for _, test := range tests {
        t.Run(test.Title, func(t *testing.T) {
            p2d := &Pr2Doc{}
            prNum, err := p2d.findPRNumber(test.Input)
            if err != nil {
                t.Fatal("error findPRNumber", err.Error())
            }
            if g, w := prNum, test.Output; g != w {
                t.Errorf("error find PR number %d, want %d", g, w)
            }
        })
    }
}
