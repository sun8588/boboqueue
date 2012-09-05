package Entity

/**
已经就绪的任务列表
从wait列表中cp过来
**/
import (
//"time"
)

type ReadyData map[string]rData
type rData struct {
	Value   interface{}
	Expired uint
}

/**
添加一个准备好的任务
**/
func (wd ReadyData) Add(key string, value interface{}, expired uint) error {
	wd[key] = rData{value, expired}
	return nil
}

/**
删除一个任务
**/
func (wd ReadyData) Del(key string) error {
	delete(wd, key)
	return nil
}
/**
检查值是否存在
**/
func (wd ReadyData)Isset(key string)bool{
	if _,exist:=wd[key];exist{
		return true;
	}
	return false
}


/**
更新任务值
**/
func (wd ReadyData) UpdateValue(key string, value interface{}) error {
	data := wd[key]
	data.Value = value
	return nil
}
/**
从就绪列表中，取出对应数量的任务
**/
func (wd ReadyData) Cut(num int)*ReadyData {
	data:=make(ReadyData)
	for key, value := range wd {
		if num > 0 {
			data[key] = value
			delete(wd, key)
			num--
		} else {
			break
		}
	}
	return &data
}
