/*
   服务器端程序
   接收客户端连接
   将客户端发送的数据写入记录文件中
   每个新连接都会创建新记录文件
*/

package main

import (
	"io"
	"log"
	"net"
	"os"
)

const (
	NETWORK string = "tcp"   //socket网络协议
	LADDR   string = ":8080" //绑定的本机地址和端口

)

var (
	logger *log.Logger = log.New(os.Stdout, "", log.Ldate|log.Ltime) //创建日志对象
)

func main() {
	listener, err := net.Listen(NETWORK, LADDR) //创建socket,绑定绑定端口,实现监听
	if err != nil {
		logger.Printf("监听端口失败!error:%s", err) //将错误写入日志文件中
		os.Exit(1)
	}
	logger.Print("服务启动...")
	defer listener.Close() //延迟关闭socket
	for {
		conn, err := listener.Accept() //创建连接
		logger.Print("收到客户端连接:", conn.RemoteAddr().String())

		if err != nil {
			logger.Printf("创建连接失败!error:err%s", err)
			os.Exit(1)
		}
		go connHandle(conn) //创建协程处理连接
	}

}

func connHandle(conn net.Conn) {
	defer conn.Close() //延迟关闭连接

	Raddr := conn.RemoteAddr().String() //客户端的IP和端口号

	prefix := Raddr + ":"
	logger.Println("客户端:" + Raddr + "已连接")
	var buf []byte = make([]byte, 4096)
	for {
		n, err := conn.Read(buf) //将客户端发送的数据写入buf中
		if err != nil {
			if err == io.EOF {
				logger.Print(prefix, "socket连接已关闭!")
				break
			} else {
				logger.Print(prefix, "写入数据失败!error:"+err.Error())
				break
			}
		}
		logger.Print(prefix, string(buf[:n]))
		if string(buf[:n]) == "ping" {
			conn.Write([]byte("pong"))
		}
	}
}
