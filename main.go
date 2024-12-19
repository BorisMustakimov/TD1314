package main

import (
	"fmt"
	"os"

)

func main() {
	a, err := server.New()
	if err != nil {
		fmt.Println("ошибка запуска приложения", err)
		os.Exit(1)
	}

	if err = a.Run(); err != nil {
		fmt.Println("ошибка запуска приложения", err)
	}
}
