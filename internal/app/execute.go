package app

import (
	"fmt"
	"os"
	"encoding/csv"
	"log"
	"io"
	//"github.com/Jeffail/gabs/v2"
)

var Endpoints map[string]bool

type Options struct{
	BurpFile string
	Thread int
}

func Run(options *Options) error {
	fmt.Printf("Path of file is %v\n", options.BurpFile)
	fmt.Printf("Number of Threads are  %v\n", options.Thread)
	
	Endpoints = make(map[string]bool)
	
	file, err := os.Open(options.BurpFile)
	if err != nil {
		log.Fatal(err)
	}
	

	defaultFile ,err := os.Open("D:\\GoProj\\Bu2Swag\\default.yaml")
      if err != nil{
              return err
      }
      defer defaultFile.Close()

      filePath := "D:\\GoProj\\Bu2Swag\\swagger.yaml"
      fYml, err := os.OpenFile(filePath,os.O_CREATE, 0600)
      if err != nil {
              return err
      }
        
      _, err = io.Copy(fYml, defaultFile)
      if err != nil{
              return err
      }
	fYml.Close()	
	
	reader := csv.NewReader(file)
	lineNum := 1
	for {
		eachRecord , err := reader.Read()
		if err != nil || err == io.EOF{
			break
		}
		if lineNum == 1 {
			lineNum++
			continue
		} else {
			swaggerJsonBytes:=  base64ToSwaggerJsonAndBodyJson([]byte(eachRecord[0]))
			if swaggerJsonBytes == nil {
				lineNum++
				continue
			}
			err = jsonToYaml(swaggerJsonBytes, "swagger.yaml")
			if err != nil{
				return err
			}		
		}
		lineNum++
	}


	return nil
}

