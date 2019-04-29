package main

import (
	"io/ioutil"
	"os"
	"strings"
	"time"
)

func parseDir(curPath string, paths *[]string) {
	fileInfos, err := ioutil.ReadDir(curPath)
	if err != nil {
		return
	}

	useDir := false
	for _, v := range fileInfos {

		if v.IsDir() && v.Name()[0] != '.' {
			parseDir(curPath+"/"+v.Name(), paths)
			continue
		}
		if useDir == true {
			continue
		}

		*paths = append(*paths, curPath)
		useDir = true
	}
}

func getFileModTime(path string) int64 {
	path = strings.Replace(path, "\\", "/", -1)
	f, err := os.Open(path)
	if err != nil {
		return time.Now().Unix()
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return time.Now().Unix()
	}

	return fi.ModTime().Unix()
}

func checkIsIgnore(filename string) bool {
	if len(cfg.Ignores) > 0 {
		for _, v := range cfg.Ignores {
			if v == filename {
				return true
			} else {
				continue
			}
		}
	}

	if len(cfg.IgnoreRegexp) > 0 {
		for _, r := range cfg.IgnoreRegexp {
			if r.MatchString(filename) {
				return true
			} else {
				continue
			}
		}
	}

	return false
}

func checkIsExt(name string) bool {
	for _, s := range cfg.Exts {
		if strings.HasSuffix(name, s) {
			return true
		}
	}
	return false
}

func fileExist(fpath string) bool {
	_, err := os.Stat(fpath)
	return err == nil || os.IsExist(err)
}

func isIn(sl []string, tg string) bool {
	for _, v := range sl {
		if v == tg {
			return true
		}
	}
	return false
}
