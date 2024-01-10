package generator

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"text/template"

	"go/format"

	"github.com/vetcher/go-astra/types"
)

const (
	KafkaPackage string = "kafka github.com/segmentio/kafka-go"
	Json         string = "encoding/json"
)

type Generator struct {
	InterFaceFile string
	OutPutFile    string
	TemplateFile  string
}

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

func (g Generator) Generate() {
	builder := prePareData(g.InterFaceFile)
	templateFile, err := os.Open(g.TemplateFile)
	if err != nil {
		log.Fatal(err)
	}
	defer templateFile.Close()

	templateBuffer, _ := io.ReadAll(templateFile)
	tmp, err := template.New("kafka_publisher").Parse(string(templateBuffer))
	if err != nil {
		log.Fatal(err)
	}

	buff := &bytes.Buffer{}
	err = tmp.Execute(buff, builder)
	if err != nil {
		log.Fatal(err)
	}

	formatFile(buff.Bytes(), g.OutPutFile)
}

func build(f types.File) (b Builder) {
	// validate the given file before accessing the content
	if err := validate(f); err != nil {
		log.Fatal(err)
	}

	intrFace := f.Interfaces[0]

	b.PackageName = f.Name
	b.InterfaceName = intrFace.Name
	b.Imports = parseImports(f.Imports)

	for _, method := range intrFace.Methods {
		args := parseVariable(method.Args)
		results := parseVariable(method.Results)
		b.Methods = append(b.Methods, Method{
			Name:   method.Name,
			Args:   args,
			Result: results,
		})
	}
	return
}

func validate(f types.File) error {
	if f.Interfaces == nil || len(f.Interfaces) == 0 {
		return errors.New("interface is not provided")
	}

	if len(f.Interfaces) > 1 {
		return errors.New("more than one interfaces are defined")
	}

	return nil
}

func prePareData(path string) Builder {
	file, err := parseFile(path)
	if err != nil {
		log.Fatal(err)
	}

	// TODO check for max interfaces to be allowed in give path file
	builder := build(file)
	fmt.Printf("%+v\n", builder)
	return builder
}

func formatFile(buff []byte, path string) {
	formatedBuff, err := format.Source(buff)
	if err != nil {
		log.Fatal(err)
	}

	file, _ := os.Create(path)
	file.Write(formatedBuff)
	file.Close()
}
