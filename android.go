package fb

import (
	"github.com/HaiChecker/go-fb/ipc"
	_ "github.com/HaiChecker/go-fb/ipc/method"
	_ "github.com/HaiChecker/go-fb/ipc/method/image"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Start		Java app_process
func Start(host string, devFile string) {
	_, _ = exec.Command("/system/bin/su", "-c", "/system/bin/sh", "chmod 777 "+devFile).CombinedOutput()
	//syscall.Umask(0777)
	os.Remove(host)
	addr, err := net.ResolveUnixAddr("unix", host)
	if err != nil {
		panic("Cannot resolve unix addr: " + err.Error())
	}
	listen, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Printf("跨进程通信启动失败:%v", err)
		return
	}
	defer listen.Close()
	log.Printf("跨进程通信启动成功 -> %v", host)
	for listen != nil {
		accept, err := listen.Accept()
		if err != nil {
			if strings.Contains(err.Error(), "closed") {
				log.Printf("关闭")
				break
			}
			log.Printf("通信失败")
			continue
		}
		go func(conn net.Conn) {
			defer func(conn net.Conn) {
				err := conn.Close()
				if err != nil {
					log.Printf("关闭当前通信失败")
				}

			}(conn)
			conn.SetDeadline(time.Now().Add(time.Second * 5))
			method := ipc.CreateNewMethod(conn)
			info, err := method.HandlerInfo()
			if err != nil {
				log.Printf("解析函数信息失败：%v", err)
				return
			}
			log.Printf("函数:%v", info.Name())

			res, err := info.Invoke(method.Uuid, method.Request.Pairs...)
			if err != nil {
				log.Printf("函数:%v 执行时发生错误:%v", info.Name(), err)
				return
			}
			err = method.Return(res)
			if err != nil {
				log.Printf("函数:%v 处理返回结果错误:%v", info.Name(), err)
				return
			}
			log.Printf("函数:%v 返回值:%v", info.Name(), res.Data.String())
		}(accept)
	}
	log.Printf("结束")

}
