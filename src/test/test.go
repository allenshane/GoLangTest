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
	DirecoriesToCopy struct {
		Source      []string `json:"source"`
		Destination string   `json:"destination"`
	} `json:"directoriestocopy"`
	Mysqldb struct {
		Username    string   `json:"username"`
		Password    string   `json:"password"`
		Dbs         []string `json:"dbs"`
		Destination string   `json:"destination"`
	} `json:"mysqldb"`
}

func LoadConfiguration(file string) (config Config, err error) {
	fmt.Println("Loading Config....")
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		return
	}
	dec := json.NewDecoder(configFile)
	err = dec.Decode(&config)
	return
}

func CopyDirectories(sources []string, destination string) {
	fmt.Println("Copying Directories....")
	for i, source := range sources {
		fmt.Printf("Copying [%d] ==> [%s] to [%s]\n", i, source, destination)
		copy.CopyDirectory(source, destination)

	}
	fmt.Println("Copying Completed!!!")
}

func MysqlDump(username string, password string, dbs []string, destination string) {
	os.Mkdir(destination, 0700)
	fmt.Printf("Starting MySQL Dump")
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
		result := destination + "/" + db + ".sql"
		err = ioutil.WriteFile(result, bytes, 0644)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("MySQL Dump Completed")
}

func main() {
	fmt.Println("Starting the app....")
	config, _ := LoadConfiguration("E:/go/config.json")
	CopyDirectories(config.DirecoriesToCopy.Source, config.DirecoriesToCopy.Destination)
	MysqlDump(config.Mysqldb.Username, config.Mysqldb.Password, config.Mysqldb.Dbs, config.Mysqldb.Destination)
}
