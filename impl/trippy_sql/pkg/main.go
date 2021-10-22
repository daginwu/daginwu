package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"trippy_sql/pkg/parser"

	"github.com/davecgh/go-spew/spew"
)

func main() {

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Trippy SQL >")
		query, _ := reader.ReadString('\n')
		fmt.Println(query)
		ast, err := parser.Parse(query)
		if err != nil {
			log.Println(err)
		}
		spew.Dump(ast)
	}

}
