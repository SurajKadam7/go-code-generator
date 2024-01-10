package generator

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/vetcher/go-astra"
	"github.com/vetcher/go-astra/types"
)

func parseFile(fpath string) (res types.File, err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return
	}
	path := filepath.Join(currentDir, fpath)
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments|parser.AllErrors)
	if err != nil {
		return
	}
	file, err := astra.ParseAstFile(f)
	if err != nil {
		fmt.Println(err)
	}
	return *file, nil
}

// TODO implimenation for the other types as well
func findTheType(arg types.Type) string {
	switch tt := arg.(type) {
	case types.TMap:
		fmt.Println("TMap")
	case types.TName:
		return tt.TypeName
	case types.TPointer:
		fmt.Println("TPointer")
	case types.TImport:
		return fmt.Sprintf("%s.", tt.Import.Name) + findTheType(tt.NextType())
	case types.TEllipsis:
		fmt.Println("TEllipsis")
	case types.TChan:
		fmt.Println("TChan")
	}
	return ""
}

func parseVariable(args []types.Variable) []Arg {
	arguments := []Arg{}
	for _, arg := range args {
		name := arg.Name
		dtype := findTheType(arg.Type)

		arguments = append(arguments, Arg{
			Name:  name,
			Dtype: dtype,
		})
	}

	return arguments
}

func parseImports(imports []*types.Import) []string {
	result := []string{KafkaPackage, Json}
	for _, imp := range imports {
		pkg := strings.Split(imp.Package, "/")
		modiFiedImp := ""
		if len(pkg) > 0 && pkg[len(pkg)-1] != imp.Name {
			modiFiedImp += imp.Name + " "
		}
		modiFiedImp += imp.Package
		result = append(result, modiFiedImp)
	}

	// adding all the packages under double inverted comma
	for i, val := range result {
		result[i] = fmt.Sprintf("\"%s\"", val)
	}
	return result
}
