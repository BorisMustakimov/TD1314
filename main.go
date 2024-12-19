package main

import (
	"fmt"
	"os"

	"github.com/BorisMustakimov/TD1314/server"
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
