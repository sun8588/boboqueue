package main

import (
"net"
"Utils"
"runtime"
"bufio"
//"io"
)
/**

**/
func main() {
	service := ":843"
	tcpAddr, err := net.ResolveTCPAddr("ip4", service)
	Utils.LogErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	Utils.LogErr(err)
	i := 0
	runtime.GOMAXPROCS(runtime.NumCPU())
	for {
		conn, err := listener.Accept()
		if err!=nil{
			Utils.LogErr(err)
			return 
		}
		Utils.LogInfo("start conn\n")
		go sendPolicy(conn)
//		go Handle.Test(i,roomList, conn)
//		roomChan <- roomList
		i++
	}
}
func sendPolicy(conn net.Conn){
	defer conn.Close()
	str:="<cross-domain-policy><site-control permitted-cross-domain-policies=\"master-only\"/>  <allow-access-from domain=\"*\" to-ports=\"9876\" /></cross-domain-policy>"
//	str+="\0"
	Utils.LogInfo("send data=%s\n",str)
	_,err :=Utils.ReadConn(23,conn)
	if err!=nil{
		Utils.LogErr(err)
	}
	bufk := bufio.NewWriter(conn)
	
	bufk.WriteString(str)
//	bufk.WriteString("\\0")
	bufk.Flush()
}