package app

import (
	"fmt"
	"log"
	"encoding/base64"
	"github.com/Jeffail/gabs/v2"
	"strings"
	"regexp"
	"os"
	"bufio"
)


func base64ToSwaggerJsonAndBodyJson(base64Value []byte) (swaggerJsonData []byte){
    jsonObject := gabs.New()
    var gabDataContainer *gabs.Container
    var gabMainContainer *gabs.Container
    var path string
    var methodType string
    var contentType string


//  var jsonBody string
    reachedToBody := false
    dst := make([]byte, base64.StdEncoding.DecodedLen(len(base64Value)))
    n, err := base64.StdEncoding.Decode(dst, base64Value)
    if err != nil {
            log.Fatal(err)
            return
    }
    dstReader := strings.NewReader(string(dst[:n]))
    scanner := bufio.NewScanner(dstReader)
    for scanner.Scan(){
            text := scanner.Text()
            if reachedToBody {
                            if len(text) > 0{
				params := gabs.New()
				if contentType == "application/x-www-form-urlencoded" {
					applicationMap := make(map[string]bool)
					textArray := strings.Split(text, "&")
					for _, val := range textArray {
						eachData := strings.Split(val, "=")
						if _, ok := applicationMap[eachData[0]]; ok{
							continue
						}
						applicationMap[eachData[0]] = true
						params.SetP(eachData[0], "name")
						params.SetP("formData","in")
						params.SetP("string","type")
						gabDataContainer.ArrayAppend(params, "parameters")
					}
				}else {
					swag , _ := jsonToSwagger([]byte(text))
					parseSwaggerBytes, err := gabs.ParseJSON(swag)
					if err != nil {
						log.Fatal(err)
						os.Exit(1)
					}
					params.SetP("body", "name")
					params.SetP("body", "in")
					params.SetP(parseSwaggerBytes, "schema")
	
					gabDataContainer.ArrayAppend(params,"parameters")
			    }
            		}
		} else{
                    if strings.Contains(text, "HTTP"){
                            httpHeaderLine := strings.Split(text, " ")
                            methodType = httpHeaderLine[0]

                            pathWithQueryParams := strings.Split(httpHeaderLine[1],"?")
			    initpath := pathWithQueryParams[0]
			    
                            regexId := regexp.MustCompile(`\d+\w+`)
                            path = regexId.ReplaceAllString(initpath, "{id}")

			    key := fmt.Sprintf("%s [%s]", path, methodType) 
			
			    if _, ok := Endpoints[key]; ok {
				return nil
			    }

			    Endpoints[key] = true

                            gabMainContainer , _ = jsonObject.Object(key)
                            gabDataContainer , _ = gabMainContainer.Object(strings.ToLower(methodType))
			    gabDataContainer.Array("parameters")
			    
			    if len(pathWithQueryParams) > 1{
				queryParameterArray := strings.Split(pathWithQueryParams[1], "&")
				for _ , value := range queryParameterArray{
					queryParams := gabs.New()
					paramName := strings.Split(value, "=")[0]
					paramValue := strings.Split(value, "=")[1]
					queryParams.SetP(paramName, "name")
					queryParams.SetP("query", "in")
					queryParams.SetP( fmt.Sprintf("%T", paramValue), "type")

					gabDataContainer.ArrayAppend(queryParams, "parameters")
				}
			    }

                            
                            params := gabs.New()
                            if regexId.Match([]byte(initpath)) {
                                    params.SetP("id", "name")
                                    params.SetP("path", "in")
                                    params.SetP("string", "type")
				    params.SetP(true, "required")
                            	    gabDataContainer.ArrayAppend(params,"parameters")
                            }
		    }else if strings.Contains(text, "Content-Type:"){
			    regexConsumes := regexp.MustCompile(`: ([^;]+)`)
			    contentType = regexConsumes.FindStringSubmatch(text)[1]
                            gabDataContainer.ArrayAppend (contentType, "consumes")
                    }

            }

            if text == "" {
                    reachedToBody = true
            }
    }

    // Response addition
    responseObj := gabs.New()
    responseStatusObj, _ := responseObj.Object("200")
    responseStatusObj.SetP("Correct Data", "description")

    gabDataContainer.SetP(responseObj, "responses")

    return []byte(jsonObject.String())
}
