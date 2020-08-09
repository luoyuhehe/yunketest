package utils

import (
	"math/rand"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

var snowflakeNode *snowflake.Node
var newNodeOnce sync.Once

// UniqueID 唯一标识
type UniqueID snowflake.ID

// newNodeID 生成随机 node 码，node 影响的是 E10-bit。
func newNodeID() int64 {
	s := rand.NewSource(time.Now().UnixNano())
	return s.Int63() % 1024
}

// GenerateID 生成 ID
// MySQL 中应该使用 BIGINT 字段
func GenerateID() UniqueID {
	newNodeOnce.Do(func() {
		var err error
		snowflakeNode, err = snowflake.NewNode(newNodeID())
		if err != nil {
			panic(err)
		}
	})
	id := snowflakeNode.Generate()
	return UniqueID(id)
}

// Int64 转成int64
func (v UniqueID) Int64() int64 {
	return int64(v)
}

func (v UniqueID) UInt64() uint64 {
	return uint64(v)
}

func (v UniqueID) Equal(id UniqueID) bool {
	return v == id
}

func (v UniqueID) IsEmpty() bool {
	if v == 0 {
		return true
	}
	return false
}

func ParseUniqueID(id uint64) UniqueID {
	return UniqueID(id)
}
