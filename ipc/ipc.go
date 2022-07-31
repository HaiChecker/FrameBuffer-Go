package ipc

import (
	"github.com/HaiChecker/go-fb/ipc/pack"
	"net"
)

type BasicMethod struct {
	pack.Request
	net.Conn
}

func CreateNewMethod(conn net.Conn) BasicMethod {
	return BasicMethod{Conn: conn}
}

// Request
// +-------------+---------------+--------------------+---------+------------+--------------+
// | Request ID  | FunctionName  | FunctionPairCount  | Pair    | Pair Data  | Return       |
// +=============+===============+====================+=========+============+==============+
// | 16          | 1             | 2                  | 2       | N          | 1            |
// +-------------+---------------+--------------------+---------+------------+--------------+
// | UUID        | 函数名称       | Count              | Length  | Data（x）   | Return Type  |
// +-------------+---------------+--------------------+---------+------------+--------------+

func (b *BasicMethod) HandlerInfo() (h HandlerMethod, err error) {
	b.Request, err = pack.ParseRequest(b)
	if err != nil {
		return nil, err
	}
	return GetMethodFromName(b.Request.FunctionName)
}

func (b *BasicMethod) Return(response pack.Response) error {
	_, err := b.Write(response.ToByte())
	if err != nil {
		return err
	}
	return nil
}
