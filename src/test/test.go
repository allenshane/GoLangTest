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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
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
	S3Info struct {
		Bucket    string `json:"bucket"`
		Region    string `json:"region"`
		AccessKey string `json:"accesskey"`
		SecretKey string `json:"secretkey"`
	} `json:"s3info"`
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

	outFile, err := os.Create(destination)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	w := zip.NewWriter(outFile)

	AddFiles(w, baseFolder, "")

	if err != nil {
		fmt.Println(err)
	}

	err = w.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Compression Completed!!!")
}

func AddFiles(w *zip.Writer, basePath, baseInZip string) {
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

			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			newBase := basePath + file.Name() + "/"
			fmt.Println("Recursing and Adding SubDir: " + file.Name())
			fmt.Println("Recursing and Adding SubDir: " + newBase)

			AddFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}

func UploadToAws(filepath string, bucket string, region string, accesskey string, secretkey string) {
	filename := filepath

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accesskey, secretkey, "")},
	)

	uploader := s3manager.NewUploader(sess)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
		Body:   file,
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Successfully uploaded %q\n", filename)
}

func main() {
	fmt.Println("Starting the app....")
	config, _ := LoadConfiguration(os.Args[1])
	CopyDirectories(config.DirecoriesToCopy.Source, config.DirecoriesToCopy.Destination)
	MysqlDump(config.Mysqldb.Username, config.Mysqldb.Password, config.Mysqldb.Dbs, config.Mysqldb.Destination)
	ZipWriter(config.Compression.Source, config.Compression.Destination)
	UploadToAws(config.Compression.Destination, config.S3Info.Bucket, config.S3Info.Region, config.S3Info.AccessKey, config.S3Info.SecretKey)
}
