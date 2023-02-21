package main

import (
	"fmt"
	"github.com/ikigai-gh/git/lib"
)

func main() {
	repo := lib.Repository{Path: ".git"}
	fmt.Println(repo.ListObjects())
}
