package cmd

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gni.dev/gni/cmd/generator"
)

func Gen(args []string) {
	genCommand := flag.NewFlagSet("gen", flag.ExitOnError)
	javaFlag := genCommand.String("java", "", "Generate java source files.")

	if err := genCommand.Parse(args); err != nil {
		log.Fatal(err)
	}

	if len(genCommand.Args()) < 1 {
		log.Fatal("no interface files specified")
	}

	src := genCommand.Args()[0]
	srcPath := filepath.Dir(src)
	srcBase := filepath.Base(src)
	srcExt := filepath.Ext(srcBase)
	srcName := strings.TrimSuffix(srcBase, srcExt)

	srcf, err := os.Open(src)
	if err != nil {
		log.Fatal(err)
	}
	defer srcf.Close()

	d, err := generator.Parse(srcf)
	if err != nil {
		log.Fatal(err)
	}

	dstf, err := os.Create(filepath.Join(srcPath, srcName+".gni."+srcExt))
	if err != nil {
		log.Fatal(err)
	}
	defer dstf.Close()

	if err := generator.Golang(d, dstf); err != nil {
		log.Fatal(err)
	}

	if javaFlag != nil {
		dstf, err := os.Create(filepath.Join(*javaFlag, srcName+".java"))
		if err != nil {
			log.Fatal(err)
		}
		defer dstf.Close()

		if err := generator.Java(d, dstf); err != nil {
			log.Fatal(err)
		}
	}
}
