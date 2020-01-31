package main

import (
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

type safeDir string

func (sd safeDir) Open(path string) (http.File, error) {
	var err error
	defer func() {
		logger.WithFields(logrus.Fields{"path": path, "error": err}).Debug("handling request")
	}()
	dir := http.Dir(sd)
	file, err := dir.Open(path)
	if err != nil {
		if os.IsNotExist(err) && conf.Path404 != "" {
			return dir.Open(conf.Path404)
		}
		return nil, err
	}
	if conf.ShowList {
		return file, nil
	}
	return safeFile{file}, nil
}

type safeFile struct {
	http.File
}

func (sf safeFile) Readdir(count int) ([]os.FileInfo, error) {
	// always show empty dir-tree
	return []os.FileInfo{}, nil
}
