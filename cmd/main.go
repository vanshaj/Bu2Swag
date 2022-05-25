package main

import (
	"fmt"
	"os"
	"flag"
	"github.com/vanshaj/Bu2Swag/internal/app"
	"log"
)

var count int = 0

func main(){
	options := &app.Options{}
	flag.IntVar(&options.Thread,"thread", 1, "this is the number of concurrent threads you want to run for this process")
	flag.StringVar(&options.BurpFile, "filepath","D:\\TaskNotes\\CareApis\\LoggerPlusPlus.csv", "this is the csv file which will contain the burp history")
	flag.Parse()

	if len(options.BurpFile) == 0 { 
		fmt.Println("Burp file is not mentioned")
		os.Exit(1)
	}

	err := app.Run(options)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

}
