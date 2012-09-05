package main

import (
	"Config"
	"Store2"
	"os"
	//	"Entity"
	"Handle"
	"Utils"
	"log"
	"net"
	"os/signal"
	"runtime"
	"syscall"
)

/**

**/
func main() {

	service := ":9800"
	tcpAddr, err := net.ResolveTCPAddr("ip4", service)
	Utils.LogErr(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	Utils.LogErr(err)
	i := 0
	readyData, waitData, err := Handle.LoadTask()
	if err != nil {
		Utils.LogErr(err)
		return
	}
	log.Println("load data=%v", readyData)
	for key, value := range *readyData {
		log.Printf("load ready key=%s, data=%v\n", key, value)
	}
	for key, value := range *waitData {
		log.Printf("load wait key=%s, data=%v\n", key, value)
	}

	//	taskList := make([]Entity.Task, Config.GetRoomMaxNum())
	runtime.GOMAXPROCS(runtime.NumCPU())
	//go 程处理数据写入文件
	go Store2.EntityDump(Config.GetTickTime())
	//go程处理waittask检索
	go Handle.CheckWaitTask(readyData, waitData)
	go signalHandle()
	for {
		conn, err := listener.Accept()
		if err != nil {
			Utils.LogErr(err)
			return
		}
		go Handle.HandleClient(readyData, waitData, conn)
		//		go Handle.Test(i,roomList, conn)
		//		roomChan <- roomList
		i++
	}
}

/**
信号量处理函数
**/
func signalHandle() {
	for {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGUSR1)
		sig := <-ch
		Utils.LogInfo("Signal received: %v", sig)
		switch sig {
		case syscall.SIGINT:
			os.Exit(1)
		case syscall.SIGUSR1:
			Config.ParseXml(Config.ConfigFile)
		}
	}
}
