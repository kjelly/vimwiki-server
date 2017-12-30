package utils

import (
	"errors"
	"fmt"
	"github.com/russross/blackfriday"
	"html/template"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func walkDir(path string) []string {
	var ret []string
	filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".md" {
			ret = append(ret, path)
		}
		return nil
	})
	return ret
}

func getHomeDir() string {
	usr, err := user.Current()
	if err != nil {
		return os.Getenv("HOME")
	}
	return usr.HomeDir
}

// ReadFile todo
func ReadFile(name string) (string, error) {
	homeDir := getHomeDir()
	basePath := filepath.Join(homeDir, "Dropbox/vimwiki/")
	var fileList = walkDir(basePath)
	fmt.Printf("%v\n", fileList)
	path := filepath.Join(basePath, name)
	realPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(realPath, basePath) {
		return "", errors.New("Don't allow to access")
	}
	content, err := ioutil.ReadFile(realPath)
	return string(content), err
}

// MarkdownToHTML todo
func MarkdownToHTML(input string) template.HTML {
	output := blackfriday.MarkdownCommon([]byte(input))
	return template.HTML(string(output))
	return template.HTML(input)
}
