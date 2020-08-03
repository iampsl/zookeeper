package main

import (
	"log"
	"time"
	"zookeeper/myzk"
)

type testNode struct {
}

func (t testNode) OnCreated(path string) {
	log.Printf("%s OnCreated\n", path)
}

func (t testNode) OnDeleted(path string) {
	log.Printf("%s OnDeleted\n", path)
}

func (t testNode) OnData(path string, exist bool, data []byte) {
	log.Printf("%s exist=%v %s\n", path, exist, string(data))
}

func (t testNode) OnChildren(parent string, exist bool, children []string) {
	log.Println(parent, exist, children)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Llongfile)
	c, err := myzk.NewZkConn([]string{"192.168.31.3:2181"}, time.Second*10)
	if err != nil {
		panic(err)
	}
	var t testNode
	c.LoopWatchChildren("/bbb", &t)
}
