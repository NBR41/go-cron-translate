package main

import (
	"fmt"
	"log"
	"os"

	"github.com/NBR41/gocrontranslate/translator"
)

func main() {
	if len(os.Args) == 1 {
		log.Fatal(`usage: gocrontranslate "[CRON expr]"`)
		return
	}
	val, err := translator.GetTranslation(os.Args[1])

	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println(val)
}
