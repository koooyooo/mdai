/*
Copyright © 2025 koooyooo
*/
package file

import (
	"io/fs"
	"os"
	"path/filepath"
)

func LoadFiles(path string) (map[string][]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	filePathContents := make(map[string][]byte)
	if info.IsDir() {
		err := filepath.WalkDir(path, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				b, err := os.ReadFile(p)
				if err != nil {
					return err
				}
				filePathContents[p] = b
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		return filePathContents, nil
	} else {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		filePathContents[path] = b
		return filePathContents, nil
	}
}

// LoadContent loads the content of a file as a string
func LoadContent(path string) (string, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
