package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	url := getUrlFromGitRemote()
	openUrlInBrowser(url)
}

func openUrlInBrowser(url string) {
	cmd := exec.Command("open", url)
	output, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error executing 'open %s' command\n", url)
		fmt.Println(err.Error())
	}
	if string(output) != "" {
		fmt.Println("open output")
		fmt.Println(string(output))
	}
}
func getUrlFromGitRemote() string {
	cmd := exec.Command("git", "remote", "-v")
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing 'git remote -v' command")
		fmt.Println(err.Error())
	}
	url := parseUrl(string(output))
	return url
}

type githubOrigin struct {
	domain    string
	repoName  string
	repoOwner string
}

func (gh *githubOrigin) getUrl() string {
	return fmt.Sprintf("http://%s/%s/%s", gh.domain, gh.repoOwner, gh.repoName)
}

type gitUrlRepresentation int

const (
	ssh gitUrlRepresentation = iota
	https
)

func (gur gitUrlRepresentation) String() string {
	switch gur {
	case ssh:
		return "ssh"
	case https:
		return "https"
	default:
		return "unsupported git url rep"
	}
}

func getGitUrlRepr(gitRemoteOutput string) gitUrlRepresentation {
	match, _ := regexp.MatchString("https://", gitRemoteOutput)
	if match {
		return https
	} else {
		return ssh
	}
}

// parseUrl get url from 'git remote -v' output
func parseUrl(gitRemoteOutput string) string {
	gh := githubOrigin{}
	gitUrlRepr := getGitUrlRepr(gitRemoteOutput)
	switch gitUrlRepr {
	case ssh:
		gh.domain = parseGithubDomainSsh(gitRemoteOutput)
		gh.repoName = parseGithubRepoNameSsh(gitRemoteOutput)
		gh.repoOwner = parseGithubRepoOwnerSsh(gitRemoteOutput)
		return gh.getUrl()
	case https:
		gh.domain = parseGithubDomainHttps(gitRemoteOutput)
		gh.repoName = parseGithubRepoNameHttps(gitRemoteOutput)
		gh.repoOwner = parseGithubRepoOwnerHttps(gitRemoteOutput)
		println("https git remotes not supported yet")
		os.Exit(1)
		return gh.getUrl()
	default:
		return "invalid gitUrlRepr"
	}
}

func parseGithubDomainSsh(gitSshUrl string) string {
	r, _ := regexp.Compile(`.*fetch`)
	result1 := r.FindString(gitSshUrl)
	r, _ = regexp.Compile(`@.*\.com`)
	result2 := r.FindString(result1)
	r, _ = regexp.Compile(`[^@].*`)
	result3 := r.FindString(result2)
	return result3
}

func parseGithubRepoOwnerSsh(gitSshUrl string) string {
	r, _ := regexp.Compile(`:.*/`)
	result1 := r.FindString(gitSshUrl)
	r, _ = regexp.Compile(`[^:][a-z]*`)
	result2 := r.FindString(result1)
	return result2
}

func parseGithubRepoNameSsh(gitSshUrl string) string {
	r, _ := regexp.Compile(`/.*\.git`)
	result1 := r.FindString(gitSshUrl)
	r, _ = regexp.Compile(`[^/][\w|\d|-]*`)
	result2 := r.FindString(result1)
	return result2
}

func parseGithubDomainHttps(gitHttpsUrl string) string {
	r, _ := regexp.Compile(`https://\w*\.com`)
	result1 := r.FindString(gitHttpsUrl)
	return result1
}

func parseGithubRepoOwnerHttps(gitHttpsUrl string) string {
	r, _ := regexp.Compile(`com/\w*/`)
	result1 := r.FindString(gitHttpsUrl)
	r, _ = regexp.Compile(`[^com/]\w*`)
	result2 := r.FindString(result1)
	return result2
}

func parseGithubRepoNameHttps(gitHttpsUrl string) string {
	r, _ := regexp.Compile(`\.*/.*\.git`)
	result1 := r.FindString(gitHttpsUrl)
	r, _ = regexp.Compile(`[^com/]\w*`)
	result2 := r.FindString(result1)
	return result2
}
