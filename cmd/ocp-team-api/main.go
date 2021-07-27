package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("The ocp-team-api project.")

	for i := 0; i < 5; i++ {
		if err := ReadFile("go.mod"); err != nil {
			fmt.Println(err)
			break
		}
	}
}

func ReadFile(path string) (err error) {
	file, err := os.Open(path)

	if err != nil {
		return
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		} else {
			fmt.Printf("File %s successfully closed\n", file.Name())
		}
	}(file)

	return nil
}