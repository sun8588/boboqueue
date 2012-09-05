package Entity

import (
"container/list"
//"Store"
)

type TaskIndex map[string]TaskAddr
type TaskAddr struct{
	Key string 
	EleAddr *list.Element
	FileAddr int64
}
type TaskList *list.List
type TaskData struct{
	Key string
	Value interface{}
	Expired int
	FileAddr int64
}
/**
初始化task列表
**/
//func GTaskList()(TaskIndex,*list.List,error){
//	taskAddrList,err:=Store.GetAllTaskIndex()
//	if err!=nil{
//		return nil,nil,err
//	}
//	taskDataList,err:=Store.GetAllTaskData()
//	if err!=nil{
//		return nil,nil,err
//	}
//	taskIndex:=make(TaskIndex)
//	
//	for _,taskAddrInfo:=range taskAddrList.([]TaskAddr){
////		taskAddrInfo:=taskAddr.(TaskAddr)
//		taskIndex[taskAddrInfo.Key]=taskAddrInfo
//	}	
//	taskList :=list.New() 
//	for _,taskData:=range taskDataList.([]TaskData){
//		p:=taskList.PushBack(taskData)
//	//重新设置list.ele的内存地址
//		m:=taskIndex[taskData.Key]
//		m.EleAddr=p
//	}
//	return taskIndex,taskList,nil
//}
func (ti TaskIndex)Add(key string,eleAddr *list.Element ,fileAddr int64)TaskIndex{
	ti[key]=TaskAddr{key,eleAddr,fileAddr}
	return ti
}
/**
处理添加任务入口
**/
func AddTask(taskIndex TaskIndex,taskList *list.List,key string,value interface{},expired int){
//	taskData:=TaskData{key,value,expired,0}
//	taskIndex.Add(key,taskList.PushBack(taskData),0)
//	Store.AddIndex(taskIndex[key])
//	Store.AddData(taskData)
}

