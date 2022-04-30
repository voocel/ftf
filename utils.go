package main

import (
	"os"
	"strings"
)

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

func isFile(path string) bool {
	return !isDir(path)
}

func filterDuplicate(src []string) []string {
	i := 0
	m := make(map[string]struct{})
	for _, item := range src {
		if _, ok := m[item]; !ok && item != "" {
			item = strings.TrimSpace(item)
			m[item] = struct{}{}
			src[i] = item
			i++
		}
	}
	return src[:i]
}
