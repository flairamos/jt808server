package main

import (
	"bufio"
	"fmt"
	"jt808server/jt808"
	"net"
)

// 处理函数
func process(conn net.Conn) {
    defer conn.Close() // 关闭连接
    for {
        reader := bufio.NewReader(conn)
        var buf [128]byte
        n, err := reader.Read(buf[:]) // 读取数据
        if err != nil {
            fmt.Println("read from client failed, err:", err)
            break
        }
        recvStr := string(buf[:n])
        fmt.Println("收到client端发来的数据：", recvStr)
        fmt.Println("解析后的数据为",jt808.ParseJT808Message(recvStr))
        send := jt808.ServerCommonReplyMessage(recvStr)
        conn.Write([]byte(send)) // 发送数据
    }
}

func main() {
    listen, err := net.Listen("tcp", "127.0.0.1:9999")
    if err != nil {
        fmt.Println("listen failed, err:", err)
        return
    }
	fmt.Println("server is runing on 9999 port...")
    for {
        conn, err := listen.Accept() // 建立连接
        if err != nil {
            fmt.Println("accept failed, err:", err)
            continue
        }
        go process(conn) // 启动一个goroutine处理连接
    }
}