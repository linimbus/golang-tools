package main

import (
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	HELP      bool
	SRC_DIR   string
	Separator string
	HMAC_LIST map[string]string
)

func init() {
	HMAC_LIST = make(map[string]string, 1024)
	Separator = fmt.Sprintf("%c", os.PathSeparator)
	flag.BoolVar(&HELP, "help", false, "This help.")
	flag.StringVar(&SRC_DIR, "src", "", "The directory path to be searched.")
}

func hmacCalc(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("%s open failed! (%s)\n", filePath, err.Error())
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Printf("%s read failed! (%s)\n", filePath, err.Error())
		return "", err
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}

func matchFile(filename string, hmac string) bool {
	file, b := HMAC_LIST[hmac]
	if b {
		log.Printf("find match between %s <-> %s\n", file, filename)
		return true
	}
	HMAC_LIST[hmac] = filename
	return false
}

func searchedFile(srcdir string) {
	files, _ := ioutil.ReadDir(srcdir)
	for _, file := range files {
		if file.IsDir() {
			searchedFile(srcdir + Separator + file.Name())
		} else {
			filename := srcdir + Separator + file.Name()
			hmac, err := hmacCalc(filename)
			if err != nil {
				continue
			}
			// log.Printf("FILE: %s HMAC: %s\n", filename, hmac)

			if matchFile(filename, hmac) {
				err := os.Remove(filename)
				if err != nil {
					log.Printf("%s remove failed! (%s)\n", filename, err.Error())
				}
			}
		}
	}
}

func main() {

	flag.Parse()
	if HELP || SRC_DIR == "" {
		flag.Usage()
		os.Exit(-1)
	}

	searchedFile(SRC_DIR)
}
