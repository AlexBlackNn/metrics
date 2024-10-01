package main

import (
	"go/ast"
	"go/parser"
	"go/token"
)

func main() {
	// исходный код, который будем разбирать
	src := `package main
import "fmt"

func main() {
    fmt.Println("Hello, world!")
}`

	// дерево разбора AST ассоциируется с набором исходных файлов FileSet
	fset := token.NewFileSet()
	// парсер может работать с файлом
	// или исходным кодом, переданным в виде строки
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		panic(err)
	}
	// печатаем дерево
	ast.Print(fset, f)
}
