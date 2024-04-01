package growhelper

import (
	"fmt"
	"net/http"
	"os"
	"strings"
)

var (
	Port               = ":8080"
	Filename_growtable = "../data/growtable.json"
)

const (
	Filename_expertpage                 = "../html/expert.html"
	Filename_growsniffer                = "../scripts/growsniffer.py"
	Filename_favicon                    = "../html/favicon.ico"
	Filename_growmpage                  = "../html/growmpage.html"
	Filename_growganizer                = "../data/growganizer.json"
	Filename_git_private_version_origin = "../scripts/git_private_version_origin.sh"
	Filename_git_private_version_local  = "../scripts/git_private_version_local.sh"
	Filename_git_public_version_origin  = "../scripts/git_public_version_origin.sh"
	Filename_github_public_install      = "../scripts/github_public_install.sh"
	Filename_git_private_reset          = "../scripts/git_private_reset.sh"
	Filename_pictures                   = "../data/camera/"
)

func Url(suffix string) string {
	if !strings.HasPrefix(suffix, "/") {
		suffix = "/" + suffix
	}
	name, err := os.Hostname()
	if err == nil && ToInt(strings.Split(Port, ":")[1]) < 9999{
		return "http://" + name + Port + suffix
	} else {
		return internalUrl(suffix)
	}
}

func internalUrl(suffix string) string {
	return "http://127.0.0.1" + Port + suffix
}

func Post(suffix string, body string) string {
	answer := extract(http.Post(Url(suffix), "application/text", strings.NewReader(body)))
	return answer
}

func Get(suffix string) string {
	answer := extract(http.Get(Url(suffix)))
	return answer
}

func extract(response *http.Response, err error) string {
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	if response != nil && response.Body != nil {
		return ReadString(response.Body)
	}
	return ""
}
