package Model

import (
"os"
//"log"
"bufio"
)
type FileInfo struct{
	Fd *os.File
//	IModel
}
func New(fileName string)(*FileInfo,error){
		files, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0775)
		if err != nil {
			return nil,err
		}
	return &FileInfo{Fd:files},nil
}
func (f *FileInfo)Set(data []byte)(int,error){
	return f.Fd.Write(data)
}
func (f *FileInfo)Get() (*bufio.Reader){
	return bufio.NewReader(f.Fd)
	
}
func (f *FileInfo)GetFd()interface{}{
	return f.Fd
}
