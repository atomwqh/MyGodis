package set

import "github.com/atomwqh/MyGodis/datastruct/dict"

// Set is a set of elements base on hash table
type Set struct {
	dict dict.Dict
}

// Make creates a new set
func Make(members ...string) *Set {
	set := &Set{
		dict: dict.MakeSimple(),
	}
	for _, member := range members {
		set.Add(member)
	}
	return set
}

func (s *Set) Add(val string) int {
	return s.dict.Put(val, nil)
}

func (s *Set) Remove(val string) int {
	_, res := s.dict.Remove(val)
	return res
}

func (s *Set) Has(val string) bool {
	if s == nil || s.dict == nil {
		return false
	}
	_, exist := s.dict.Get(val)
	return exist
}

func (s *Set) Len() int {
	if s == nil || s.dict == nil {
		return 0
	}
	return s.dict.Len()
}

// ToSlice convert set to []string
func (s *Set) ToSlice() []string {
	slice := make([]string, s.Len())
	i := 0
	s.dict.ForEach(func(key string, val interface{}) bool {
		if i < len(slice) {
			slice[i] = key
		} else {
			// set 在转换成slice的时候有可能变大，因为这里实现使用的是simpleDict
			slice = append(slice, key)
		}
		i++
		return true
	})
	return slice
}

func (s *Set) ForEach(consumer func(member string) bool) {
	if s == nil || s.dict == nil {
		return
	}
	s.dict.ForEach(func(key string, val interface{}) bool {
		return consumer(key)
	})
}

// 浅拷贝到另一个set，只拷贝引用位置
func (s *Set) ShallowCopy() *Set {
	res := Make()
	s.ForEach(func(member string) bool {
		res.Add(member)
		return true
	})
	return res
}

func Intersect(s ...*Set) *Set {
	result := Make()
	if len(s) == 0 {
		return result
	}

	countMap := make(map[string]int)
	for _, set := range s {
		set.ForEach(func(member string) bool {
			countMap[member]++
			return true
		})
	}
	for k, v := range countMap {
		if v == len(s) {
			result.Add(k)
		}
	}
	return result
}

// Union adds two sets
func Union(sets ...*Set) *Set {
	result := Make()
	for _, set := range sets {
		set.ForEach(func(member string) bool {
			result.Add(member)
			return true
		})
	}
	return result
}

// 差集
func Diff(sets ...*Set) *Set {
	if len(sets) == 0 {
		return Make()
	}
	// 选一个集合然后把遍历另一个集合来删除=>差集
	result := sets[0].ShallowCopy()
	for i := 1; i < len(sets); i++ {
		sets[i].ForEach(func(member string) bool {
			result.Remove(member)
			return true
		})
		if result.Len() == 0 {
			break
		}
	}
	return result
}

// RandomMembers randomly returns keys of the given number, may contain duplicated key
func (set *Set) RandomMembers(limit int) []string {
	if set == nil || set.dict == nil {
		return nil
	}
	return set.dict.RandomKeys(limit)
}

// RandomDistinctMembers randomly returns keys of the given number, won't contain duplicated key
func (set *Set) RandomDistinctMembers(limit int) []string {
	return set.dict.RandomDistinctKeys(limit)
}
