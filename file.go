package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type FileInfo struct {
	Name     string
	Size     int64
	FullPath string
	Md5      string
}

func (f *FileInfo) GetFiles(dir string) []*FileInfo {
	files := make([]*FileInfo, 0)
	rd, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}
	for _, info := range rd {
		if info.IsDir() {
			fs := f.GetFiles(filepath.Join(dir, info.Name()))
			if fs != nil {
				files = append(files, fs...)
			}
		} else {
			file := &FileInfo{
				Name:     info.Name(),
				Size:     info.Size(),
				FullPath: filepath.Join(dir, info.Name()),
				Md5:      f.fileMD5(filepath.Join(dir, info.Name())),
			}
			files = append(files, file)
		}
	}
	return files
}

func (f *FileInfo) fileMD5(path string) string {
	file, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer file.Close()

	m := md5.New()
	_, err = io.Copy(m, file)
	if err != nil {
		return ""
	}
	return hex.EncodeToString(m.Sum(nil))
}
