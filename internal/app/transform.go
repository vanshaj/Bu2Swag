package app

import (
	"bufio"
	"strings"
	"encoding/json"
	"gopkg.in/yaml.v3"
	"github.com/Jeffail/gabs/v2"
	"fmt"
	"log"
	"os"
)

func jsonToYaml(jsonString []byte, filePath string) (err error){
	var jsonObj interface{}

	err = json.Unmarshal(jsonString, &jsonObj)

	if err != nil {
		return err
	}

	yamlOut , err := yaml.Marshal(jsonObj)
	if err != nil{
		return err
	}
	currentDir, err := os.Getwd()
	if err != nil{
		return err
	}
	filePath = fmt.Sprintf("%s/swagger.yaml",currentDir)
	fYml, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer fYml.Close()
	
	scanner := bufio.NewScanner(strings.NewReader(string(yamlOut)))
	fYml.WriteString("\n")
	for scanner.Scan(){
		text := scanner.Text()
		data := fmt.Sprintf("  %s\n",text)
		fYml.WriteString(data)
	}

	return nil
	//reOturn yamlOut, err
}

func jsonToSwagger(jsonString []byte) (out []byte, err error){
	var unmarshalledJson map[string]interface{}
	json.Unmarshal(jsonString, &unmarshalledJson)
	var gabObj *gabs.Container = gabs.New()
	gabMainObj , _ := gabObj.Object("properties")
	fmt.Println(unmarshalledJson)
	iterate(unmarshalledJson, gabMainObj)
	return []byte(gabObj.String()), nil
}

func iterate(jsonString map[string]interface{}, gabObj *gabs.Container ) {
	log.SetFlags(log.Llongfile)
	for key, value := range jsonString {
		switch v := value.(type) {
		case bool:
			gabBoolObject, err := gabObj.Object(key)
			if err != nil{
				log.Fatal(err)
				return
			}
			gabBoolObject.Set("boolean","type")
		case float64:
			gabFloat64Object, err := gabObj.Object(key)
			if err != nil{
				log.Fatal(err)
				return
			}
			gabFloat64Object.Set("integer","type")
		case map[string]interface{}:
			gabMapObject, err := gabObj.Object(key)
			if err != nil{
				log.Fatal(err)
				return
			}
			gabMapObject.Set("object","type")
			gabMapProperties, err := gabMapObject.Object("properties")
			iterate(v, gabMapProperties)
		case []interface{}:
			gabArrayObjectM, err := gabObj.Object(key)
			if err != nil{
				log.Fatal(err)
				return
			}
			gabArrayObjectM.Set("array", "type")
			gabItemsObjectM, err := gabArrayObjectM.Object("items")
			if err != nil{
				log.Fatal(err)
				return
			}
			typeOfArray := fmt.Sprintf("%T",v[0])
			if typeOfArray == "float64" {
				gabItemsObjectM.Set("integer","type")
			} else if typeOfArray == "bool"{
				gabItemsObjectM.Set("boolean","type")
			} else if typeOfArray == "string" {
				 gabItemsObjectM.Set("string", "type")
			} else{
				gabItemsObjectM.SetP("object", "type")
				gabArrayProperties, _ := gabItemsObjectM.Object("properties")
				iterate(v[0].(map[string]interface{}) , gabArrayProperties)
			}
		default:
			gabStringObject, err := gabObj.Object(key)
			if err != nil{
				log.Fatal(err)
				return
			}
			gabStringObject.Set(fmt.Sprintf("%T",v), "type")
		}
	}
}
