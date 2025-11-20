package main

import (
	"fmt"
	"os"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	filePath, choosed := chooseFile(wd)
	if choosed {
		fmt.Println("no file choose, exiting")
		return
	}

	err = parseFile(filePath)
	if err != nil {
		panic(err)
	}
}
