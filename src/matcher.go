package main

import (
	"strings"
)

type MatchResult struct {
	Platform string
	Owner    string
	Repo     string
	Extra    string
}

type rule struct {
	Prefix       string
	TrimPrefix   bool
	MinParts     int
	Platform     string
	HasSpaces    bool
	AllowExtra   bool
	ExtraIsRest  bool
}

var rules = []rule{
	{"github.com/", true, 2, "github", true, true, false},
	{"raw.githubusercontent.com/", true, 2, "github_raw", false, true, false},
	{"gist.githubusercontent.com/", true, 1, "gist", false, false, false},
	{"api.github.com/repos/", true, 2, "github_api", false, false, false},
	{"huggingface.co/", true, 2, "huggingface", true, false, false},
	{"cdn-lfs.hf.co/", true, 2, "cdn-lfs", true, true, true},
	{"download.docker.com/", true, 1, "docker_download", false, false, false},
	{"github.githubassets.com/", true, 1, "githubassets", false, false, false},
	{"opengraph.githubassets.com/", true, 1, "githubassets", false, false, false},
}

func trimGitSuffix(s string) string {
	return strings.TrimSuffix(s, ".git")
}

// MatchURL 函数
func MatchURL(raw string) (res MatchResult, ok bool) {
	raw = strings.TrimPrefix(raw, "https://")
	raw = strings.TrimPrefix(raw, "http://")

	for _, r := range rules {
		if strings.HasPrefix(raw, r.Prefix) {
			path := raw
			if r.TrimPrefix {
				path = raw[len(r.Prefix):]
			}
			if r.HasSpaces && strings.HasPrefix(path, "spaces/") {
				path = path[len("spaces/"):]
			}

			parts := strings.SplitN(path, "/", r.MinParts+1)
			if len(parts) < r.MinParts {
				continue
			}

			res.Platform = r.Platform
			res.Owner = parts[0]
			if r.MinParts > 1 {
				res.Repo = trimGitSuffix(parts[1])
			}
			if r.AllowExtra && len(parts) > r.MinParts {
				res.Extra = parts[r.MinParts]
			}
			return res, true
		}
	}
	return
}
