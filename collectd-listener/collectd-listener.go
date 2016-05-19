package main

import (
  "os"
  "fmt"
  "net"
  "flag"
  "github.com/nateseay/collectd-listener/parse"
  "github.com/nateseay/collectd-listener/logging"
)

type args_struct struct {
  udp_port *int
  typesdb_path *string
  log_level *int
}

func parse_args() args_struct{
  var args args_struct
  args.udp_port = flag.Int("udp-port", 8096, "UDP port to listen for collectd data" )
  args.typesdb_path = flag.String("typesdb-path", "./typedb", "Path to the collectd typesdb file")
  args.log_level = flag.Int("log-level", 2, "Log level, 1-4. 1 is least verbose, 4 is most")
  flag.Parse()
  return args
}

func check_error(logger logging.LoggerStruct, err error){
  if err != nil {
    logger.Error.Println("Error! ", err)
    os.Exit(1)
  }
}

func run_server(logger logging.LoggerStruct,
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
      logger.Info.Println(fmt.Sprintf("Received data from %v, length of %v", addr, n))
      if err != nil {
          logger.Error.Println("Error: ",err)
      }
      parse.ParseBuffer(logger, buf)
  }

}

func main(){
  args := parse_args()
  logger := logging.LogInit(*args.log_level)
  logger.Info.Println("[Config] udp-port:", *args.udp_port)
  logger.Info.Println("[Config] typedb_path:", *args.typesdb_path)
  logger.Info.Println("[Config] log_level:", *args.log_level)
  run_server(logger, args)
}
