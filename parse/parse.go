// Package collectdparse is to parse collectd binary protocol.
// Spec is here https://collectd.org/wiki/index.php/Binary_protocol


package parse

import (
  "bytes"
  "encoding/binary"
  "fmt"
  "github.com/nateseay/collectd-listener/logging"
  "time"
)


type PartType uint16

//go:generate stringer -type=PartType
const (
  Host PartType = 0x0000
  Time PartType = 0x0001
  TimeHighRes PartType = 0x0008
  Plugin PartType = 0x0002
  PluginInstance PartType = 0x0003
  CPartType PartType = 0x0004
  CPartTypeInstance PartType = 0x0005
  Values PartType = 0x0006
  Interval PartType = 0x0007
  IntervalHighRes PartType = 0x0009
  Message PartType = 0x0100
  Severity PartType = 0x0101
  Signature PartType = 0x0200 //TODO in the spec
  Encryption PartType = 0x0210 // TODO in the spec
)

type DataType uint8
const (
  COUNTER DataType = 0
  GAUGE DataType = 1
  DERIVE DataType = 2
  ABSOLUTE DataType = 3
)

//http://stackoverflow.com/questions/28012952/golang-variable-length-array-in-struct-for-use-with-binary-read
type BinaryProtocolHeader struct {
  BinaryType PartType
  DataLength uint16

}

type BinaryValuesHeade struct {


}

func timeDecode( logger logging.LoggerStruct,
                  IsHighRes bool,
                  value uint64,
                  TimeZoneLoc time.Location) string {
  var secs int64
  var usecs uint64
  secs = 0
  usecs = 0
  if IsHighRes{
    usecs = value >> 30
  } else{
    usecs = value
  }
  secs = int64(usecs)
  utc := time.Unix(secs, 0)
  local := utc.In(&TimeZoneLoc)
  ret := fmt.Sprintf("Epoch: %v Timestamp UTC: %v Timestamp Local: %v", secs, utc, local)
  return ret
}

func parsePart( logger logging.LoggerStruct,
                part []byte,
                header BinaryProtocolHeader,
                TimeZoneLoc time.Location){
  switch header.BinaryType {
    // These are the string types
    case Host, Plugin, PluginInstance, CPartType, CPartTypeInstance, Message:
      PartString := string(part)
      logger.Info.Println(fmt.Sprintf("Decoded Part: type=%v value=%v", header.BinaryType, PartString))
    // These are the number seconds since epoch values
  case Time, Interval, TimeHighRes, IntervalHighRes:
      Epoch := binary.BigEndian.Uint64(part)
      var HighRes bool
      HighRes = false
      if (header.BinaryType == TimeHighRes) || (header.BinaryType == IntervalHighRes){
        HighRes = true
      }
      TimeString := timeDecode(logger, HighRes, Epoch, TimeZoneLoc)
      logger.Info.Println(fmt.Sprintf("Decoded Part: type=%v raw_value=%v value/s=%v", TimeString, Epoch, TimeString ))
    default:
      logger.Info.Println(fmt.Sprintf("Decoded Part: type %v not implemented yet!", header.BinaryType))
  }
 }

func ParseBuffer ( logger logging.LoggerStruct,
                    buf []byte ,
                    TimeZoneLoc time.Location) {
  logger.Trace.Println("In ParseBuffer")
  // create a reader
  reader := bytes.NewReader(buf)
  logger.Trace.Println(reader.Len())
  for reader.Len() > 0 {
    logger.Trace.Println(fmt.Sprintf("Moving on to next part, %v remaining", reader.Len()))
    var header BinaryProtocolHeader
    err := binary.Read(reader, binary.BigEndian,  &header)
    logger.Trace.Println(fmt.Sprintf("Decoded Header: err=%v Type=%v DataLength=%v", err, header.BinaryType, header.DataLength))
    part := make([]byte, (header.DataLength-4)) // -4 is to account for the header which has been read, but is included in DataLength
    err2 := binary.Read(reader, binary.BigEndian, part)
    logger.Trace.Println(fmt.Sprintf("Read Part err=%v buffer length remaining=%v", err2, reader.Len()))
    parsePart(logger, part, header, TimeZoneLoc)
  }
}
