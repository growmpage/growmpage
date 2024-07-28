package growhelper

import (
	"fmt"
	"net/http"
	"os/exec"
	"strings"
)

func GitPrivateVersion(remote bool) string {
	command := &exec.Cmd{}
	if remote {
		command = exec.Command("/bin/bash", "-c", Filename_git_private_version_origin)
	} else {
		command = exec.Command("/bin/bash", "-c", Filename_git_private_version_local)
	}
	out, err := command.CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return fmt.Sprintf("ERROR-could not fetch origin version with remote %v: %v", remote, err.Error())
	}
	fmt.Printf("GetGitVersion(%v): %v\n", remote, string(out))
	lines := strings.Split(string(out), "\n")
	fmt.Printf("lines: %v\n", lines)
	return lines[len(lines)-2]
}

func GitPublicVersion() string {
	out, err := exec.Command("/bin/bash", "-c", Filename_git_public_version_origin).CombinedOutput()
	if err != nil {
		fmt.Println(err)
		return "ERROR: could not fetch origin version"
	}
	fmt.Printf("gitPublicVersion(): %v\n", string(out))
	lines := strings.Split(string(out), "\n")
	return lines[len(lines)-2]
}

func GitPublicUpgrade(w http.ResponseWriter, local string) { //TODO: return string
	origin := GitPublicVersion()
	if origin != local {
		GithubPublicInstall(nil, origin)
		fmt.Fprintf(w, "upgraded from %v to %v, reloading growmpage...", local, origin)
	} else {
		fmt.Fprint(w, "binary allready up to date: "+local)
	}
}

func GithubPublicInstall(w http.ResponseWriter, version string) {
	command := exec.Command("bash", "-c", Filename_github_public_install+" "+version)
	_, err := command.CombinedOutput()
	if err != nil && w != nil {
		fmt.Fprint(w, err)
	}
	// lines := strings.Split(string(out), "\n")
	// fmt.Fprint(w, lines[len(lines)-2])
}

func GitPrivateReset(w http.ResponseWriter) {
	local := GitPrivateVersion(false)
	origin := GitPrivateVersion(true)
	if origin == local {
		fmt.Fprint(w, "allready up to date: "+local)
		return
	} else {
		command := exec.Command("/bin/bash", "-c", Filename_git_private_reset)
		_, err := command.CombinedOutput()
		if err != nil {
			fmt.Fprint(w, err)
		}
		fmt.Fprint(w, "updated from "+local+" to "+origin+", reloading growmpage...")
	}
}

func GitPrivateBackup() {
	command := exec.Command("/bin/bash", "-c", Filename_git_private_backup)
		_, err := command.CombinedOutput()
		if err != nil {
			fmt.Println(err.Error())
		}
}