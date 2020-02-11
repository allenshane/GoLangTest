package main

import (
	"archive/zip"
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
	Compression struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
	} `json:"compression"`
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
		cmd := exec.Command("mysqldump", "-u"+username, "-p"+password, db)
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

func ZipWriter(source string, destination string) {
	fmt.Println("Starting Compression...")
	baseFolder := source

	// Get a Buffer to Write To
	outFile, err := os.Create(destination)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, baseFolder, "")

	if err != nil {
		fmt.Println(err)
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Compression Completed!!!")
}

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(basePath + file.Name())
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := basePath + file.Name() + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}

func main() {
	fmt.Println("Starting the app....")
	config, _ := LoadConfiguration(os.Args[1])
	CopyDirectories(config.DirecoriesToCopy.Source, config.DirecoriesToCopy.Destination)
	MysqlDump(config.Mysqldb.Username, config.Mysqldb.Password, config.Mysqldb.Dbs, config.Mysqldb.Destination)
	ZipWriter(config.Compression.Source, config.Compression.Destination)
}
