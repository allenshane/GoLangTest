package main

import (
	"encoding/json"
	"fmt"
	"github/copy"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Config struct {
	Folders struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	} `json:"folders"`
	Mysqldb struct {
		Username string   `json:"username"`
		Password string   `json:"password"`
		Dbs      []string `json:"dbs"`
	} `json:"mysqldb"`
}

func LoadConfiguration(file string) (config Config, err error) {
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return
	}
	dec := json.NewDecoder(configFile)
	err = dec.Decode(&config)
	return
}

func MysqlDump(username string, password string, dbs []string) {

	for i, db := range dbs {
		fmt.Printf("Dumping [%d] ==> [%s]\n", i, db)
		cmd := exec.Command("mysqldump", "-P3306", "-h127.0.0.1", "-uroot", "-ppassword", db)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal(err)
		}

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}

		bytes, err := ioutil.ReadAll(stdout)
		if err != nil {
			log.Fatal(err)
		}
		err = ioutil.WriteFile("./out.sql", bytes, 0644)
		if err != nil {
			panic(err)
		}
	}
}

func main() {
	fmt.Println("Starting the app....")
	config, _ := LoadConfiguration("E:/go/config.json")
	copy.CopyDirectory(config.Folders.Source, config.Folders.Destination)
	MysqlDump(config.Mysqldb.Username, config.Mysqldb.Password, config.Mysqldb.Dbs)
}
