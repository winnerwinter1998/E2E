package server

import (
	"fmt"
	"net"
	"strings"
)

// User can create a user
type User struct {
	Username      string
	OtherUsername string
	Msg           string
	ServerMsg     string
}

var (
	userMap = make(map[string]net.Conn)
	user    = new(User)
)

// Run runs the tcp server
func Run() {
	addr, _ := net.ResolveTCPAddr("tcp4", "localhost:8899")
	list, _ := net.ListenTCP("tcp4", addr)
	for {
		conn, _ := list.Accept()
		go func() {
			for {
				b := make([]byte, 1024)
				count, _ := conn.Read(b)
				array := strings.Split(string(b[:count]), "-")
				user.Username = array[0]
				user.OtherUsername = array[1]
				user.Msg = array[2]
				user.ServerMsg = array[3]
				userMap[user.Username] = conn

				if v, ok := userMap[user.OtherUsername]; ok && v != nil {
					n, err := v.Write([]byte(fmt.Sprintf("%s-%s-%s-%s", user.Username, user.OtherUsername, user.Msg, user.ServerMsg)))
					if n <= 0 || err != nil {
						delete(userMap, user.OtherUsername)
						conn.Close()
						v.Close()
						break
					}
				} else {
					user.ServerMsg = "对方不在线"
					conn.Write([]byte(fmt.Sprintf("%s-%s-%s-%s", user.Username, user.OtherUsername, user.Msg, user.ServerMsg)))
				}
			}
		}()
	}

}
