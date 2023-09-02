package main

import (
	"os"
)

func noErrInNotMainFunc() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	os.Exit(1)
}

func main() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	os.Exit(1) // want "os.Exit returns"
}
