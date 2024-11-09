package pkg1

import (
	"fmt"
	"os"
)

func errCheckFunc() {
	fmt.Println("hello world")
	os.Exit(1) // want "using os.Exit is prohibbited"
}
