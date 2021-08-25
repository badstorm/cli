package module

import (
	"errors"
	"io/ioutil"
	"strings"

	"golang.org/x/mod/modfile"
)

var (
	ErrModuleNameNotFound = errors.New("module name not found")
)

// Returns module name for folder name
func FolderName() (string, error) {
	name, err := Name()
	if err != nil {
		return "", err
	}

	parts := strings.Split(name, "/")
	return parts[len(parts)-1], nil
}

//Name of the module specified in the modfile
func Name() (string, error) {
	content, err := ioutil.ReadFile("go.mod")
	if err != nil {
		return "", err
	}

	path := modfile.ModulePath(content)
	if path == "." {
		return "", ErrModuleNameNotFound
	}

	return path, nil
}
