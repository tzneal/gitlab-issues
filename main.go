package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/xanzy/go-gitlab"
)

func main() {
	baseURL := flag.String("url", "https://gitlab.com/", "base URL to connect to")
	token := flag.String("token", "", "gitlab token from User Settings->Account page")
	project := flag.String("project", "", "project name")
	milestone := flag.String("milestone", "", "project milestone")
	label := flag.String("label", "", "issues with any maching label will be returned, multiple labels can be separated by a comma")
	exLabel := flag.String("xlabel", "", "issues with any matching label will be excluded, multiple labels can be separated with a comma")
	outputCSV := flag.Bool("csv", false, "format output as CSV")
	file := flag.String("o", "", "specify filename to write output to instead of stdout")
	flag.Parse()

	if *token == "" || *project == "" {
		flag.Usage()
		os.Exit(1)
	}
	git := gitlab.NewClient(nil, *token)

	// ensure the URI is terminated with a slash
	if !strings.HasSuffix(*baseURL, "/") {
		*baseURL = *baseURL + "/"
	}

	git.SetBaseURL(*baseURL + "api/v3")

	// setup list filtering options
	opts := &gitlab.ListProjectIssuesOptions{}
	opts.PerPage = 100
	if *milestone != "" {
		opts.Milestone = milestone
	}
	if *label != "" {
		opts.Labels = strings.Split(*label, ",")
	}

	allIssues := []*gitlab.Issue{}
	maxPages := 1
	for page := 0; page < maxPages; page++ {
		opts.Page = page + 1
		if page != 0 {
			fmt.Println("fetching issue page", opts.Page, "of", maxPages)
		}

		issues, rsp, err := git.Issues.ListProjectIssues(*project, opts)
		if err != nil {
			log.Printf("error retrieving issues: %s", err)
			os.Exit(1)
		}

		maxPages = rsp.LastPage
		allIssues = append(allIssues, issues...)
	}

	// filter out any issues excluded by label
	allIssues = filterOutLabels(allIssues, *exLabel)

	// sort issues by the project specific issue ID
	sort.Slice(allIssues, func(i int, j int) bool {
		return allIssues[i].IID < allIssues[j].IID
	})

	of := os.Stdout
	if *file != "" {
		var err error
		of, err = os.Create(*file)
		if err != nil {
			log.Fatalf("error creating %s: %s", *file, err)
		}
	}

	headers := []string{"ID", "State", "Assignee", "Labels", "Description"}
	fmt.Printf("found %d issues\n", len(allIssues))

	if *outputCSV {
		cw := csv.NewWriter(of)
		defer cw.Flush()

		cw.Write(headers)
		for _, issue := range allIssues {
			cw.Write(fieldsFrom(issue))
		}
	} else {
		tw := tabwriter.NewWriter(of, 4, 4, 2, ' ', 0)
		defer tw.Flush()

		fmt.Fprintf(tw, "%s\n", strings.Join(headers, "\t"))
		for _, issue := range allIssues {
			fmt.Fprintf(tw, "%s\n", strings.Join(fieldsFrom(issue), "\t"))
		}
	}
}

func fieldsFrom(issue *gitlab.Issue) []string {
	return []string{strconv.Itoa(issue.IID), issue.State, issue.Assignee.Name, strings.Join(issue.Labels, ","), issue.Title}
}

// filterOutLabels removes issues that have are marked with a label
func filterOutLabels(issues []*gitlab.Issue, exLabels string) []*gitlab.Issue {
	if exLabels == "" {
		return issues
	}
	labels := map[string]struct{}{}
	for _, l := range strings.Split(exLabels, ",") {
		labels[l] = struct{}{}
	}
	for i := 0; i < len(issues); {
		issue := issues[i]
		skip := false
		for _, l := range issue.Labels {
			if _, ok := labels[l]; ok {
				skip = true
			}
		}
		if skip {
			issues[i] = issues[len(issues)-1]
			issues = issues[0 : len(issues)-1]
		} else {
			i++
		}
	}
	return issues
}
