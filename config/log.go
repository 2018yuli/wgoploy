package config

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type LogConfig struct {
	Path string `toml:"path"`
}

func (l *LogConfig) OnChange() error {
	l.SetDefault()
	return nil
}

func (l *LogConfig) SetDefault() {
	var logFile io.Writer
	logPathEnv := Toml.Log.Path
	if strings.ToLower(logPathEnv) == "stdout" {
		logFile = os.Stdout
	} else {
		logPath, err := filepath.Abs(logPathEnv)
		if err != nil {
			fmt.Println(err.Error())
		}
		if _, err := os.Stat(logPath); err != nil && os.IsNotExist(err) {
			if err := os.Mkdir(logPath, os.ModePerm); nil != err {
				panic(err.Error())
			}
		}
		logFile, err = os.OpenFile(logPath+"/goploy.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if nil != err {
			panic(err.Error())
		}
	}
	log.SetReportCaller(true)

	log.SetFormatter(&log.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return fmt.Sprintf("%s()", path.Base(f.Function)), fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		},
	})

	stdoutWriter := os.Stdout                            // 若要禁用控制台输出，改为ioutil.Discard；否则留空或使用os.Stdout
	log.SetOutput(io.MultiWriter(stdoutWriter, logFile)) // 若要同时输出到控制台和文件，使用io.MultiWriter
	//log.SetOutput(logFile)

	log.SetLevel(log.TraceLevel)

}
