package set

// Set 使用map实现一个set集合(string)
type Set struct {
	items map[string]struct{}
}

func NewStringSet() *Set {
	return &Set{
		items: make(map[string]struct{}),
	}
}

// Add 添加元素
func (s *Set) Add(item string) {
	s.items[item] = struct{}{}
}

// Remove 删除元素
func (s *Set) Remove(item string) {
	delete(s.items, item)
}

// Contains 判断元素是否存在
func (s *Set) Contains(item string) bool {
	_, exists := s.items[item]
	return exists
}

// List 获取所有元素
func (s *Set) List() []string {
	items := make([]string, 0, len(s.items))
	for item := range s.items {
		items = append(items, item)
	}
	return items
}
