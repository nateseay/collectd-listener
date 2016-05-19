
package logging

import (
  "log"
  "os"
  "io/ioutil"
  "io"
)

type LoggerStruct struct {
  Trace *log.Logger
  Info *log.Logger
  Warning *log.Logger
  Error *log.Logger
}

func LogInit(log_level int) LoggerStruct{
  var logger LoggerStruct
  var errorHandle, warningHandle, infoHandle, traceHandle io.Writer
  errorHandle = os.Stderr

  switch log_level {
  case 2:
    warningHandle = os.Stdout
    infoHandle = ioutil.Discard
    traceHandle = ioutil.Discard
  case 3:
    warningHandle = os.Stdout
    infoHandle = os.Stdout
    traceHandle = ioutil.Discard
  default:
    // fall back to all logging on
    warningHandle = os.Stdout
    infoHandle = os.Stdout
    traceHandle = os.Stdout
  }

  logger.Trace = log.New(traceHandle,
      "TRACE: ",
      log.Ldate|log.Ltime|log.Lshortfile)

  logger.Info = log.New(infoHandle,
      "INFO: ",
      log.Ldate|log.Ltime|log.Lshortfile)

  logger.Warning = log.New(warningHandle,
      "WARNING: ",
      log.Ldate|log.Ltime|log.Lshortfile)

  logger.Error = log.New(errorHandle,
      "ERROR: ",
      log.Ldate|log.Ltime|log.Lshortfile)

  logger.Error.Println("Logger initialized")
  logger.Warning.Println("Logger initialized")
  logger.Info.Println("Logger initialized")
  logger.Trace.Println("Logger initialized")
  return logger
}
