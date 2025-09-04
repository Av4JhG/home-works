package hw04lrucache

type Key string

// Интерфейс для работы с кэшем.
type Cache interface {
	Set(key Key, value interface{}) bool // добавить значение в кэш по ключу.
	Get(key Key) (interface{}, bool)     // получить значение из кэша по ключу.
	Clear()                              // очистить кэш.
}

// Структура элемента кэша.
type cacheItem struct {
	key   Key         // ключ элемента кэша.
	value interface{} // значение элемента кэша.
}

// Структура LRU кэша.
type lruCache struct {
	capacity int               // емкость кэша.
	queue    List              // очередь.
	items    map[Key]*ListItem // мап элементов кэша.
}

// Передача значения элемента списка в элемент списка.
func listItemToCacheItem(item *ListItem) *cacheItem {
	return item.Value.(*cacheItem)
}

// Создание элемента кэша.
func newCacheItem(key Key, value interface{}) *cacheItem {
	return &cacheItem{key, value}
}

// Реализация добавления элемента в кэш.
func (l *lruCache) Set(key Key, value interface{}) bool {
	item, ok := l.items[key]

	if ok {
		l.queue.MoveToFront(item)
		listItemToCacheItem(item).value = value
		return true
	}

	if l.queue.Len() == l.capacity {
		removable := l.queue.Back()
		l.queue.Remove(removable)
		delete(l.items, listItemToCacheItem(removable).key)
	}

	item = l.queue.PushFront(newCacheItem(key, value))
	l.items[key] = item

	return false
}

// Реализация получения элемента из кэша.
func (l *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := l.items[key]

	if ok {
		l.queue.MoveToFront(item)
		return listItemToCacheItem(item).value, true
	}

	return nil, false
}

// Реализация очистки кэша.
func (l *lruCache) Clear() {
	for l.queue.Len() > 0 {
		l.queue.Remove(l.queue.Back())
	}
	l.items = make(map[Key]*ListItem, l.capacity)
}

// Функция инициализации кэша (конструктор).
func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
