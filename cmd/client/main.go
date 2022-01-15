package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	file := flag.String("file", "", "A template file to be processed.")
	varsFile := flag.String("vars-file", "", "A YAML file that contains variables to substitute into the template file.")

	flag.Parse()

	if *file == "" {
		log.Fatal("A template file must be specified with the --file parameter.")
	}

	fmt.Println(*varsFile)

	flag.Parse()
}
