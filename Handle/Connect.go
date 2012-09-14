package Handle

import (
	"Entity"
	"Utils"
	"log"

	"net"

	"Conn"
	"io"
)

func HandleClient(readyData *Entity.ReadyData, waitData *Entity.WaitData, conn net.Conn) {
	defer func() {
		conn.Close()
	}()
	handleConn, cerr := Conn.New(conn)
	if cerr != nil {
		return
	}
	for {
		cmd, err := handleConn.ReadInt(4)
		if err != nil {
			if err == io.EOF {
				break
			}
			continue
		}
		if cmd == 0 {
			break
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
				handleConn.Write(errorRet(err))
				//				conn.Write( + "\n"))
				continue
			}
			key, err := handleConn.ReadStr(keyLen)

			Utils.LogInfo("get key=%v\n", key)
			valueLen, err := handleConn.ReadInt(4)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			value, err := handleConn.ReadStr(valueLen)
			Utils.LogInfo("get value=%v\n", value)
			expired, err := handleConn.ReadInt(4)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			Utils.LogInfo("args=%s,%s,%d\n", key, value, expired)
			err = Add(readyData, waitData, key, value, uint(expired))
			if err != nil {
				Utils.LogInfo("err\n")
				//				return 
				handleConn.Write(errorRet(err))
				continue
			} else {
				Utils.LogInfo("done\n")
				//			return
				handleConn.Write(doneRet([]byte("1")))
				continue
			}

			//得到数据
		case 101:
			num, err := handleConn.ReadInt(4)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			Utils.LogInfo("num=%d\n", num)
			return
			data, err := Get(readyData, waitData, num)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			handleConn.Write(doneRet(data))
			//删除数据
		case 102:
			keyLen, err := handleConn.ReadInt(4)
			Utils.LogInfo("get keylen=%v\n", keyLen)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			key, err := handleConn.ReadStr(keyLen)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			Utils.LogInfo("args=%s\n", key)
			err = Del(readyData, waitData, key)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			handleConn.Write(doneRet([]byte("1")))
			//增加到期时间
		case 103:
			keyLen, err := handleConn.ReadInt(4)
			Utils.LogInfo("get keylen=%v\n", keyLen)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			key, err := handleConn.ReadStr(keyLen)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			expired, err := handleConn.ReadInt(4)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			Utils.LogInfo("key=%s,expired=%d\n", key, expired)
			err = AddExpired(readyData, waitData, key, uint(expired))
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			handleConn.Write(doneRet([]byte("1")))
			//减少到期时间
		case 104:
			keyLen, err := handleConn.ReadInt(4)
			Utils.LogInfo("get keylen=%v\n", keyLen)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			key, err := handleConn.ReadStr(keyLen)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			expired, err := handleConn.ReadInt(4)
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			Utils.LogInfo("key=%s,expired=%d\n", key, expired)
			err = DecExpired(readyData, waitData, key, uint(expired))
			if err != nil {
				handleConn.Write(errorRet(err))
				continue
			}
			handleConn.Write(doneRet([]byte("1")))
		default:
			Utils.LogInfo("command not found....")
			handleConn.Write(errorRet(Utils.LogErr(100)))
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

/**
组合错误结果，返回给client
**/
func errorRet(err error) []byte {
	data := []byte("1" + err.Error())
	dataLenByte := Utils.Uint32ToBytes(uint32(len(data)))
	return append(dataLenByte[:], data...)
}

/**
完成结果，返回给client
**/
func doneRet(data []byte) []byte {
	dataRet := append([]byte("0"), data...)
	dataLenByte := Utils.Uint32ToBytes(uint32(len(dataRet)))
	return append(dataLenByte[:], dataRet...)
}
