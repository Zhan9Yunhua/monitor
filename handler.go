package main

import (
	"github.com/Zhan9Yunhua/logger"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	fsn "github.com/howeyc/fsnotify"
)

var (
	cmd          *exec.Cmd
	mutex        sync.Mutex
	modTimeMap   = make(map[string]int64)
	scheduleTime time.Time
)

func mon(paths []string) {
	watcher, err := fsn.NewWatcher()
	if err != nil {
		logger.Fatalln("Fail to new watcher")
	}
	logger.Successln("Initializing ...")

	go monHandler(watcher)

	for _, path := range paths {
		err = watcher.Watch(path)
		if err != nil {
			logger.Fatalln("Fail to watch directory: " + err.Error())
		}
	}
}

func monHandler(w *fsn.Watcher) {
	for {
		select {
		case e := <-w.Event:
			isRun := true

			if checkIsIgnore(e.Name) {
				continue
			}
			if !checkIsExt(e.Name) {
				continue
			}

			mt := getFileModTime(e.Name)
			if t := modTimeMap[e.Name]; mt == t {
				isRun = false
			}

			modTimeMap[e.Name] = mt

			if isRun {
				go func() {
					scheduleTime = time.Now().Add(1 * time.Second)
					for {
						time.Sleep(scheduleTime.Sub(time.Now()))
						if time.Now().After(scheduleTime) {
							break
						}
						return
					}
					run()
				}()
			}
		case err := <-w.Error:
			logger.Fatalln(err)
		}
	}
}

func run() {
	if err := os.Chdir(curPath); err != nil {
		panic(err)
	}

	if cfg.Lang == "go" {
		build()
	} else {
		restart()
	}
}

func build() {
	mutex.Lock()
	defer mutex.Unlock()

	logger.Infoln("Start building...")

	bcmd := exec.Command("go", argsHandler()...)
	bcmd.Env = append(os.Environ(), "GOGC=off")
	bcmd.Stdout = os.Stdout
	bcmd.Stderr = os.Stderr
	if err := bcmd.Run(); err != nil {
		logger.Fatalln("!! Build failed !!")
		return
	}

	logger.Successln("Build was successful")
	restart()
}

func restart() {
	kill()
	go start()
}

func kill() {
	defer func() {
		if e := recover(); e != nil {
			logger.Infoln("Kill recover -> ", e)
		}
	}()

	if cmd != nil && cmd.Process != nil {
		if err := cmd.Process.Kill(); err != nil {
			logger.Infof("** Kill -> Pid: %d. %s\n", cmd.Process.Pid, err)
		}
	}
}

func start() {
	logger.Infof("Restarting %s ...\n", cfg.AppName)

	if cfg.Lang == "go" {
		cmd = exec.Command(cfg.Output)
		cmd.Args = append([]string{cfg.Output}, cfg.CmdArgs...)
	} else {
		scr := strings.Split(cfg.Script, " ")
		cmd = exec.Command(scr[0], strings.Join(scr[1:], " "))
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), cfg.Envs...)

	go cmd.Run()
	logger.Infof("%s is running...\n", cfg.AppName)
}

// 处理args
func argsHandler() []string {
	args := []string{"build"}
	args = append(args, "-o", cfg.Output)
	if cfg.BuildTags != "" {
		args = append(args, "-tags", cfg.BuildTags)
	}
	return args
}
