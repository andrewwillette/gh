package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"runtime/pprof"

	"github.com/andrewwillette/gocommon"
	"github.com/rs/zerolog/log"
)

type gitUrlRepresentation int

const (
	ssh gitUrlRepresentation = iota
	https
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
)

func main() {
	flag.Parse()
	gocommon.ConfigureConsoleZerolog()
	configureProfiling()
	defer pprof.StopCPUProfile()
	url := getUrlFromGitRemote()
	openUrlInBrowser(url)
}

func configureProfiling() {
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Err(err)
		}
		runtime.SetCPUProfileRate(1000000)
		pprof.StartCPUProfile(f)
	}
}

func openUrlInBrowser(url string) {
	cmd := exec.Command("open", url)
	output, err := cmd.Output()
	if err != nil {
		log.Err(err)
		log.Error().Msgf("url: %s", url)
	}
	if string(output) != "" {
		log.Warn().Msgf("open url returned output: %s", string(output))
	}
}

func getUrlFromGitRemote() string {
	cmd := "git remote -v | grep push"
	out, err := exec.Command("bash", "-c", cmd).Output()
	log.Debug().Msgf("git remote -v output: %s", out)
	if err != nil {
		log.Error().Msg("Error executing 'git remote -v' command")
		log.Err(err)
	}
	url := parseUrl(string(out))
	return url
}

type githubOrigin struct {
	domain    string
	repoName  string
	repoOwner string
}

func (gh *githubOrigin) getUrl() string {
	log.Debug().Msg("githubOrigin#getUrl")
	log.Debug().Msgf("%+v", gh)
	return fmt.Sprintf("http://%s/%s/%s", gh.domain, gh.repoOwner, gh.repoName)
}

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
		log.Debug().Msg("git remote is of type https.")
		return https
	} else {
		log.Debug().Msg("git remote is of type ssh.")
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
	log.Debug().Msg("parseGithubDomainSsh")
	r, _ := regexp.Compile(`.*push`)
	result1 := r.FindString(gitSshUrl)
	log.Debug().Msgf("result1: %s", result1)
	r, _ = regexp.Compile(`@([^:].)*`)
	result2 := r.FindString(result1)
	log.Debug().Msgf("result2: %s", result2)
	r, _ = regexp.Compile(`[^@].*`)
	result3 := r.FindString(result2)
	log.Debug().Msgf("result3: %s", result3)
	return result3
}

func parseGithubRepoOwnerSsh(gitSshUrl string) string {
	log.Debug().Msg("parseGithubRepoOwnerSsh")
	r, _ := regexp.Compile(`:.*/`)
	result1 := r.FindString(gitSshUrl)
	log.Debug().Msgf("result1: %s", result1)
	r, _ = regexp.Compile(`[^:][\w|\d|-|\.]*`)
	result2 := r.FindString(result1)
	log.Debug().Msgf("result2: %s", result2)
	return result2
}

func parseGithubRepoNameSsh(gitSshUrl string) string {
	log.Debug().Msgf("parseGithubRepoNameSsh %s", gitSshUrl)
	r, _ := regexp.Compile(`/.*\.git`)
	result1 := r.FindString(gitSshUrl)
	log.Debug().Msgf("result1: %s", result1)
	r, _ = regexp.Compile(`[^/](\w|\d|-|\.)*[^(.git)]`)
	result2 := r.FindString(result1)
	log.Debug().Msgf("result2: %s", result2)
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
