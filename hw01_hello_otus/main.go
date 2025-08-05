package main

// Импортируем необходимые библиотеки.
import (
	"fmt"

	"golang.org/x/example/hello/reverse"
)

// Основная функция программы.
func main() {
	// Печать в стандартный вывод перевернутой строки.
	fmt.Println(reverse.String("Hello, OTUS!"))
}
