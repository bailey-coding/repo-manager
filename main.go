package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"sigs.k8s.io/yaml"
)

type State struct {
	Repos map[string]*Repo `json:"repos"`
}

type Repo struct {
	Remotes map[string]string`json:"remotes"`
}

func readState() State {
	var result State
	data, err := ioutil.ReadFile("state.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			return State{}
		}
		panic(err)
	}
	err = yaml.Unmarshal(data, &result)
	if err != nil {
		panic(err)
	}
	return result
}

func writeState(data State) {
	y, err := yaml.Marshal(data)
	if err != nil {
		panic(err)
	}
	ioutil.WriteFile("state.yaml", y, 0666)

}

func main() {
	state := readState()
	if state.Repos == nil {
		state.Repos = make(map[string]*Repo)
	}
	path := "~/src"
	maxDepth := 1

	usr, _ := user.Current()
	dir := usr.HomeDir

	if path == "~" {
		// In case of "~", which won't be caught by the "else if"
		path = dir
	} else if strings.HasPrefix(path, "~/") {
		// Use strings.HasPrefix so we don't match paths like
		// "/something/~/something/"
		path = filepath.Join(dir, path[2:])
	}
	startDepth := strings.Count(path, "/")
	maxDepth = maxDepth + startDepth

	fmt.Println("On Unix:")
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if strings.Count(path, "/") > maxDepth {
			fmt.Printf("skipping a dir because beyond max depth: %+v \n", info.Name())
			return filepath.SkipDir
		}

		gitDir := filepath.Join(path, ".git")
		if f, err := os.Stat(gitDir); err == nil {
			ft := "file"
			if f.IsDir() {
				ft = "dir"
			}
			fmt.Printf("found git %s: %q\n", ft, path)
			r, err := git.PlainOpen(gitDir)
			if err != nil {
				fmt.Printf("error opening git repo at %q: %v", gitDir, err)
			}
			remotes, err := r.Remotes()
			if err != nil {
				fmt.Printf("error listing remotes for git repo at %q: %v", gitDir, err)
			}
			for _, remote := range remotes {
				config := remote.Config()
				if len(config.URLs) != 1 {
					fmt.Printf("Found multiple URLs for %+v at %s\n", config.Name, path)
					for _, url := range config.URLs {
						fmt.Printf("%+v\n", url)
					}
				}
				path = unexpandHome(usr.HomeDir, path)
				if state.Repos[path] == nil {
					state.Repos[path] = &Repo{}
				}
				if state.Repos[path].Remotes == nil {
					state.Repos[path].Remotes = make(map[string]string)
				}
				state.Repos[path].Remotes[config.Name] = config.URLs[0]
			}
			return filepath.SkipDir
		} else if !os.IsNotExist(err) {
			fmt.Printf("error checking for gitDir %q: %v", gitDir, err)
		}
		return nil
	})
	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}
	writeState(state)
}

func unexpandHome(home, path string) string {
	return strings.ReplaceAll(path, home, "~")
}
