package main

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"runtime"
)

const defCfgFile = "./monitor.json"

type Config struct {
	// app name
	AppName string `json:"appName"`
	// go打包目录
	Output string `json:"output"`
	// 监听的文件后缀
	Exts []string `json:"exts"`
	// BuildPkg  string
	BuildTags string `json:"buildTags"`
	// 执行的cmd
	CmdArgs []string `json:"cmdArgs"`
	// 环境变量
	Envs []string `json:"envs"`
	// 变动执行的脚本
	Script string `json:"script"`
	// 语言
	Lang string `json:"lang"`
	// 忽略文件
	Ignores []string `json:"ignores"`
	// 忽略文件的正则
	IgnoreRegs   []string `json:"ignoreRegs"`
	IgnoreRegexp []*regexp.Regexp
}

func cfgHandler() {
	c := &Config{}

	if cfgFile == "" {
		cfgFile = defCfgFile
	}
	fpath, _ := filepath.Abs(cfgFile)
	if fileExist(fpath) {
		if fby, err := ioutil.ReadFile(fpath); err == nil {
			if json.Unmarshal(fby, c) != nil {
				panic("json unmarshal failed!")
			}
		} else {
			panic("read config file failed!")
		}
	}

	if c.Output == "" {
		outputExt := ""
		if runtime.GOOS == "windows" {
			outputExt = ".exe"
		}
		c.Output = "./" + c.AppName + outputExt
	}
	if c.Lang == "" {
		c.Lang = "go"
	}
	if c.AppName == "" {
		c.AppName = filepath.Base(curPath)
	}
	if len(c.Exts) == 0 || !isIn(c.Exts, "."+c.Lang) {
		c.Exts = append(c.Exts, "."+c.Lang)
	}
	if len(c.IgnoreRegs) > 0 {
		for _, regex := range c.IgnoreRegs {
			if r, err := regexp.Compile(regex); err == nil {
				c.IgnoreRegexp = append(c.IgnoreRegexp, r)
			} else {
				panic("Compile regex error: " + regex)
			}
		}
	}

	cfg = c
}
