package database

import (
	"github.com/atomwqh/MyGodis/interface/redis"
	"time"
)

type CmdLine = [][]byte

type DB interface {
	Exec(client redis.Connection, cmdLine [][]byte) redis.Reply
	AfterClientClose(c redis.Connection)
	Close()
	//TODO: add RDB mode
	//LoadRDB(dec *core.Decoder) error
}

type KeyEventCallback func(dbIndex int, key string, entity *DataEntity)

// DBEngine is the embedding storage engine exposing more methods for complex application
// 扩展存储引擎功能，暴露更多复杂应用程序可能需要的方法
type DBEngine interface {
	DB
	ExecWithLock(conn redis.Connection, cmdLine [][]byte) redis.Reply
	ExecMulti(conn redis.Connection, watching map[string]uint32, cmdLines []CmdLine) redis.Reply
	GetUndoLogs(dbIndex int, cmdLine [][]byte) redis.Reply
	ForEach(dbIndex int, cb func(key string, data *DataEntity, expiration *time.Time) bool)
	RWLocks(dbIndex int, writeKeys []string, readKeys []string)
	RWUnlocks(dbIndex int, writeKeys []string, readKeys []string)
	GetDBSize(dbIndex int) (int, int)
	GetEntity(dbIndex int, key string) (*DataEntity, bool)
	GetExpiration(dbIndex int, key string) *time.Time
	SetKeyInsertedCallback(cb KeyEventCallback)
	SetKeyDeletedCallback(cb KeyEventCallback)
}

type DataEntity struct {
	Data interface{}
}
