package main

import (
	"fmt"
	"github.com/ikigai-gh/git/lib"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Expected one argument!")
		os.Exit(1)
	}

	repo := lib.Repository{Path: ".git"}

	switch os.Args[1] {
	case "log":
		{
			repo.Log()
		}
	default:
		{
			os.Exit(0)
		}
	}

}
