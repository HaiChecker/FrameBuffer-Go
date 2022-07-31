package image

import (
	"github.com/HaiChecker/go-fb/ipc"
	"github.com/HaiChecker/go-fb/ipc/pack"
	"image/png"
	"log"
	"os"
)

func init() {
	ipc.RegisterMethod(ipc.SaveImage, NewSaveImageMethod)
}

type SaveImageMethod struct {
	name byte
}

func NewSaveImageMethod(name byte) (ipc.HandlerMethod, error) {
	log.Printf("注册函数:%v - 保存当前图片到系统指定路径", name)
	return &SaveImageMethod{name: name}, nil
}

func (method *SaveImageMethod) Name() string {
	return "saveImage"
}

func (method *SaveImageMethod) Invoke(uuid []byte, args ...*pack.FBString) (result pack.Response, err error) {
	file, err := os.OpenFile(args[0].String(), os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return pack.CreateResponse(uuid, pack.CreateBool(false)), err
	}
	err = png.Encode(file, ipc.Framebuffer.Image)
	if err != nil {
		return pack.CreateResponse(uuid, pack.CreateBool(false)), err
	}
	return pack.CreateResponse(uuid, pack.CreateBool(true)), nil
}
