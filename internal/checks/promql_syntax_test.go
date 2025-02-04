package checks_test

import (
	"testing"

	"github.com/cloudflare/pint/internal/checks"
)

func newSyntaxCheck(_ string) checks.RuleChecker {
	return checks.NewSyntaxCheck()
}

func TestSyntaxCheck(t *testing.T) {
	testCases := []checkTest{
		{
			description: "valid recording rule",
			content:     "- record: foo\n  expr: sum(foo)\n",
			checker:     newSyntaxCheck,
			problems:    noProblems,
		},
		{
			description: "valid alerting rule",
			content:     "- alert: foo\n  expr: sum(foo)\n",
			checker:     newSyntaxCheck,
			problems:    noProblems,
		},
		{
			description: "no arguments for aggregate expression provided",
			content:     "- record: foo\n  expr: sum(\n",
			checker:     newSyntaxCheck,
			problems: func(uri string) []checks.Problem {
				return []checks.Problem{
					{
						Fragment: "sum(",
						Lines:    []int{2},
						Reporter: "promql/syntax",
						Text:     "syntax error: no arguments for aggregate expression provided",
						Severity: checks.Fatal,
					},
				}
			},
		},
		{
			description: "unclosed left parenthesis",
			content:     "- record: foo\n  expr: sum(foo) by(",
			checker:     newSyntaxCheck,
			problems: func(uri string) []checks.Problem {
				return []checks.Problem{
					{
						Fragment: "sum(foo) by(",
						Lines:    []int{2},
						Reporter: "promql/syntax",
						Text:     "syntax error: unclosed left parenthesis",
						Severity: checks.Fatal,
					},
				}
			},
		},
	}
	runTests(t, testCases)
}
