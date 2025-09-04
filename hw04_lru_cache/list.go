package hw04lrucache

// Интерфейс работы со списком.
type List interface {
	Len() int                          // получение длины списка.
	Front() *ListItem                  // получение первого элемента списка.
	Back() *ListItem                   // получение последнего элемента списка.
	PushFront(v interface{}) *ListItem // добавить значение в начало.
	PushBack(v interface{}) *ListItem  // добавить значение в конец.
	Remove(i *ListItem)                // удалить элемент.
	MoveToFront(i *ListItem)           // переместить элемент в начало.
}

// Структура элемента списка
type ListItem struct {
	Value interface{} // значение элемента списка.
	Next  *ListItem   // предыдущий элемент списка.
	Prev  *ListItem   // следующий элемент списка.
}

// Структура списка
type list struct {
	len   int       // длина списка.
	front *ListItem // первый элемент списка.
	back  *ListItem // последний элемент списка.
}

// Реализация получения длины списка.
func (l *list) Len() int {
	return l.len
}

// Реализация получения первого элемента списка.
func (l *list) Front() *ListItem {
	return l.front
}

// Реализация получения последнего элемента списка.
func (l *list) Back() *ListItem {
	return l.back
}

// Реализация добавления элемента в начало списка.
func (l *list) PushFront(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.len == 0 {
		l.back = item
	} else {
		item.Next = l.front
		l.front.Prev = item
	}
	l.front = item

	l.len++
	return item
}

// Реализация добавления элемента в конец списка.
func (l *list) PushBack(v interface{}) *ListItem {
	item := &ListItem{Value: v}

	if l.len == 0 {
		l.front = item
	} else {
		item.Prev = l.back
		l.back.Next = item
	}
	l.back = item

	l.len++
	return item
}

// Реализация удаления элемента списка.
func (l *list) Remove(i *ListItem) {
	if i.Prev == nil {
		l.front = i.Next
	} else {
		i.Prev.Next = i.Next
	}

	if i.Next == nil {
		l.back = i.Prev
	} else {
		i.Next.Prev = i.Prev
	}

	l.len--
}

// Реализация перемещения элемента в начало списка.
func (l *list) MoveToFront(i *ListItem) {
	if i.Prev == nil {
		return
	}

	l.Remove(i)

	i.Prev = nil
	i.Next = l.front
	l.front.Prev = i
	l.front = i
	l.len++
}

// Функция инициализации списка (конструктор)
func NewList() List {
	return new(list)
}
