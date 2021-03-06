
package golog

import (
    "io"
    "io/ioutil"
    "log"
    "os"
    "fmt"
    "strings"
    "time"
    "sync/atomic"
    "errors"
)
const (
	// everything
	LevelTrace int32 = 1

	// Info, Warnings and Errors
	LevelInfo int32 = 2

	// Warning and Errors
	LevelWarn int32 = 4

	// Errors
	LevelError int32 = 8
)

// GoLogStruct provides support to write to log files.
type GoLogStruct struct {
	LogLevel           int32
	MyTrace              *log.Logger
	MyInfo               *log.Logger
	MyWarning            *log.Logger
	MyError              *log.Logger
	File               *log.Logger
	LogFile            *os.File
}

// log maintains a pointer to a singleton for the logging system.
var logger GoLogStruct

// Called to init the logging system.
func (lS *GoLogStruct) Init(logLevel int32, baseFilePath string) error {
	log.SetPrefix("TRACE: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	err := startFile(logLevel, baseFilePath)
	if err != nil {
		return err;
	}
	lS = &logger
	return err
}

func GetLoggerInstance() (GoLogStruct, error) {
	if &logger == nil {
		return logger, errors.New("Logger not initialized")
	}
	return logger, nil
}

// StartFile initializes GoLogStruct and only displays the specified logging level
// and creates a file to capture writes.
func startFile(logLevel int32, baseFilePath string) error {
	baseFilePath = strings.TrimRight(baseFilePath, "/")
	currentDate := time.Now().UTC()
	dateDirectory := time.Now().UTC().Format("2006-01-02")
	dateFile := currentDate.Format("2006-01-02T15-04-05")

	filePath := fmt.Sprintf("%s/%s/", baseFilePath, dateDirectory)
	fileName := strings.Replace(fmt.Sprintf("%s.txt", dateFile), " ", "-", -1)

	err := os.MkdirAll(filePath, os.ModePerm)
	if err != nil {
		log.Fatalf("main : Start : Failed to Create log directory : %s : %s\n", filePath, err)
		return err
	}

	logf, err := os.Create(fmt.Sprintf("%s%s", filePath, fileName))
	if err != nil {
		log.Fatalf("main : Start : Failed to Create log file : %s : %s\n", fileName, err)
		return err
	}

	
	turnOnLogging(logLevel, logf)
	return err
	
}

// Stop will release resources and shutdown all processing.
func (lS *GoLogStruct) Stop() error {
	var err error
	if lS.LogFile != nil {
		//Trace("main", "Stop", "Closing File")
		err = lS.LogFile.Close()
	}
	return err
}


// LogLevel returns the configured logging level.
func (lS *GoLogStruct) GetLogLevel() int32 {
	return atomic.LoadInt32(&lS.LogLevel)
}

// turnOnLogging configures the logging writers.
func turnOnLogging(logLevel int32, fileHandle io.Writer) {
	traceHandle := ioutil.Discard
	infoHandle := ioutil.Discard
	warnHandle := ioutil.Discard
	errorHandle := ioutil.Discard

	if logLevel&LevelTrace != 0 {
		traceHandle = os.Stdout
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel&LevelInfo != 0 {
		infoHandle = os.Stdout
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel&LevelWarn != 0 {
		warnHandle = os.Stdout
		errorHandle = os.Stderr
	}

	if logLevel&LevelError != 0 {
		errorHandle = os.Stderr
	}

	if fileHandle != nil {
		if traceHandle == os.Stdout {
			traceHandle = io.MultiWriter(fileHandle, traceHandle)
		}

		if infoHandle == os.Stdout {
			infoHandle = io.MultiWriter(fileHandle, infoHandle)
		}

		if warnHandle == os.Stdout {
			warnHandle = io.MultiWriter(fileHandle, warnHandle)
		}

		if errorHandle == os.Stderr {
			errorHandle = io.MultiWriter(fileHandle, errorHandle)
		}
	}

	logger.MyTrace = log.New(traceHandle, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.MyInfo = log.New(infoHandle, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.MyWarning = log.New(warnHandle, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.MyError = log.New(errorHandle, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)

	atomic.StoreInt32(&logger.LogLevel, logLevel)
}



//** TRACE

// Trace writes to the Trace destination
func (lS *GoLogStruct) Trace(format string, a ...interface{}) {
	lS.MyTrace.Output(2, fmt.Sprintf("%s\n", fmt.Sprintf(format, a...)))
}

//** INFO

// Info writes to the Info destination
func (lS *GoLogStruct) Info(format string, a ...interface{}) {
	lS.MyInfo.Output(2, fmt.Sprintf(fmt.Sprintf(format, a...)))
}

//** WARNING

// Warning writes to the Warning destination
func (lS *GoLogStruct) Warning(format string, a ...interface{}) {
	lS.MyWarning.Output(2, fmt.Sprintf(fmt.Sprintf(format, a...)))
}

//** ERROR

// Error writes to the Error destination and accepts an err
func (lS *GoLogStruct) Error(format string, a ...interface{}) {
	lS.MyError.Output(2, fmt.Sprintf(fmt.Sprintf(format, a...)))
}

//writes to the Error and exit(1)
func (lS *GoLogStruct) Fatal(format string, a ...interface{}) {
	lS.MyError.Output(2, fmt.Sprintf(fmt.Sprintf(format, a...)))
	os.Exit(1)
}

