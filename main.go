package main

import (
	"os"
	"io/ioutil"
	"fmt"
	"encoding/json"
	//"google.golang.org/appengine"
)

// configuration : configuration for behind the scenes information
type configuration struct {
	Server string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port int `json:"port"`
}

func check(e error) {
	if e != nil {
		panic (e)
	}
}

func main() {
	file, err := os.Open("config.json")
	check(err)

	fmt.Println("Read config file")

	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var config map[string]interface{}
	json.Unmarshal([]byte(byteValue), &config)
	
	fmt.Println("Connecting to database")
	db, err := connect(config["username"].(string), config["password"].(string), config["server"].(string))
	check(err)
	
	fmt.Println("Database connected")
	
	
	defer db.Close()
	
	tableExists := checkTable(db)
	if tableExists {
		getPosts(db)
		serve(int(config["port"].(float64)))
	}
	//appengine.Main()
}