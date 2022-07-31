package image

import (
	"bytes"
	"github.com/HaiChecker/go-fb/ipc"
	"github.com/HaiChecker/go-fb/ipc/pack"
	"image/png"
	"log"
)

func init() {
	ipc.RegisterMethod(ipc.GetImage, NewGetImageMethod)
}

type GetImageMethod struct {
	name byte
}

func NewGetImageMethod(name byte) (ipc.HandlerMethod, error) {
	log.Printf("注册函数:%v - 保存当前图片到系统指定路径", name)
	return &GetImageMethod{name: name}, nil
}

func (method *GetImageMethod) Name() string {
	return "getImage"
}

func (method *GetImageMethod) Invoke(uuid []byte, args ...*pack.FBString) (result pack.Response, err error) {
	buf := new(bytes.Buffer)
	err = png.Encode(buf, ipc.Framebuffer.Image)
	return pack.CreateResponse(uuid, pack.CreateByte(buf.Bytes())), err
}
