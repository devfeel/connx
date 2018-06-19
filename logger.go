package connx

import (
	"fmt"
	"time"
	"path/filepath"
	"os"
	"syscall"
	"log"
	"strings"
)

const (
	LogLevel_Debug = "debug"
	LogLevel_Info  = "info"
	LogLevel_Warn  = "warn"
	LogLevel_Error = "error"
)


const (
	defaultTimeLayout = "2006_01_02"
	defaultDateFormatForFileName = "2006_01_02"
	defaultFullTimeLayout        = "2006-01-02 15:04:05.9999"

)

var connLogger Logger

func init(){
	connLogger = newXLog()
}

// SetLogger set outer logger replace xLog
func SetLogger(logger Logger){
	connLogger = logger
}

type Logger interface {
	SetEnabledLog(enabledLog bool)
	Debug(log string)
	Info(log string)
	Warn(log string)
	Error(log string)
}

type chanLog struct {
	Content   string
	LogLevel  string
}

type xLog struct {
	logRootPath    string
	logChan_Custom chan chanLog
	enabledLog     bool
}

//create new xLog
func newXLog() *xLog {
	l := &xLog{logChan_Custom: make(chan chanLog, 1000)}
	go l.handleCustom()
	return l
}

//SetEnabledLog set enabled log
func  (l *xLog) SetEnabledLog(isLog bool) {
	l.enabledLog = isLog
}



// Debug debug log with default format
func (l *xLog) Debug(log string) {
	l.log(log, LogLevel_Debug)
}

// Info info log with default format
func (l *xLog) Info(log string) {
	l.log(log, LogLevel_Info)
}

// Warn warn log with default format
func (l *xLog) Warn(log string) {
	l.log(log, LogLevel_Warn)
}

// Error error log with default format
func (l *xLog) Error(log string) {
	l.log(log, LogLevel_Error)
}

// log push log into chan
func (l *xLog) log(log string, logLevel string) {
	if l.enabledLog {
		chanLog := chanLog{
			Content:   log,
			LogLevel:  logLevel,
		}
		l.logChan_Custom <- chanLog
	}
}


//处理日志内部函数
func (l *xLog) handleCustom() {
	for {
		log := <-l.logChan_Custom
		l.writeLog(log, "custom")
	}
}

func (l *xLog) writeLog(chanLog chanLog, level string) {
	fileName := getCurrentDirectory() + chanLog.LogLevel
	switch level {
	case "custom":
		fileName = fileName + "_" + time.Now().Format(defaultDateFormatForFileName) + ".log"
		break
	}
	log := chanLog.Content
	log = fmt.Sprintf("%s [%s] %s", time.Now().Format(defaultFullTimeLayout), chanLog.LogLevel, chanLog.Content)
	fmt.Println(log)
	writeFile(fileName, log)
}

func writeFile(fileName string, log string) {
	pathDir := filepath.Dir(fileName)
	if !existsFile(pathDir) {
		//create path
		err := os.MkdirAll(pathDir, 0777)
		if err != nil {
			fmt.Println("xlog.writeFile create path error ", err)
			return
		}
	}

	var mode os.FileMode
	flag := syscall.O_RDWR | syscall.O_APPEND | syscall.O_CREAT
	mode = 0666
	logstr := log + "\r\n"
	file, err := os.OpenFile(fileName, flag, mode)
	defer file.Close()
	if err != nil {
		fmt.Println(fileName, err)
		return
	}
	file.WriteString(logstr)
}


func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

//check filename is exist
func existsFile(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}