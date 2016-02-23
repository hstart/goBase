package socket

/*
基础组件包:
socket相关的通讯处理相关函数
*/

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"time"
)

func init() {
	fmt.Println("Simple Main Test @ goBase/socket")
}

//启动监听处理传入的参数是端口好，和连接处理函数
func StartListen(ListenPort string, FuncHandleConn interface{}) error {

	fmt.Println("服务启动监听:", ListenPort)

	//定义一个函数指针
	pfuncHandle := func(conn net.Conn) {}

	//将接口转换为处理函数
	switch pfuncHandleTmp := FuncHandleConn.(type) {
	case func(conn net.Conn):
		pfuncHandle = pfuncHandleTmp
		//fmt.Println(pfuncHandle)
	default:
		fmt.Println(" 启动监听错误,参数错误 ")
		return errors.New("参数错误@第二个参数非处理函数")
	}

	ln, err := net.Listen("tcp", ":"+ListenPort)
	if err != nil {
		fmt.Println("start listen ERR:", err)
		panic("start listen")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("有新连接接入:", conn)
		go pfuncHandle(conn) //启动一个线程处理连接
	}
	return nil
}

//短连接测试
func ShortHandleConnTest(conn net.Conn) {
	defer conn.Close()
	fmt.Println("ShortHandleConnTest @ 接收数据: ", time.Now())
	buffer := make([]byte, 4*1024) //31个字节的长度头
	reader := bufio.NewReader(conn)
	bpi := []byte{0x45, 0x2E}
	n, err := reader.Read(buffer)
	if err != nil {
		_, err = conn.Write(bpi)
		if err != nil {
			fmt.Println("ErrWrite=", err)
		}
		return
	}
	fmt.Println("接收数据@ ", n, string(buffer))
	conn.Write(bpi)
	fmt.Println("短连接处理完成")
	return
}

func SocketTest() {

	//定义一个处理函数
	f := func(conn net.Conn) {
		ShortHandleConnTest(conn)
	}

	StartListen("8510", f)

}
