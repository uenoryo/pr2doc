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
