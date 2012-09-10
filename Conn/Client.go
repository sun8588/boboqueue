package Conn

import (
	"Utils"
	"net"
//	"strconv"
	"strings"
)

type CConn struct {
	conn net.Conn
//	IConn
}

/**
初始化conn
**/
func NewClient(conn net.Conn) *CConn {
	return &CConn{conn: conn}
}

func (cc *CConn) Read(num int) ([]byte, error) {
	data, err := Utils.ReadConn(num, cc.conn)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (cc *CConn) Write(data []byte) (int, error) {
	return cc.conn.Write(data)
}
func (cc *CConn) ReadStr(num int) (string, error) {
	data, err := Utils.ReadConn(num, cc.conn)
	if err != nil {
		return "", err
	}
//	Utils.LogInfo("read datastr=%v\n",data)
	return strings.Trim(string(data), "\r\n\t "), nil
}
func (cc *CConn) ReadInt(num int) (int, error) {
	data, err := Utils.ReadConn(num, cc.conn)
	if err != nil {
		return 0, err
	}
//	Utils.LogInfo("read data=%#v\n",data)
	return  int(Utils.BytesToUint32(data)),nil
	
}
