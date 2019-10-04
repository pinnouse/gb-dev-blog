package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	//"google.golang.org/appengine"
)

// configuration : configuration for behind the scenes information
type configuration struct {
	Server   string `json:"server"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func getConfig() map[string]interface{} {
	file, err := os.Open("config.json")
	check(err)

	defer file.Close()

	byteValue, _ := ioutil.ReadAll(file)

	var config map[string]interface{}
	json.Unmarshal([]byte(byteValue), &config)

	return config
}

func main() {
	config := getConfig()
	fmt.Println("Read config file")

	fmt.Println("Connecting to database")
	db, err := connect(config["username"].(string), config["password"].(string), config["server"].(string))
	check(err)

	fmt.Println("Database connected")

	defer db.Close()

	tableExists := checkTable(db)
	if tableExists {
		getPosts(db)
	}
	serve(int(config["port"].(float64)), db)
	//appengine.Main()
}
