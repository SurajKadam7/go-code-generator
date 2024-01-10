package main

import (
	"github.com/SurajKadam7/go-auto/generator"
)

func main() {
	opPath := "generated_templates/publisher"	
	tmplPath := "templates/publisher"
	interFacePath := "sample_interfaces/service.go"
	gen := generator.Generator{
		OutPutFile: opPath,
		TemplateFile: tmplPath,
		InterFaceFile: interFacePath,
	}

	gen.Generate()
}


