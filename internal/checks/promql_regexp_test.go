package checks_test

import (
	"testing"

	"github.com/cloudflare/pint/internal/checks"
)

func newRegexpCheck(_ string) checks.RuleChecker {
	return checks.NewRegexpCheck()
}

func TestRegexpCheck(t *testing.T) {
	testCases := []checkTest{
		{
			description: "ignores rules with syntax errors",
			content:     "- alert: foo\n  expr: sum(foo) without(\n",
			checker:     newRegexpCheck,
			problems:    noProblems,
		},
		{
			description: "static match",
			content:     "- record: foo\n  expr: foo{job=\"bar\"}\n",
			checker:     newRegexpCheck,
			problems:    noProblems,
		},
		{
			description: "valid regexp",
			content:     "- record: foo\n  expr: foo{job=~\"bar.+\"}\n",
			checker:     newRegexpCheck,
			problems:    noProblems,
		},
		{
			description: "valid regexp",
			content:     "- record: foo\n  expr: foo{job!~\"(.*)\"}\n",
			checker:     newRegexpCheck,
			problems:    noProblems,
		},
		{
			description: "valid regexp",
			content:     "- record: foo\n  expr: foo{job=~\"a|b|c\"}\n",
			checker:     newRegexpCheck,
			problems:    noProblems,
		},
		{
			description: "unnecessary regexp",
			content:     "- record: foo\n  expr: foo{job=~\"bar\"}\n",
			checker:     newRegexpCheck,
			problems: func(uri string) []checks.Problem {
				return []checks.Problem{
					{
						Fragment: `foo{job=~"bar"}`,
						Lines:    []int{2},
						Reporter: checks.RegexpCheckName,
						Text:     `unnecessary regexp match on static string job=~"bar", use job="bar" instead`,
						Severity: checks.Bug,
					},
				}
			},
		},
		{
			description: "unnecessary negative regexp",
			content:     "- record: foo\n  expr: foo{job!~\"bar\"}\n",
			checker:     newRegexpCheck,
			problems: func(uri string) []checks.Problem {
				return []checks.Problem{
					{
						Fragment: `foo{job!~"bar"}`,
						Lines:    []int{2},
						Reporter: checks.RegexpCheckName,
						Text:     `unnecessary regexp match on static string job!~"bar", use job!="bar" instead`,
						Severity: checks.Bug,
					},
				}
			},
		},
		{
			description: "empty regexp",
			content:     "- record: foo\n  expr: foo{job=~\"\"}\n",
			checker:     newRegexpCheck,
			problems: func(uri string) []checks.Problem {
				return []checks.Problem{
					{
						Fragment: `foo{job=~""}`,
						Lines:    []int{2},
						Reporter: checks.RegexpCheckName,
						Text:     `unnecessary regexp match on static string job=~"", use job="" instead`,
						Severity: checks.Bug,
					},
				}
			},
		},
		{
			description: "unnecessary regexp anchor",
			content:     "- record: foo\n  expr: foo{job=~\"^.+$\"}\n",
			checker:     newRegexpCheck,
			problems: func(uri string) []checks.Problem {
				return []checks.Problem{
					{
						Fragment: `foo{job=~"^.+$"}`,
						Lines:    []int{2},
						Reporter: checks.RegexpCheckName,
						Text:     `prometheus regexp matchers are automatically fully anchored so match for job=~"^.+$" will result in job=~"^^.+$$", remove regexp anchors ^ and/or $`,
						Severity: checks.Bug,
					},
				}
			},
		},
		{
			description: "unnecessary regexp anchor",
			content:     "- record: foo\n  expr: foo{job=~\"(foo|^.+)$\"}\n",
			checker:     newRegexpCheck,
			problems: func(uri string) []checks.Problem {
				return []checks.Problem{
					{
						Fragment: `foo{job=~"(foo|^.+)$"}`,
						Lines:    []int{2},
						Reporter: checks.RegexpCheckName,
						Text:     `prometheus regexp matchers are automatically fully anchored so match for job=~"(foo|^.+)$" will result in job=~"^(foo|^.+)$$", remove regexp anchors ^ and/or $`,
						Severity: checks.Bug,
					},
				}
			},
		},
	}
	runTests(t, testCases)
}
