package Conn

import (
	"Utils"
	"net"
)

type IConn interface {
	Read(int) ([]byte, error)
	Write([]byte) (int, error)
	ReadInt(int) (int, error)
	ReadStr(int) (string, error)
}
/**
根据约定的数据包结构，区分出是telnet还是client
**/
func New(conn net.Conn) (IConn, error) {
	data, err := Utils.ReadConn(1, conn)
	Utils.LogInfo("new data=%v\n",string(data))
	if err != nil {
		return nil, err
	}
	var retn IConn
	if string(data) == "1" {
	Utils.LogInfo("telnet\n")
		retn = NewTelnet(conn)
	} else {
	Utils.LogInfo("client\n")
		retn = NewClient(conn)
	}
	return retn, nil
}
