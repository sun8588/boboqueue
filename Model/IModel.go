package Model

import (
"bufio"

)

type IModel interface{
	GetFd()interface{}
	Set([]byte)(int,error)
	Get()(*bufio.Reader)
}
