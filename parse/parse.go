// Package collectdparse is to parse collectd binary protocol.
// Spec is here https://collectd.org/wiki/index.php/Binary_protocol


package parse

import (
  "bytes"
  "github.com/nateseay/collectd-listener/logging"
)

func ParseBuffer ( logger logging.LoggerStruct,
                    buf []byte ) {
  logger.Trace.Println("In ParseBuffer")
  // create a reader
  reader := bytes.NewReader(buf)
  logger.Trace.Println(reader.Len())
}
