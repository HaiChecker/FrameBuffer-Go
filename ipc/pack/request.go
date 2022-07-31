package pack

import (
	"encoding/binary"
	"io"
	"log"
)

type Request struct {
	Uuid              []byte
	FunctionName      byte
	FunctionPairCount uint64
	Pairs             []*FBString
}

func ParseRequest(reader io.Reader) (request Request, err error) {
	defer func() {
		i := recover()
		if i != nil {
			log.Printf("解析函数调用失败")
		}
	}()
	request = Request{}

	// UUID
	header := make([]byte, 25)
	_, err = reader.Read(header)
	if err != nil {
		return request, err
	}
	request.Uuid = header[:16]
	log.Printf("UUID:%v", len(request.Uuid))
	request.FunctionName = header[16]
	request.FunctionPairCount = binary.BigEndian.Uint64(header[17:])

	request.Pairs = make([]*FBString, request.FunctionPairCount)
	// Function
	for i := 0; i < int(request.FunctionPairCount); i++ {
		readString, errs := StringConvert(reader)
		if errs != nil {
			return request, errs
		}
		request.Pairs[i] = readString
	}
	return request, nil
}
