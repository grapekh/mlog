//
// mlog is simple logging module for go, Provide a standard interface for output to the screen and a logfile. 
// It includes a rotating file feature and console logging.
// The message in the log file will be prepended with the severity and also the standard LOG time/date stamp and of course, the message
// The message displayed on the screen, however, will not include the time/datestamp - only the severity and the message itself.  
//
// Options to this package - (added 07/27/18)
// You can call this with or without a log file - if a log file is specified: 
// By default, the log will be rolled over to a backup file when its size reaches 10Mb and 10 such files will be created (and eventually reused).
// Alternatively, you can specify the max size of the log file before it gets rotated, and the number of backup files you want to create, with the StartEx function.
//
// Some example calls to this pagage: 
//
// You can decide what to log when calling mlog.Start
//   a) LevelTrace logs:   everything
//   b) LevelInfo  logs:   info, Warnings and Errors
//   c) LevelWarn  logs:   Warning and Errors
//   d) LevelError logs:   just Errors
//
// Example for startx - Write everything to console and don't write to a file.
//	mlog.Start(mlog.LevelInfo, "")
//
// Example for Start - Write everything to console and write Info, Warning and Error messages to file. 
//			since it is not specified, the default file size will be max of 10Mb and rotate to 10 backup log files
//		mlog.Start(mlog.LevelTrace, "app.log")
//
// Example StartEx - Start with Extra information: 
//					Write everything to console;
//					write Info, Warning and Error messages to file
//					set log file to 5mb and only keep 4 copies of old logs. 
// 		mlog.StartEx(mlog.LevelInfo, "app.log", 5*1024*1024, 4)
//
// Sample calls: 
//
// Info: 	
//	 // write to screen and logfile. 
//	 mlog.Info("This is a piece of information")
// 
// Warning: 
//	 // write to screen and logfile. 
//	 foo := "foo"
//   mlog.Warning("BlaBla - %s", foo)
//
// Trace: 
//	 // write to screen and logfile. Note, this has newlines around it. 
//   mlog.Trace"This is a trace")
//
// Error (use standard Go errors structures
//	 err := errors.New("This is an error to display")
//	 mlog.Error(err))
//
// Fatal - display the message and STOP execution!
//   mlog.Fatal("Testing fatal error message.  This should hault the program!!")
//
package mlog

import (
	"fmt"
	//"io"					// needed if we use multiWriter
	"io/ioutil"
	"log"
	"os"
	"path"
	"sync/atomic"
)

// LogLevel type
type LogLevel int32

const (
	// LevelTrace logs everything
	LevelTrace LogLevel = (1 << iota)

	// LevelInfo logs Info, Warnings and Errors
	LevelInfo

	// LevelWarn logs Warning and Errors
	LevelWarn

	// LevelError logs just Errors
	LevelError
)

const MaxBytes int = 10 * 1024 * 1024
const RotateCount int = 10

type mlog struct {
	LogLevel int32

	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	Fatal   *log.Logger

	LogFile *RotatingFileHandler
}

var Logger mlog

// DefaultFlags used by created loggers
var DefaultFlags = log.Ldate | log.Ltime | log.Lshortfile

//
// RotatingFileHandler writes log a file, if file size exceeds maxBytes,
// It will rotate current file and add a version number to it and reopen a new one.
//
// max backup file number is set by rotateCount, 
// it will delete oldest file if number of rotated files exceeds number specified
//
type RotatingFileHandler struct {
	fd *os.File

	fileName    string
	maxBytes    int
	rotateCount int
}

// 
// NewRotatingFileHandler creates dirs and opens the logfile
//
func NewRotatingFileHandler(fileName string, maxBytes int, rotateCount int) (*RotatingFileHandler, error) {
	dir := path.Dir(fileName)
	os.Mkdir(dir, 0777)

	h := new(RotatingFileHandler)

	if maxBytes <= 0 {
		return nil, fmt.Errorf("invalid max bytes")
	}

	h.fileName = fileName
	h.maxBytes = maxBytes
	h.rotateCount = rotateCount

	var err error
	h.fd, err = os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return h, nil
}

func (h *RotatingFileHandler) Write(p []byte) (n int, err error) {
	h.doRollover()
	return h.fd.Write(p)
}

// Close simply closes the File
func (h *RotatingFileHandler) Close() error {
	if h.fd != nil {
		return h.fd.Close()
	}
	return nil
}

//
// utility function to handle roll over of logflies. 
// Append number to end of filename
//
func (h *RotatingFileHandler) doRollover() {
	f, err := h.fd.Stat()
	if err != nil {
		return
	}

	// log.Println("The logfile size: ", f.Size())

	if h.maxBytes <= 0 {
		return
	} else if f.Size() < int64(h.maxBytes) {
		return
	}

	if h.rotateCount > 0 {
		h.fd.Close()

		for i := h.rotateCount - 1; i > 0; i-- {
			sfn := fmt.Sprintf("%s.%d", h.fileName, i)
			dfn := fmt.Sprintf("%s.%d", h.fileName, i+1)

			os.Rename(sfn, dfn)
		}

		dfn := fmt.Sprintf("%s.1", h.fileName)
		os.Rename(h.fileName, dfn)

		h.fd, _ = os.OpenFile(h.fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	}
}

// Start starts the logging
func Start(level LogLevel, path string) {
	doLogging(level, path, MaxBytes, RotateCount)
}

func StartEx(level LogLevel, path string, maxBytes, rotateCount int) {
	doLogging(level, path, maxBytes, rotateCount)
}

// Stop stops the logging
func Stop() error {
	if Logger.LogFile != nil {
		return Logger.LogFile.Close()
	}

	return nil
}

//Sync commits the current contents of the file to stable storage.
//Typically, this means flushing the file system's in-memory copy
//of recently written data to disk.
func Sync() {
	if Logger.LogFile != nil {
		Logger.LogFile.fd.Sync()
	}
}

func doLogging(logLevel LogLevel, fileName string, maxBytes, rotateCount int) {
	traceHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	errorHandle := ioutil.Discard
	fatalHandle := ioutil.Discard

	var fileHandle *RotatingFileHandler

	switch logLevel {
	case LevelTrace:
		traceHandle = os.Stdout
		fallthrough
	case LevelInfo:
		infoHandle = os.Stdout
		fallthrough
	case LevelWarn:
		warnHandle = os.Stdout
		fallthrough
	case LevelError:
		errorHandle = os.Stderr
		fatalHandle = os.Stderr
	}

	if fileName != "" {
		var err error
		fileHandle, err = NewRotatingFileHandler(fileName, maxBytes, rotateCount)
		if err != nil {
			log.Fatal("mlog: unable to create RotatingFileHandler: ", err)
		}

		// MultiWriter acts like the "tee" unix function - write to two streams. 
		// we can write to screen as well as logfile, but for our case, we don't 
		// want the timestamp on the console.  I'm leaving this here commented out just
		// in case I might want to add it back later. 

		if traceHandle == os.Stdout {
			//traceHandle = io.MultiWriter(fileHandle, traceHandle)
			traceHandle = fileHandle
		}

		if infoHandle == os.Stdout {
			//infoHandle = io.MultiWriter(fileHandle, infoHandle)
			infoHandle = fileHandle
		}

		if warnHandle == os.Stdout {
			//warnHandle = io.MultiWriter(fileHandle, warnHandle)
			warnHandle = fileHandle
		}

		if errorHandle == os.Stderr {
			//errorHandle = io.MultiWriter(fileHandle, errorHandle)
			errorHandle = fileHandle
		}

		if fatalHandle == os.Stderr {
			//fatalHandle = io.MultiWriter(fileHandle, fatalHandle)
			fatalHandle = fileHandle
		}
		
	}

	Logger = mlog{
		Trace:   log.New(traceHandle, " - [TRACE]: ", DefaultFlags),
		Info:    log.New(infoHandle,  " - [INFO]: ", DefaultFlags),
		Warning: log.New(warnHandle,  " - [WARN]: ", DefaultFlags),
		Error:   log.New(errorHandle, " - [ERROR]: ", DefaultFlags),
		Fatal:   log.New(errorHandle, " - [FATAL]: ", DefaultFlags),
		LogFile: fileHandle,
	}

	atomic.StoreInt32(&Logger.LogLevel, int32(logLevel))
}

// Trace writes to the Trace destination as well as the console
func Trace(format string, a ...interface{}) {
	Logger.Trace.Output(2, fmt.Sprintf(format, a...))
	fmt.Println(" - [TRACE]: ", fmt.Sprintf(format, a...))	
}

// Info writes to the Info destination as well as the console
func Info(format string, a ...interface{}) {
	Logger.Info.Output(2, fmt.Sprintf(format, a...))
	fmt.Println(" - [INFO]: ", fmt.Sprintf(format, a...))	
}

// Warning writes to the Warning destination as well as the console
func Warning(format string, a ...interface{}) {
	Logger.Warning.Output(2, fmt.Sprintf(format, a...))
	fmt.Println(" - [WARNING]: ", fmt.Sprintf(format, a...))	
}

// Error writes to the Error destination as well as the console and accepts an err
func Error(err error) {
	Logger.Error.Output(2, fmt.Sprintf("%s\n", err))
	fmt.Println(" - [ERROR]: ", fmt.Sprintf("%s", err))	
}

// IfError is a shortcut function for log.Error if error
// writes to destination as well as the console
func IfError(err error) {
	if err != nil {
		Logger.Error.Output(2, fmt.Sprintf("%s\n", err))
		fmt.Println(" - [ERROR]: ", fmt.Sprintf("%s", err))	 
	}
}

// Fatal writes to the Fatal destination and exits with an error 255 code
func Fatal(a ...interface{}) {
	fmt.Println(" - [FATAL]: ", fmt.Sprintf("%s", a...))	 
	Logger.Fatal.Output(2, fmt.Sprint(a...))
	Sync()
	os.Exit(255)
}

// Fatalf writes to the Fatal destination and exits with an error 255 code
func Fatalf(format string, a ...interface{}) {
	fmt.Println(" - [FATAL]: ", fmt.Sprintf("%s", a...))	 
	Logger.Fatal.Output(2, fmt.Sprintf(format, a...))
	Sync()
	os.Exit(255)
}

// FatalIfError is a shortcut function for log.Fatalf if error and
// exits with an error 255 code
func FatalIfError(err error) {
	if err != nil {
		fmt.Println(" - [FATAL]: ", fmt.Sprintf("%s", err))	 
		Logger.Fatal.Output(2, fmt.Sprintf("%s\n", err))
		Sync()
		os.Exit(255)
	}
}