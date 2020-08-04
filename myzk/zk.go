package myzk

import (
	"log"
	"time"

	"github.com/go-zookeeper/zk"
)

type zkconn = zk.Conn

type conn struct {
	*zkconn
}

//ConnPtr 连接指针
type ConnPtr = *conn

//NewZkConn 创建myzk连接
func NewZkConn(servers []string, sessionTimeout time.Duration) (ConnPtr, error) {
	c, _, err := zk.Connect(servers, sessionTimeout)
	if err != nil {
		return nil, err
	}
	pconn := new(conn)
	pconn.zkconn = c
	return pconn, nil
}

//NodeNotify 节点变更通知
type NodeNotify interface {
	OnCreated(string)
	OnDeleted(string)
	OnData(string, bool, []byte)
}

//ChildrenNotify 孩子变更通知
type ChildrenNotify interface {
	OnChildren(string, bool, []string)
}

//LoopWatchChildren 循环观察孩子节点的添加与删除
func (c *conn) LoopWatchChildren(parent string, cb ChildrenNotify) {
	for {
		children, _, ch, err := c.ChildrenW(parent)
		if err != nil && err != zk.ErrNoNode {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		if err == zk.ErrNoNode {
			cb.OnChildren(parent, false, []string{})
			exist, _, ch, err := c.ExistsW(parent)
			if err != nil {
				log.Println(err)
				continue
			}
			if exist {
				continue
			}
			e, ok := <-ch
			log.Println(e, ok)
			continue
		}
		cb.OnChildren(parent, true, children)
		e, ok := <-ch
		log.Println(e, ok)
	}
}

//LoopWatchNode 循环观察节点
func (c *conn) LoopWatchNode(path string, cb NodeNotify) {
	for {
		data, _, ch, err := c.GetW(path)
		if err != nil && err != zk.ErrNoNode {
			log.Println(err)
			time.Sleep(time.Second)
			continue
		}
		if err == zk.ErrNoNode {
			cb.OnData(path, false, []byte{})
			exist, _, ch, err := c.ExistsW(path)
			if err != nil {
				log.Println(err)
				continue
			}
			if exist {
				continue
			}
			e, ok := <-ch
			log.Println(e, ok)
			if ok && e.Err == nil {
				switch e.Type {
				case zk.EventNodeCreated:
					cb.OnCreated(e.Path)
				case zk.EventNodeDeleted:
					cb.OnDeleted(e.Path)
				}
			}
			continue
		}
		cb.OnData(path, true, data)
		e, ok := <-ch
		log.Println(e, ok)
		if ok && e.Err == nil {
			switch e.Type {
			case zk.EventNodeDeleted:
				cb.OnDeleted(e.Path)
			}
		}
	}
}

func (c *conn) WathcNode(path string, nodeNotify NodeNotify, childrenNotify ChildrenNotify) {
	go c.LoopWatchNode(path, nodeNotify)
	go c.LoopWatchChildren(path, childrenNotify)
}
