package main 

import (	"fmt"
		"io/ioutil"
		"strings"
		"os"
		"regexp"
		"strconv"
		)


type Json struct {
	 Bool bool
	 String string
	 Number int
	 Float float64 
	 Array []Json
	 Object map[string]Json
	 Parsed bool
	 Null interface{}
	 Type string
}


var re = regexp.MustCompile("^[-+]?[0-9]*.?[0-9]+([eE][-+]?[0-9]+)?")

func stringParser(JsonData string) (Json, string) {
	//This function parses string elements and returns the string and remaining Jsondata
	var result Json
	var index int                          // stores the index of end of string
	JsonData = strings.TrimSpace(JsonData) // Trim the data for spaces
	if JsonData[0] == '"' {  
		JsonData = JsonData[1:]
		for JsonData[index] != '"' { //looping over till the end of string
			index++
		}
		return Json{String: JsonData[:index], Parsed: true, Type: "String"}, JsonData[index+1:]
	}
	return result, JsonData
}


func numberParser(JsonData string) (Json, string) {
	//Parses numbers in Jsonstring and returns the number and remaining JsonData
	JsonData = strings.Trim(JsonData, " ")
	var number Json
	var result [][]int
	result1 := re.FindAllString(JsonData, -1)                          //Library function to find occurrence of string and returns a slice of succe elem
	if len(result1) > 0 {
		result = re.FindAllStringIndex(JsonData, -1) //Gives start and end index of the string
		index := result[0][1] //index is the end  index of the number
		j,err := strconv.Atoi(result1[0])
		if err == nil{
			return Json{Number: j, Parsed: true, Type: "Number"}, JsonData[index:]
		}
		i,err := strconv.ParseFloat(result1[0],64)
		if err == nil {
			return Json{Float: i,Parsed: true, Type: "Float"}, JsonData[index:]
		}
	} else {
		return number, JsonData
	}
	return number,JsonData
}

func arrayParser(JsonData string) (Json,string) {
	//Parses Json arrays and stores them into the user defined type else returns the data
	var result Json
	parsed := make([]Json,0)
	if JsonData[0] == '[' {
		JsonData = JsonData[1:]
		for len(JsonData) > 0 {
			result, JsonData = elementParser(JsonData)
			if result.Parsed == true {
				parsed = append(parsed,result)
				JsonData = commaParser(JsonData)
			}
			if JsonData[0] == ']'{
				fmt.Println("End of array")
				return Json{Array: parsed, Parsed: true, Type: "Array"},JsonData[1:]
			}
		}
	}
	return result,JsonData
}

func objParser(JsonData string) (Json,string) {
	//Parses Json arrays and stores them into the user defined type else returns the data
	var result Json
	var key string
	parsed := make(map[string]Json)
	if JsonData[0] == '{' {
		JsonData = JsonData[1:]
		for len(JsonData) > 0 {
			result,JsonData = stringParser(JsonData)
			if result.Parsed == true {
				key = result.String
			}
			JsonData = colonParser(JsonData)
			result, JsonData = elementParser(JsonData)
			if result.Parsed == true {
				parsed[key] = result
				JsonData = commaParser(JsonData)
			}
			if JsonData[0] == '}'{
				fmt.Println("End of Object")
				return Json{Object: parsed, Parsed: true, Type: "Object"},JsonData[1:]
			}
		}
	}
	return result,JsonData
}

func colonParser(JsonData string) string {
	//Parses colon in Json objects and returns the colon and the remaining JsonData
	JsonData = strings.Trim(JsonData, " ")
	if JsonData[0] == ':' {
		return JsonData[1:]
	}
	return JsonData
}

func commaParser(JsonData string) string {
	//Parses comma which saperate elements and returns comma and the remaining JsonData
	JsonData = strings.Trim(JsonData, " ")
	if JsonData[0] == ',' {
		return JsonData[1:]
	}
	return JsonData
}

func booleanParser(JsonData string) (Json, string) {
	// Function parses boolean elements and returns the bool value and the remaining Jsondata
	var result Json
	JsonData = strings.Trim(JsonData, " ")
	if len(JsonData) > 4 {
		if JsonData[0:4] == "true" {
			return Json{Bool: true, Parsed: true, Type: "Bool"}, JsonData[4:] //slicing the JsonData
		} else if JsonData[0:5] == "false" {
			return Json{Bool: false, Parsed: true ,Type:"Bool"}, JsonData[5:]
		}
	}
	return result, JsonData
}

func nullParser(JsonData string)  (Json ,string) {
	//Parses null values form Jsonstring and returns the nil value and remaining JsonData
	JsonData = strings.Trim(JsonData, " ")
	var result Json
	if len(JsonData) > 4 {
		if JsonData[0:4] == "null" {
			return Json{Null:nil, Parsed: true, Type: "Null"}, JsonData[4:]
		}
	}
	return result, JsonData
}


func elementParser(JsonData string) (Json, string) {
	//Function tries all the parsers one by one on each element and returns the result and the remaining JsonData
	var result Json
	result, JsonData = stringParser(JsonData)
	if result.Parsed == true {
		return result, JsonData
	}
	result, JsonData = numberParser(JsonData)
	if result.Parsed == true {
		return result, JsonData
	}
	
	result, JsonData = booleanParser(JsonData)
	if result.Parsed == true {
		return result, JsonData
	}
	result, JsonData = nullParser(JsonData)
	if result.Parsed == true {
		return result, JsonData
	}
	result, JsonData = arrayParser(JsonData)
	if result.Parsed == true {
		return result,JsonData
	}
	result, JsonData = objParser(JsonData)
	if result.Parsed == true {
		return result,JsonData
	}
	return result, JsonData
}

//Utility function to acess the parsed json
func (m Json) getElement() interface{} {
	typeof := m.Type
	switch typeof {
		case "String":
			return m.String
		case "Object":
			return m.Object
		case "Number":
			return m.Number
		case "Array":
			return m.Array
		case "Null":
			return m.Null
		case "Float":
			return m.Float
		case "Bool":
			return m.Bool
		default:
			fmt.Println("Don't know how to handle the type")
			os.Exit(1)
	}
	return ""
}

func main(){
	buf, err := ioutil.ReadFile("test.txt") //Reading the file completely
	if err != nil {                          //error check
		fmt.Println(err)
		os.Exit(1)
	}
	data := string(buf) //bytes to string
	if data[0] == '['{
		result,_ :=arrayParser(data)
		final := result.getElement()
		fmt.Println(final)          
	
	}else if data[0] =='{' {
		result,_ :=objParser(data)
		fmt.Println(result)
		k := result.getElement()
		fmt.Println(k)		
	}
	os.Exit(0)
}
