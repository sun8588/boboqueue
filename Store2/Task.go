package Store2

import (
	"Entity"

	"Model"
	"Utils"
	"bytes"
	"encoding/gob"
	"os"
	"time"
)

/**
数据model对象
**/
var DataStore Model.IModel

/**
数据是否变化，用于dump方式存储数据
**/
var dataChange bool

/**
用于方便存数据到文件中
**/
type taskData struct {
	ReadyData *Entity.ReadyData
	WaitData  *Entity.WaitData
}

var TaskData taskData

/**
初始化index和data的model对象
**/
func New(dataModel Model.IModel) {
	DataStore = dataModel
}

/**
取出readyData和waitData数据
**/
func GetData() (*Entity.ReadyData, *Entity.WaitData, error) {
	//得到存储方式的句柄
	fd := DataStore.GetFd().(*os.File)
	//得到文件状态
	fileStat, err := fd.Stat()
	if err != nil {
		return nil, nil, Utils.LogErr(err)
	}
	//如果文件是空的，直接返回
	fileSize:=fileStat.Size()
	if fileSize==0{
		return &Entity.ReadyData{},&Entity.WaitData{},nil
	}
	//得到文件大小的buf
	dataBuf := make([]byte,fileSize )
	//得到文件读取buf
	bufReader := DataStore.Get()
	_, err = bufReader.Read(dataBuf)
	Utils.LogInfo("read file data=%#v\n", dataBuf)
	if err != nil {
		Utils.LogInfo("err=%d\n", err)
		return nil, nil, err
	}
	//吧读到的[]byte转换成结构数据,使用全局变量TaskData
	rgob := gob.NewDecoder(bytes.NewBuffer(dataBuf))
	err = rgob.Decode(&TaskData)
	Utils.LogInfo("gob data =%v\n", TaskData)
	if err != nil {
		Utils.LogInfo("gob data err=%d\n", err)
		return nil, nil, Utils.LogErr(err)
	}
	return TaskData.ReadyData, TaskData.WaitData, nil
}

/**
写入数据
**/
func SetData(readyData *Entity.ReadyData, waitData *Entity.WaitData) error {
	TaskData.ReadyData = readyData
	TaskData.WaitData = waitData
	dataChange=true
	return nil
}

/**
定时tick，把数据写入文件
**/
func EntityDump(tickTime time.Duration) {
	c := time.Tick(tickTime*time.Second)
	for _= range c {
//		Utils.LogInfo("%v\n", now)
		//如果在规定时间内没有更改，就不更新
		if dataChange==false{
			continue
		}
		//得到存储方式的句柄
		fd := DataStore.GetFd().(*os.File)
		//清空文件
		fd.Truncate(0)
		//得到gob转换出来的[]byte
		var wbuf bytes.Buffer
		egob := gob.NewEncoder(&wbuf)
		Utils.LogInfo("gob data=%#v\n", TaskData)
		err := egob.Encode(TaskData)
		if err != nil {
			Utils.LogInfo("data gob err")
			Utils.LogErr(err)
			continue
		}
		_, err = DataStore.Set(wbuf.Bytes())
		if err != nil {
			Utils.LogErr(err)
			continue
		}
		dataChange=false
	}
}
