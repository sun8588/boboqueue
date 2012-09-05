package Handle

//_,import "Entity"
import (
	"Entity"
	"Model"
	Store "Store2"
	"Utils"
	//	"container/list"
	"net"
	"time"

//	"strconv"
)

/**
从文件中加载数据
**/
func LoadTask() (*Entity.ReadyData, *Entity.WaitData, error) {
	file, err := Model.New("/home/dingbo/workspace/message/src/data.log")
	if err != nil {
		return nil, nil, err
	}
	Store.New(file)
	readyData, waitData, err := Store.GetData()
	if err != nil {
		return nil, nil, err
	}
	return readyData, waitData, nil
}

/**
检查进入等待状态的任务是否到期，到期后放入readyData
**/
func CheckWaitTask(readyData *Entity.ReadyData, waitData *Entity.WaitData) {
	c := time.Tick(time.Second)
	dataChange := false
	for _ = range c {
		for key, value := range *waitData {
			//已经到期的任务，放入readyData
			if value.Expired <= uint(time.Now().Unix()) {
				Utils.LogInfo("check do=%s\n", key)
				readyData.Add(key, value.Value, value.Expired)
				waitData.Del(key)
				dataChange = true
			}
		}
		if dataChange {
			//存储数据
			Store.SetData(readyData, waitData)
			dataChange = false
		}

	}

}

/**
添加数据
**/
func Add(readyData *Entity.ReadyData, waitData *Entity.WaitData, key string, value interface{}, expired uint) error {
	//如果到期时间为0,表示是一个需要立即执行的任务，直接放入readyData列表中
	//不为0,放入waitData，等待时间到期
	if expired == 0 {
		readyData.Add(key, value, expired)
	} else {
		waitData.Add(key, value, expired)
	}
	//存储数据
	Store.SetData(readyData, waitData)
	return nil
}

/**
返回对应num的任务json
**/
func Get(readyData *Entity.ReadyData, waitData *Entity.WaitData, num int) ([]byte, error) {
	data := readyData.Cut(num)
	jsonByte, err := Utils.EncodeJson(data)
	if err != nil {
		return nil, err
	}
	//存储数据
	Store.SetData(readyData, waitData)
	return jsonByte, nil
}

/**
删除一个任务，在两个列表里删，ready和wait
**/
func Del(readyData *Entity.ReadyData, waitData *Entity.WaitData, key string) error {
	if readyData.Isset(key) {
		readyData.Del(key)
		//存储数据
		Store.SetData(readyData, waitData)
		return nil
	}
	if waitData.Isset(key) {
		waitData.Del(key)
		//存储数据
		Store.SetData(readyData, waitData)
		return nil
	}
	return Utils.LogErr(1000)
}

/**
增加任务到期时间
**/
func AddExpired(readyData *Entity.ReadyData, waitData *Entity.WaitData, key string, expired uint) error {
	//增加到期时间，就从就绪队列中删除，并添加到wait队列中
	if readyData.Isset(key) {
		waitData.Add(key, (*readyData)[key].Value, (*readyData)[key].Expired)
		readyData.Del(key)
		//存储数据
		Store.SetData(readyData, waitData)
		return nil
	}
	//在wait队列中的，只更新到期时间
	if waitData.Isset(key) {
		waitData.AddExpired(key, expired)
		//存储数据
		Store.SetData(readyData, waitData)
		return nil
	}
	return Utils.LogErr(1001)
}

/**
减少时间
**/
func DecExpired(readyData *Entity.ReadyData, waitData *Entity.WaitData, key string, expired uint) error {
	//如果任务已经在ready中，就不需要处理
	if readyData.Isset(key) {
		return nil
	}
	//在wait队列中的，只更新到期时间
	if waitData.Isset(key) {
		waitData.DecExpired(key, expired)
		//存储数据
		Store.SetData(readyData, waitData)
		return nil
	}
	return Utils.LogErr(1001)
}

/**
合并数据成[]byte
**/
func ComposeData(command int, stats bool) []byte {
	return nil
}

/**
发送数据
**/
func Send(sendData []byte, conn net.Conn) {

}
