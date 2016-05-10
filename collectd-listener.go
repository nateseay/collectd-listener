package main

import (
  "flag"
  "log"
  "os"
  "io/ioutil"
  "io"
  "fmt"
  "net"
)

type args_struct struct {
  udp_port *int
  typesdb_path *string
  log_level *int
}

type logger_struct struct {
  Trace *log.Logger
  Info *log.Logger
  Warning *log.Logger
  Error *log.Logger
}

func log_init(log_level int) logger_struct{
  var logger logger_struct
  var errorHandle, warningHandle, infoHandle, traceHandle io.Writer
  errorHandle = os.Stderr
  if (log_level >= 2){
    warningHandle = os.Stdout
  } else {
    warningHandle = ioutil.Discard
  }
  if (log_level >= 3){
    infoHandle = os.Stdout
  } else {
    infoHandle = ioutil.Discard
  }
  if (log_level >= 4){
    traceHandle = os.Stdout
  } else {
    traceHandle = ioutil.Discard
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

func parse_args() args_struct{
  var args args_struct
  args.udp_port = flag.Int("udp-port", 8096, "UDP port to listen for collectd data" )
  args.typesdb_path = flag.String("typesdb-path", "./typedb", "Path to the collectd typesdb file")
  args.log_level = flag.Int("log-level", 4, "Log level, 1-4. 1 is least verbose, 4 is most")
  flag.Parse()
  return args
}

func check_error(logger logger_struct, err error){
  if err != nil {
    logger.Error.Println("Error! ", err)
    os.Exit(1)
  }
}

func run_server(logger logger_struct,
                args args_struct) {
  /* Lets prepare a address at any address at port 10001*/
  logger.Info.Println(fmt.Sprintf("Preparing to bind to %v", *args.udp_port))
  ServerAddr,err := net.ResolveUDPAddr("udp",fmt.Sprintf(":%v", *args.udp_port))
  check_error(logger, err)

  /* Now listen at selected port */
  ServerConn, err := net.ListenUDP("udp", ServerAddr)
  check_error(logger, err)
  defer ServerConn.Close()

  logger.Info.Println(fmt.Sprintf("Listening at %v", *args.udp_port))

  buf := make([]byte, 1024)

  for {
      n,addr,err := ServerConn.ReadFromUDP(buf)
      logger.Info.Println("Received ",string(buf[0:n]), " from ",addr)

      if err != nil {
          logger.Error.Println("Error: ",err)
      }
  }

}

func main(){
  args := parse_args()
  logger := log_init(*args.log_level)
  logger.Info.Println("[Config] udp-port:", *args.udp_port)
  logger.Info.Println("[Config] typedb_path:", *args.typesdb_path)
  logger.Info.Println("[Config] log_level:", *args.log_level)
  run_server(logger, args)
}
