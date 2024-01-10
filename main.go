package main

import (
	"github.com/SurajKadam7/go-auto/generator"
)

func main() {
	path := "/Users/suraj.kadam/go/src/bitbucket.org/junglee_games/personal/go-auto/templates/publisher"
	generator.Generate("./test/service.go", path)
}

// currentDir, err := os.Getwd()
// if err != nil {
// 	panic(err)
// }
// path := filepath.Join(currentDir, "./test/service.go")
// fset := token.NewFileSet()
// f, err := parser.ParseFile(fset, path, nil, parser.ParseComments|parser.AllErrors)
// if err != nil {
// 	panic(fmt.Errorf("error when parse file: %v", err))
// }
// file, err := astra.ParseAstFile(f)
// if err != nil {
// 	fmt.Println(err)
// }
// var fileType types.File

// fmt.Println(fileType)
// t, err := json.Marshal(file)
// if err != nil {
// 	fmt.Println(err)
// }
// fmt.Println(string(t))
