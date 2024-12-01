package list

// 判断给定 a 是否和期望值相等
type Expected func(a any) bool

// 遍历 List
type Consumer func(i int, v any) bool

type List interface {
	Add(val any)
	Get(index int) (val any)
	Set(index int) (val any)
	Insert(index int) (val any)
	Remove(index int) (val any)
	RemoveLast() (val any)
	RemoveAllByVal(expected Expected) int
	RemoveByVal(expected Expected, count int) int
	// ReverseRemoveByVal()
	Len() int
	ForEach(consumer Consumer)
	Contains(expected Expected) bool
	Range(start int, stop int) []any
}
