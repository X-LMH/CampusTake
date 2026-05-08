package snowflake

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
	once sync.Once
)

// InitSnowFlake 初始化雪花节点
func InitSnowFlake(nodeID int64) error {
	var err error

	once.Do(func() {
		node, err = snowflake.NewNode(nodeID)
	})

	return err
}

// GenerateID 生成ID
func GenerateID() string {
	if node == nil {
		panic("snowflake node not initialized")
	}

	return node.Generate().String()
}
