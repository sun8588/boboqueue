package Handle

import (
	"Entity"
	"Utils"
	//		"bufio"
	"bytes"
	//	"encoding/binary"
	//	"io"
	"log"
	//	"container/list"
	//	"bufio"
	"net"
	//	"strings"
	//"Handle"
	//	"os"
	//	"error"
//	"strconv"

	//	"encoding/json" 
	//"sync"
	"Conn"
)

func HandleClient(readyData *Entity.ReadyData, waitData *Entity.WaitData, conn net.Conn) {
	
	handleConn,cerr:= Conn.New(conn)
	if cerr != nil {
		return
	}
	for {
		cmd, err := handleConn.ReadInt(4)
		if err != nil {
			continue
		}
		Utils.LogInfo("command is=%#v\n", cmd)
//		cmd:=handleConn.ReadInt(command)
//		cmd, err := strconv.Atoi(string(command))
//		return 
		//		command, err := parseCommandWithTelnet(conn, 4)
//		if err != nil {
//			conn.Write([]byte(err.Error() + "\n"))
//			continue
//		}

		switch cmd {
		//添加task
		case 100:
			Utils.LogInfo("hande\n")
			keyLen, err := handleConn.ReadInt(4)
			Utils.LogInfo("get key=%v\n", keyLen)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			key, err := handleConn.ReadStr(keyLen)
			
			Utils.LogInfo("get key=%v\n", key)
			valueLen, err := handleConn.ReadInt(4)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			value, err := handleConn.ReadStr(valueLen)
			Utils.LogInfo("get value=%v\n", value)
			expired, err := handleConn.ReadInt(4)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			Utils.LogInfo("args=%s,%s,%d\n", key, value, expired)
			return 
			err = Add(readyData, waitData, key, value, uint(expired))
			if err != nil {
				handleConn.Write([]byte(err.Error() + "\n"))
				continue
			} else {
				handleConn.Write([]byte("done\n"))
			}

			//得到数据
		case 101:
			num, err := handleConn.ReadInt(4)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			Utils.LogInfo("num=%d\n", num)
			return 
			data, err := Get(readyData, waitData, num)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
			}
			handleConn.Write(data)
			handleConn.Write([]byte("done\n"))
			//删除数据
		case 102:
			keyLen, err := handleConn.ReadInt(4)
			Utils.LogInfo("get keylen=%v\n", keyLen)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			key, err := handleConn.ReadStr(keyLen)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			Utils.LogInfo("args=%s\n", key)
			return
			err = Del(readyData, waitData, key)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			handleConn.Write([]byte("done\n"))
			//增加到期时间
		case 103:
			keyLen, err := handleConn.ReadInt(4)
			Utils.LogInfo("get keylen=%v\n", keyLen)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			key, err := handleConn.ReadStr(keyLen)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			expired, err := handleConn.ReadInt(4)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			Utils.LogInfo("key=%s,expired=%d\n",key, expired)
			return 
			err = AddExpired(readyData, waitData, key, uint(expired))
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			handleConn.Write([]byte("done\n"))
			//减少到期时间
		case 104:
			keyLen, err := handleConn.ReadInt(4)
			Utils.LogInfo("get keylen=%v\n", keyLen)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			key, err := handleConn.ReadStr(keyLen)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			expired, err := handleConn.ReadInt(4)
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			Utils.LogInfo("key=%s,expired=%d\n",key, expired)
			return
			err = DecExpired(readyData, waitData, key, uint(expired))
			if err != nil {
				conn.Write([]byte(err.Error() + "\n"))
				continue
			}
			handleConn.Write([]byte("done\n"))
		default:
			Utils.LogInfo("command not found....")
			handleConn.Write([]byte("command not found....\n"))
			continue
		}

		for key, value := range *readyData {
			log.Printf("ready key=%s, data=%v\n", key, value)
		}
		for key, value := range *waitData {
			log.Printf("wait key=%s, data=%v\n", key, value)
		}
	}
}
func parseCommandWithClient(conn net.Conn) (int, string, interface{}, int, error) {
	//得到包长度
	dataLenBuf, err := Utils.ReadConn(4, conn)
	if err != nil {
		Utils.LogErr(err)
		return 0, "", nil, 0, err
	}
	dataLen := Utils.BytesToUint32(dataLenBuf)
	Utils.LogInfo("dataLen=%v\n", dataLen)
	if dataLen == 0 {
		return 0, "", nil, 0, err
	}
	//得到整个包体
	data, err := Utils.ReadConn(int(dataLen), conn)
	if err != nil {
		Utils.LogErr(err)
		return 0, "", nil, 0, err
	}
	//生成bytestream
	dataBuf := bytes.NewBuffer(data)
	log.Printf("alldata=%v", data)
	//取出命令
	command := Utils.BytesToUint32(dataBuf.Next(4))
	log.Printf("get command=%d", command)
	return 0, "", nil, 0, nil
}
