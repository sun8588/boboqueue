package Conn

import (
"net"
"Utils"
"strconv"
"strings"
)

type TConn struct{
	conn net.Conn
//	IConn
}
/**
初始化一个telnet的连接处理对象
**/
func NewTelnet(conn net.Conn)*TConn{
	return &TConn{conn:conn}
}
/**
从conn中读取数据
**/
func (tc *TConn)Read(num int)([]byte,error){
	return parseCommandWithTelnet(num,tc.conn)
}
func (tc *TConn)ReadInt(num int)(int,error){
	data,err:=tc.Read(num)
	if err!=nil{
		return 0,err
	}
	retn,err:=strconv.Atoi(strings.Trim(string(data),"\r\n\t "))
	if err!=nil{
		return 0,Utils.LogErr(err)
	}
	return  retn,nil
}
func (tc *TConn)ReadStr(num int)(string,error){
	data,err:=tc.Read(num)
	if err!=nil{
		return "",err
	}
	return strings.Trim(string(data),"\r\n\t "),nil
}
/**
输出数据
**/
func (tc *TConn)Write(data []byte)(int,error){
	return tc.conn.Write(data)
}

/**
解析从telnet发送过来的命令，每次调用，解析一个参数，以" "为分割
**/
func parseCommandWithTelnet(num int,conn net.Conn) ([]byte, error) {
	var buf [1]byte
	var data []byte
	for {
		_, err := conn.Read(buf[:])
		if err != nil {
			return nil, Utils.LogErr(err)
		}
//		Utils.LogInfo("data=%siii\n",string(buf[:]))
		if string(buf[:]) == " " || string(buf[:])=="\n"{
			break
		}
		data = append(data, buf[:]...)
		num--
	}
	return data, nil

}