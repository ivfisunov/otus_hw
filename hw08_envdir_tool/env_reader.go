package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Environment map[string]string

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envVars := make(Environment)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error reading dir - %s: %w", dir, err)
	}
	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			s, _ := readEnvVar(filepath.Join(dir, fileName))
			envVars[fileName] = s
		}
	}
	return envVars, nil
}

func readEnvVar(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file - %s: %w", filePath, err)
	}
	defer f.Close()

	r := bufio.NewScanner(f)
	r.Scan()
	str := r.Text()
	str = string(bytes.ReplaceAll([]byte(str), []byte("\x00"), []byte("\n")))

	return str, nil
}
