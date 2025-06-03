/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/01/23 13:44:17
 Desc     :
*/

package util

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FileExists(filename string) bool {
	_, err := os.Stat(filename)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}

func WriteFile(filename string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.Write(content); err != nil {
		return err
	}
	return nil
}

func ReadFileToStruct(filename string, v any) error {
	file, err := ReadFile(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(content, v); err != nil {
		return err
	}
	return err
}

func ReadFile(f string) (io.ReadCloser, error) {
	if f == "-" {
		stat, err := os.Stdin.Stat()
		if err != nil {
			return nil, fmt.Errorf("can't stat stdin: %s", err.Error())
		}
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			return os.Stdin, nil
		}
		return nil, fmt.Errorf("can't read stdin")
	}
	variants := []string{f}
	for _, fn := range variants {
		if _, err := os.Stat(fn); err != nil {
			continue
		}
		fp, err := filepath.Abs(fn)
		if err != nil {
			return nil, err
		}
		file, err := os.Open(fp)
		if err != nil {
			return nil, err
		}
		return file, nil
	}

	return nil, fmt.Errorf("failed to read file: %s", f)
}

func FileToBytes(f string) ([]byte, error) {
	file, err := ReadFile(f)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func FileToString(f string, trimSpace bool) (string, error) {
	file, err := ReadFile(f)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	if trimSpace {
		return strings.TrimSpace(string(content)), nil
	}
	return string(content), nil
}

// listDirectory lists directory contents
func ListDir(path string) []string {
	var items []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return items
	}

	for _, entry := range entries {
		items = append(items, entry.Name())
	}

	return items
}
