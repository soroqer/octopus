package core

import (
	"errors"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/sirupsen/logrus"
)

var ZeroAddr = common.Address{}

type Galaxy struct {
	Path1 int // 一级衍生，为了方便扩展，目前 galaxy 服务里固定为 1001
	Path2 int // 二级衍生，galaxy 当前的二级衍生最大路径
	Nodes []Node // 衍生节点
}


func (galaxy *Galaxy)Init() error {
	galaxy.Path1 = 1001
	galaxy.Path2 = 0
	galaxy.Nodes = make([]Node,0,0)
	return nil
}


func (galaxy *Galaxy)Load(maxPath2 int, master *Master, chNode chan interface{}) error {
	galaxy.Path2 = maxPath2
	galaxy.Nodes = make([]Node,maxPath2+1,maxPath2+1)
	for i := range galaxy.Nodes {
		galaxy.Nodes[i].Path = NewPath(galaxy.Path1,i)
		galaxy.Nodes[i].Steps = make([]*Step,0,10)
	}
	err := galaxy.DeliverKeys(master)
	if err != nil {
		return err
	}
	for i := range chNode {
		switch v := i.(type) {
		case error:
			return v
		case *Node:
			paths,err := v.ParsePath()
			if err != nil {
				return err
			}
			if paths[0] != galaxy.Path1 {
				continue
			}
			if paths[1] > (len(galaxy.Nodes) -1) {
				logrus.Warnf("LvDb存储中路径 %v 的节点在超出配置中节点的数量，忽略并继续加载下一个节点",paths[0])
				continue
			}
			node := &galaxy.Nodes[paths[1]]
			if node.Path != v.Path || node.Addr != v.Addr {
				return errors.New("存储有错误，请检查代码。")
			}
			node.StepType = v.StepType
			if v.Steps != nil {
				node.Steps = v.Steps
			}
			// 在重新加载的时候，如果初始化过了，就都重置为待命状态
			node.ReSetToStandBy()
		}
	}

	return nil
}

func (galaxy *Galaxy)DeliverKeys(master *Master) (err error){
	galaxy.Range(func(node *Node) (next bool) {
		node.auth,node.key,err = master.DeriveAuth(node.Path)
		if err != nil {
			return false
		}
		node.Addr = node.auth.From
		node.sunAddr = master.auth.From
		return true
	})
	return
}

func (galaxy *Galaxy) AddNode(master *Master) (node *Node,err error) {
	galaxy.Nodes = append(galaxy.Nodes,Node{})
	galaxy.Path2 = len(galaxy.Nodes) - 1
	node = &galaxy.Nodes[galaxy.Path2]
	node.Path = NewPath(galaxy.Path1,galaxy.Path2)
	node.Steps = make([]*Step,0,10)

	node.auth,node.key,err = master.DeriveAuth(node.Path)
	if err != nil {
		return nil,err
	}
	node.Addr = node.auth.From
	node.sunAddr = master.auth.From
	return
}


func (galaxy *Galaxy) Range(fn func(node *Node)(next bool)) {
	for i := range galaxy.Nodes {
		if !fn(&galaxy.Nodes[i]) {
			return
		}
	}
}

func (galaxy *Galaxy) GetNode(path Path) (*Node,error) {
	paths,err := path.ParsePath()
	if err != nil {
		return nil, err
	}
	if paths[1] > (len(galaxy.Nodes) -1) {
		return nil,errors.New(fmt.Sprintf("child path %v not found",paths[0]))
	}
	return  &galaxy.Nodes[paths[1]],nil
}

func (galaxy *Galaxy) GetNodeByPath2(path2 int) (*Node,error) {
	if path2 > (len(galaxy.Nodes) -1) {
		return nil,errors.New(fmt.Sprintf("child path %v not found",path2))
	}
	return  &galaxy.Nodes[path2],nil
}

func (galaxy *Galaxy) GetNodeByAddr(addr common.Address) (node *Node) {
	galaxy.Range(func(v *Node) (next bool) {
		if v.Addr == addr {
			node = v
			return false
		}
		return true
	})
	return
}
