package dict

// 用来遍历dict键值对，return false means the traversal will be break
type Consumer func(key string, val interface{}) bool

// Dict 键值对的接口
type Dict interface {
	Get(key string) (val any, exixts bool)
	Len() int
	Put(key string, val any) (result int)
	PutIfAbsent(key string, val any) (result int)
	Remove(key string) (val any, result int)
	ForEach(consuemr Consumer)
	Keys() []string
	RandomKeys(limit int) []string
	RandomDistinctKeys(limit int) []string
	Clear()
	DictScan(cursor int, count int, pattern string) ([][]byte, int)
}
