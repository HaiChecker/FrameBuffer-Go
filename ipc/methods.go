package ipc

import (
	"errors"
	"fmt"
	"github.com/HaiChecker/go-fb/ipc/pack"
	"github.com/HaiChecker/go-fb/lib"
)

type HandlerMethod interface {
	// Invoke 执行函数
	Invoke(uuid []byte, args ...*pack.FBString) (result pack.Response, err error)
	// Name 函数名称
	Name() string
}

var methods = make(map[byte]MethodCreator)

var Framebuffer lib.AndroidFramebuffer

type MethodCreator func(name byte) (HandlerMethod, error)

func RegisterMethod(name MethodID, c MethodCreator) {
	methods[byte(name)] = c
}

type (
	MethodID byte
)

const (
	Init      = MethodID(0x01)
	Load      = MethodID(0x02)
	SaveImage = MethodID(0x03)
	Ping      = MethodID(0x04)
	GetImage  = MethodID(0x05)
)

func GetMethodFromName(name byte) (HandlerMethod, error) {
	result, ok := methods[name]
	if ok {
		return result(name)
	}
	return nil, errors.New(fmt.Sprintf("%v 不存在", name))
}
