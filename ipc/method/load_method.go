package method

import (
	"github.com/HaiChecker/go-fb/ipc"
	"github.com/HaiChecker/go-fb/ipc/pack"
	"log"
)

func init() {
	ipc.RegisterMethod(ipc.Load, NewGetImageMethod)
}

func NewGetImageMethod(name byte) (ipc.HandlerMethod, error) {
	log.Printf("注册函数:%v - Load FrameBuffer", name)
	return &GetImageMethod{name: name}, nil
}

type GetImageMethod struct {
	name byte
}

func (method *GetImageMethod) Name() string {
	return "load"
}

func (method *GetImageMethod) Invoke(uuid []byte, args ...*pack.FBString) (result pack.Response, err error) {
	err = ipc.Framebuffer.Load(nil)
	if err != nil {
		return pack.CreateResponse(uuid, pack.CreateBool(false)), err
	}
	return pack.CreateResponse(uuid, pack.CreateBool(true)), err
}
