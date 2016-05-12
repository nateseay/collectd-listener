// Package collectdparse is to parse collectd binary protocol.
// Spec is here https://collectd.org/wiki/index.php/Binary_protocol


package parse

import (
  "fmt"
  "bytes"
)

func ParseBuffer ( logger logger_struct,
                    buf []byte ) {
  logger.Trace.Println("In ParseBuffer")
  // create a reader
  reader := bytes.NewReader(byte)
  logger.Trace.Println(reader.Len())
}
