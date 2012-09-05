package Entity

import (
"time"
)


type wData struct {
	Value   interface{}
	Expired uint
}
type WaitData map[string]wData

/**
添加一个任务
**/
func (wd *WaitData)Add(key string,value interface{},expired uint)error{
	(*wd)[key]=wData{value,uint(time.Now().Unix())+expired}
	return nil
}
/**
删除一个任务
**/
func (wd WaitData)Del(key string)error{
	delete(wd,key)
	return nil
}
/**
检查值是否存在
**/
func (wd WaitData)Isset(key string)bool{
	if _,exist:=wd[key];exist{
		return true;
	}
	return false
}
/**
增加到期时间
**/
func (wd WaitData)AddExpired(key string,expired uint)error{
	data:=wd[key]
	data.Expired+=expired
	wd[key]=data
	return nil
}
/**
减少到期时间
**/
func (wd WaitData)DecExpired(key string,expired uint)error{
	data:=wd[key]
	//如果需要减少的到期时间是0,则意味着，设置到期时间为0
	//如果到期时间小于要减少的时间，则直接设置为0
	if expired ==0 || data.Expired<=expired {
		data.Expired=0
	}else{
		data.Expired-=expired
	} 
	wd[key]=data
	return nil
}
/**
更新任务值
**/
func (wd WaitData)UpdateValue(key string,value interface{})error{
	data:=wd[key]
	data.Value=value
	return nil
}
