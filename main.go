package main

import (
	"log"
	"time"
	"zookeeper/myzk"
)

type testNode struct {
}

func (n *testNode) OnCreated(path string) {
	log.Println(path, "create")
}

func (n *testNode) OnDeleted(path string) {
	log.Println(path, "delete")
}

func (n *testNode) OnData(path string, exist bool, data []byte) {
	log.Println(path, "data", exist, string(data))
}

func (n *testNode) OnChildren(parent string, exist bool, children []string) {
	log.Println(parent, "children", exist, children)
}

func main() {
	log.SetFlags(log.Llongfile | log.LstdFlags)
	conn, err := myzk.NewZkConn([]string{"192.168.31.3:2181"}, time.Second*10)
	if err != nil {
		panic(err)
	}
	var node testNode
	conn.WathcNode("/bbb", &node, &node)
	ch := make(chan int)
	<-ch
}
