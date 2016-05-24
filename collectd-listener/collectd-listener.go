package main

import (
  "os"
  "fmt"
  "net"
  "flag"
  "github.com/nateseay/collectd-listener/parse"
  "github.com/nateseay/collectd-listener/logging"
  "bufio"
  "time"
)

type args_struct struct {
  udp_port *int
  typesdb_path *string
  log_level *int
  timezone_file *string
}

func parse_args() args_struct{
  var args args_struct
  args.udp_port = flag.Int("udp-port", 8096, "UDP port to listen for collectd data" )
  args.typesdb_path = flag.String("typesdb-path", "./typedb", "Path to the collectd typesdb file")
  args.log_level = flag.Int("log-level", 2, "Log level, 1-4. 1 is least verbose, 4 is most")
  args.timezone_file = flag.String("timezone", "/etc/timezone", "Timezone file location. Will default to America/LosAngeles if it can't be loaded")
  flag.Parse()
  return args
}

func check_error(logger logging.LoggerStruct, err error){
  if err != nil {
    logger.Error.Println("Error! ", err)
    os.Exit(1)
  }
}

func loadTimezone(logger logging.LoggerStruct,
                  args args_struct) time.Location{
  logger.Trace.Println(fmt.Sprintf("Trying to load timezone file from %v", *args.timezone_file))
  DefTZ := "America/LosAngeles"
  file, err := os.Open(*args.timezone_file)
  if err != nil {
    logger.Error.Println(fmt.Sprintf("Unable to open timezone file. Falling back to %v. Err=%v", DefTZ, err))
    DefaultLoc, _ := time.LoadLocation(DefTZ)
    return *DefaultLoc
  }
  defer file.Close()
  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan(){
    lines = append(lines, scanner.Text())
  }
  SetLoc, _ := time.LoadLocation(lines[0])
  logger.Info.Println(fmt.Sprintf("Timezone set to %v", SetLoc))
  return *SetLoc
}

func run_server(logger logging.LoggerStruct,
                args args_struct) {
  /* Load timezone info */
  TimeZoneLoc := loadTimezone(logger, args)
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
      logger.Info.Println(fmt.Sprintf("Received data from %v, length of %v", addr, n))
      if err != nil {
          logger.Error.Println("Error: ",err)
      }
      parse.ParseBuffer(logger, buf, TimeZoneLoc)
  }

}

func main(){
  args := parse_args()
  logger := logging.LogInit(*args.log_level)
  logger.Info.Println("[Config] udp-port:", *args.udp_port)
  logger.Info.Println("[Config] typedb_path:", *args.typesdb_path)
  logger.Info.Println("[Config] log_level:", *args.log_level)
  logger.Info.Println("[Config] timezone_file:", *args.timezone_file)
  run_server(logger, args)
}
