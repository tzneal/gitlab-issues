# gitlab-issues

A simple tool to extract lists of gitlab issues in both text and CSV format.  It uses the
handy [xanzy/go-gitlab](https://github.com/xanzy/go-gitlab) library to extract
the issues from gitlab.

## Note

Requires go 1.7+

# Installation
```go install github.com/tzneal/gitlab-issues```

# Usage
```
    Usage of gitlab-issues:
      -csv
        	format output as CSV
      -label string
        	issues with any maching label will be returned, multiple labels can be separated by a comma
      -milestone string
        	project milestone
      -o string
        	specify filename to write output to instead of stdout
      -project string
        	project name
      -token string
        	gitlab token from User Settings->Account page
      -url string
        	base URL to connect to (default "https://gitlab.com/")
```

# Examples
- Extract Issues from a particular milestone

```gitlab-issues --url https://gitlab.xyz.com --token my-secret-token --project my/project --milestone "Release 1.0"```

- Export as CSV

```gitlab-issues --url https://gitlab.xyz.com --token my-secret-token --project my/project --csv -o issues.csv```

- Pull issues with particular labels

```gitlab-issues --url https://gitlab.xyz.com --token my-secret-token --project my/project --label Planned,Development```


# FAQ

- Where do I find the token to use for my gitlab instance?

  If you click Profile -> Settings in gitlab, and look at the 'Account' tab, it's listed as your 'Private token'.
