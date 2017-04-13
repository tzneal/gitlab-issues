package main

import (
	"strings"
	"testing"

	gitlab "github.com/xanzy/go-gitlab"
)

func createIssueWithLabels(labels string) *gitlab.Issue {
	issue := &gitlab.Issue{}
	issue.Labels = strings.Split(labels, ",")
	return issue
}
func TestFilterLabel(t *testing.T) {

	tests := []struct {
		FilterOut   string
		ExpectCount int
	}{{"", 3},
		{"foo", 1},
		{"bar", 2},
		{"bar,baz", 1},
		{"foo bar", 2}}

	for _, tc := range tests {
		issues := []*gitlab.Issue{createIssueWithLabels("foo"), createIssueWithLabels("foo,bar"), createIssueWithLabels("foo bar,baz")}
		if got := len(filterOutLabels(issues, tc.FilterOut)); got != tc.ExpectCount {
			t.Errorf("filtered out %s, expected %d, got %d", tc.FilterOut, tc.ExpectCount, got)
		}
	}
}
