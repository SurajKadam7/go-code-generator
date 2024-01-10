package generator

import (
	"fmt"
	"go/parser"
	"go/token"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/vetcher/go-astra"
	"github.com/vetcher/go-astra/types"
)

const (
	KafkaPackage string = "github.com/segmentio/kafka-go"
)

type Arg struct {
	Name  string
	Dtype string
}

type Method struct {
	Name   string // name of the method
	Args   []Arg
	Result []Arg
}

type Builder struct {
	PackageName   string
	Imports       []string
	InterfaceName string
	Methods       []Method
}

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
	result := []string{KafkaPackage}
	for _, imp := range imports {
		pkg := strings.Split(imp.Package, "/")
		modiFiedImp := ""
		if len(pkg) > 0 && pkg[len(pkg)-1] != imp.Name {
			modiFiedImp += imp.Name + " "
		}
		modiFiedImp += imp.Package
		result = append(result, modiFiedImp)
	}
	for i, val := range result {
		// this for import formating purpose
		result[i] = fmt.Sprintf("\"%s\"", val)
	}
	return result
}

func build(f types.File) (build Builder) {
	// validate the given file before accessing the content
	if err := validate(f); err != nil {
		log.Fatal(err)
	}

	intrFace := f.Interfaces[0]

	build.PackageName = f.Name
	build.InterfaceName = intrFace.Name
	build.Imports = parseImports(f.Imports)

	for _, method := range intrFace.Methods {
		args := parseVariable(method.Args)
		results := parseVariable(method.Results)
		build.Methods = append(build.Methods, Method{
			Name:   method.Name,
			Args:   args,
			Result: results,
		})
	}
	return
}

func validate(f types.File) error {
	if f.Interfaces == nil {
		log.Fatal("interface is not provided")
	}

	if len(f.Interfaces) > 1 {
		log.Fatal("more than one interfaces are defined")
	}

	return nil
}

func prePareData(path string) Builder {
	file, err := parseFile(path)
	if err != nil {
		log.Fatal(err)
	}

	if len(file.Interfaces) == 0 {
		log.Fatal("interface is not present")
	}
	// TODO check for max interfaces to be allowed in give path file
	builder := build(file)
	fmt.Printf("%+v\n", builder)
	return builder
}

func Generate(path string, tmplPath string) {
	op, err := os.Create("kfk_pub")
	if err != nil {
		log.Fatal(err)
	}
	defer op.Close()

	builder := prePareData(path)
	ff, err := os.Open("./templates/publisher")
	if err != nil {
		log.Fatal(err)
	}
	defer ff.Close()

	byttmp, _ := io.ReadAll(ff)
	tmp, err := template.New("kafka_publisher").Parse(string(byttmp))
	if err != nil {
		log.Fatal(err)
	}

	err = tmp.Execute(op, builder)
	if err != nil {
		log.Fatal(err)
	}
}
