package Store

import (
	"Entity"

	"Model"
	"Utils"
	"bytes"
	"encoding/gob"
	"io"
	"os"
)

/**
索引数据model对象
**/
var IndexStore Model.IModel
/**
数据model对象
**/
var DataStore Model.IModel

/**
初始化index和data的model对象
**/
func New(indexModel, dataModel Model.IModel) {
	IndexStore = indexModel
	DataStore = dataModel
}

/**
写入任意数据，都使用gob存入file,返回的int64，是该数据在文件中的起始位置
**/
func AddData(data Entity.TaskData) (int, error) {
	//得到存储方式的句柄
	fd := DataStore.GetFd().(*os.File)
	//得到文件状态
	fileStat, err := fd.Stat()
	if err != nil {
		return 0, Utils.LogErr(err)
	}
	data.FileAddr = fileStat.Size()
	//得到gob转换出来的[]byte
	var wbuf bytes.Buffer
	egob := gob.NewEncoder(&wbuf)
	Utils.LogInfo("gob data=%#v\n", data)
	err = egob.Encode(data)
	if err != nil {
		Utils.LogInfo("data gob err")
		return 0, Utils.LogErr(err)
	}
	//拼出最终存入文件的数据==数据长度+数据体
	var storeByte []byte
	storeByteBuf := bytes.NewBuffer(storeByte)
	//写入数据长度
	_, err = storeByteBuf.Write(Utils.Uint32ToBytes(uint32(wbuf.Len())))
	if err != nil {
		return 0, Utils.LogErr(err)
	}
	_, err = storeByteBuf.Write(wbuf.Bytes())
	if err != nil {
		return 0, Utils.LogErr(err)
	}
	wData := storeByteBuf.Bytes()
	Utils.LogInfo("len=%d\n", storeByteBuf.Len())
	Utils.LogInfo("data=%#v\n", wData)
	writeNum, err := DataStore.Set(wData)
	if err != nil {
		return 0, Utils.LogErr(err)
	}
	return writeNum, nil
}

/**
写入索引
**/
func AddIndex(dataAddr Entity.TaskAddr) (int, error) {
	//得到存储方式的句柄
	fd := IndexStore.GetFd().(*os.File)
	//得到文件状态
	fileStat, err := fd.Stat()
	if err != nil {
		return 0, Utils.LogErr(err)
	}
	//	//得到文件大小，这个位置就是这条记录的起始位置
	dataAddr.FileAddr = fileStat.Size()
	dataAddr.EleAddr = nil
	gob.Register(dataAddr)
	//得到gob转换出来的[]byte
	var wbuf bytes.Buffer
	egob := gob.NewEncoder(&wbuf)
	Utils.LogInfo("gob index=%#v\n", dataAddr)
	err = egob.Encode(dataAddr)
	if err != nil {
		Utils.LogInfo("index gob err")
		return 0, Utils.LogErr(err)
	}
	//拼出最终存入文件的数据==数据长度+数据体
	var storeByte []byte
	storeByteBuf := bytes.NewBuffer(storeByte)
	//写入数据长度
	_, err = storeByteBuf.Write(Utils.Uint32ToBytes(uint32(wbuf.Len())))
	if err != nil {
		return 0, Utils.LogErr(err)
	}
	_, err = storeByteBuf.Write(wbuf.Bytes())
	if err != nil {
		return 0, Utils.LogErr(err)
	}
	wData := storeByteBuf.Bytes()
	Utils.LogInfo("len=%d\n", storeByteBuf.Len())
	Utils.LogInfo("data=%#v\n", wData)
	writeNum, err := IndexStore.Set(wData)
	if err != nil {
		return 0, Utils.LogErr(err)
	}
	return writeNum, nil
}

/**
返回所有的任务数据数组
**/
func GetAllTaskData() ([]Entity.TaskData, error) {
	bufReader := DataStore.Get()
	var gobTypeList []Entity.TaskData
	for {
		var gobType Entity.TaskData
		dataLenBuf := make([]byte, 4)
		//得到数据体长度
		_, err := bufReader.Read(dataLenBuf)
		if err != nil {
			Utils.LogInfo("err=%d\n", err)
			if err == io.EOF {
				break
			}
			return nil, err
		}
		dataLen := Utils.BytesToUint32(dataLenBuf)
		dataBuf := make([]byte, dataLen)
		_, err = bufReader.Read(dataBuf)
		Utils.LogInfo("read file data=%#v\n", dataBuf)
		if err != nil {
			Utils.LogInfo("err=%d\n", err)
			return nil, err
		}
		rgob := gob.NewDecoder(bytes.NewBuffer(dataBuf))
		err = rgob.Decode(&gobType)
		if err != nil {
			Utils.LogInfo("gob data err=%d\n", err)
			return nil, Utils.LogErr(err)
		}
		gobTypeList = append(gobTypeList[:], gobType)
	}
	Utils.LogInfo("retn data=%#v\n", gobTypeList)
	return gobTypeList, nil
}

/**
返回所有任务索引数组
**/
func GetAllTaskIndex() ([]Entity.TaskAddr, error) {
	bufReader := IndexStore.Get()
	var gobTypeList []Entity.TaskAddr
	for {
		var gobType Entity.TaskAddr
		dataLenBuf := make([]byte, 4)
		//得到数据体长度
		_, err := bufReader.Read(dataLenBuf)
		if err != nil {
			Utils.LogInfo("err=%d\n", err)
			if err == io.EOF {
				break
			}
			return nil, err
		}
		dataLen := Utils.BytesToUint32(dataLenBuf)
		dataBuf := make([]byte, dataLen)
		_, err = bufReader.Read(dataBuf)
		Utils.LogInfo("read file index=%#v\n", dataBuf)
		if err != nil {
			Utils.LogInfo("err=%d\n", err)
			return nil, err
		}
		rgob := gob.NewDecoder(bytes.NewBuffer(dataBuf))
		err = rgob.Decode(&gobType)
		if err != nil {
			Utils.LogInfo("gob index err=%d\n", err)
			return nil, Utils.LogErr(err)
		}
		gobTypeList = append(gobTypeList[:], gobType)
	}
	Utils.LogInfo("retn index=%#v\n", gobTypeList)
	return gobTypeList, nil
}

