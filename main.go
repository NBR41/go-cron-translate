package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NBR41/go-cron-translate/translator"
)

func main() {
	val, err := translator.GetTranslation(os.Args[1])
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(val)
}
