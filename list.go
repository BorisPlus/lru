package lru

// List - структура двусвязного списка.
type List struct {
	len   int
	front *ListItem
	back  *ListItem
}

// Len() - получить длину двусвязного списка.
func (list *List) Len() int {
	return list.len
}

// Front() - получить первый элемент двусвязного списка.
func (list *List) Front() *ListItem {
	return list.front
}

// Back() - получить последний элемент двусвязного списка.
func (list *List) Back() *ListItem {
	return list.back
}

// PushFront() - добавить значение в начало двусвязного списка.
func (list *List) PushFront(data interface{}) *ListItem {
	item := &ListItem{
		Data: data,
	}
	if list.len == 0 {
		list.front = item
		list.back = item
	} else {
		item.Next = list.front
		list.front.Prev = item
		list.front = item
	}
	list.len++
	return item
}

// PushBack() - добавить значение в конец двусвязного списка.
func (list *List) PushBack(data interface{}) *ListItem {
	item := &ListItem{
		Data: data,
	}
	if list.len == 0 {
		list.front = item
		list.back = item
	} else {
		item.Prev = list.back
		list.back.Next = item
		list.back = item
	}
	list.len++
	return item
}

// Remove() - удалить элемент из двусвязного списка.
func (list *List) Remove(i *ListItem) {
	// TODO:
	// не наглядно, но интересно
	// Ai, Aj = Aj, Ai
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		list.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		list.back = i.Prev
	}
	list.len--
}

// MoveToFront() - переместить элемент в начало двусвязного списка.
func (list *List) MoveToFront(i *ListItem) {
	if i.Prev != nil {
		i.Prev.Next = i.Next
	} else {
		list.front = i.Next
	}
	if i.Next != nil {
		i.Next.Prev = i.Prev
	} else {
		list.back = i.Prev
	}

	i.Prev = nil
	i.Next = list.front
	if list.front != nil {
		list.front.Prev = i
		list.front = i
	} else {
		list.front = i
		list.back = i
	}
}

// Lister - интерфейс двусвязного списка.
type Lister interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}) *ListItem
	PushBack(v interface{}) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

func NewList() Lister {
	return new(List)
}

// String - наглядное представление всего двусвязного списка.
//
// Например,
//
// - пустой список:
//
//	(nil:0x0)
//	    |
//	    V
//	(nil:0x0)
//
// - список из двух элементов:
//
//	    (nil:0x0)
//	        |
//	        V
//	-------------------
//	Item: 0xc00002e3a0 <--------┐
//	-------------------         |
//	Data: 2                     |
//	Prev: 0x0                   |
//	Next: 0xc00002e380  >>>-----|---┐ Next 0xc00002e380
//	-------------------         |   | ссылается на
//	        |                   |   | блок 0xc00002e380
//	        V                   |   |
//	-------------------         |   |
//	Item: 0xc00002e380  <-----------┘
//	-------------------         | Prev 0xc00002e3a0
//	Data: 1                     | ссылается на
//	Prev: 0xc00002e3a0  >>>-----┘ блок 0xc00002e3a0
//	Next: 0x0
//	-------------------
//	        |
//	        V
//	    (nil:0x0)
func (list *List) String() string {
	result := ""
	Nill := `
    (nil:0x0)`
	delimiter := `
	|
	V`
	result += Nill
	result += delimiter
	for i := list.Front(); i != nil; i = i.Next {
		result += i.String()
		result += delimiter
	}
	result += Nill
	return result
}
