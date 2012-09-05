package Utils

import (
	"encoding/json"
	"net"
	"io"
	"log"

)

func DecodeJson(byteArr []byte) (map[string]interface{}, error) {
	var msg interface{}
	err := json.Unmarshal(byteArr, &msg)
	if err != nil {
		return nil, LogErr(err)
	}
	return msg.(map[string]interface{}), nil
}

func EncodeJson(jsonData interface{}) ([]byte, error) {
	msg, err := json.Marshal(jsonData)
	if err != nil {
		LogErr(102)
		return nil, err
	}
	return msg, err
}

// Encode uint32 to [4]byte
func Uint32ToBytes(i uint32) []byte {
	return []byte{byte((i >> 24) & 0xff), byte((i >> 16) & 0xff),
		byte((i >> 8) & 0xff), byte(i & 0xff)}
}

func BytesToUint32(buf []byte) uint32 {
	return uint32(buf[0])<<24 + uint32(buf[1])<<16 + uint32(buf[2])<<8 +
		uint32(buf[3])
}

/**
从net.conn 中读取固定长度
**/
func ReadConn(readLen int, conn net.Conn) ([]byte, error) {
	dataBuf := make([]byte, readLen)
	for {
		var dataLenTag int
		tmpNum, err := conn.Read(dataBuf[dataLenTag:readLen])
		log.Printf("read num=%d\n",tmpNum)
		log.Printf("read data=%s\n",dataBuf)
		if err != nil {
			if err==io.EOF{
			return dataBuf,nil
			}
			//				Utils.LogInfo("read num=%d\n", tmpNum)
			LogErr(err)
			return dataBuf,err
		}
		LogInfo("datalentag=%d\n",dataLenTag)
		dataLenTag += tmpNum
		if dataLenTag >= readLen {
			break
		}
	}
	return dataBuf, nil
}
