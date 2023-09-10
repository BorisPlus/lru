package lru

import "fmt"

// ListItem - элемент двусвязного списка.
type ListItem struct {
	Data interface{}
	Prev *ListItem
	Next *ListItem
}

// String - наглядное представление значения элемента двусвязного списка.
//
// Например,
//
//	-------------------             -------------------
//	Item: 0xc00002e400              Item: 0xc00002e400
//	-------------------             -------------------
//	Data: 30                или     Data: 30
//	Prev: 0xc00002e3c0              Prev: 0x0
//	Next: 0xc00002e440              Next: 0x0
//	-------------------             -------------------
func (listItem *ListItem) String() string {
	template := `
-------------------
Item: %p
-------------------
Data: %v
Prev: %p
Next: %p
-------------------`
	return fmt.Sprintf(template, listItem, listItem.Data, listItem.Prev, listItem.Next)
}
