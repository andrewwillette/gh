package main

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
)

func main() {
	cmd := exec.Command("git", "remote", "-v")
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error executing 'git remote -v' command")
		fmt.Println(err.Error())
	}
	url := parseUrlFromSsh(string(output))
	fmt.Printf("output\n%+v\n", url)
	cmd = exec.Command("open", url)
	output, err = cmd.Output()
	if err != nil {
		fmt.Printf("Error executing 'open %s' command\n", url)
		fmt.Println(err.Error())
	}
}

type githubOrigin struct {
	domain    string
	repoName  string
	repoOwner string
}

func (gh *githubOrigin) getUrl() string {
	return fmt.Sprintf("http://%s/%s/%s", gh.domain, gh.repoOwner, gh.repoName)
}

func parseUrlFromSsh(gitSshUrl string) string {
	gh := githubOrigin{}
	gh.domain = parseGithubDomain(gitSshUrl)
	gh.repoName = parseGithubRepoName(gitSshUrl)
	gh.repoOwner = parseGithubRepoOwner(gitSshUrl)
	return gh.getUrl()
}

func parseGithubDomain(gitSshUrl string) string {
	r, _ := regexp.Compile(`.*fetch`)
	result1 := r.FindString(gitSshUrl)
	r, _ = regexp.Compile(`@.*\.com`)
	result2 := r.FindString(result1)
	r, _ = regexp.Compile(`[^@].*`)
	result3 := r.FindString(result2)
	return result3
}
func parseGithubRepoOwner(gitSshUrl string) string {
	r, _ := regexp.Compile(`:.*/`)
	result1 := r.FindString(gitSshUrl)
	r, _ = regexp.Compile(`[^:][a-z]*`)
	result2 := r.FindString(result1)
	return result2
}
func parseGithubRepoName(gitSshUrl string) string {
	r, _ := regexp.Compile(`/.*\.git`)
	result1 := r.FindString(gitSshUrl)
	r, _ = regexp.Compile(`[^/][\w|\d]*`)
	result2 := r.FindString(result1)
	return result2
}
