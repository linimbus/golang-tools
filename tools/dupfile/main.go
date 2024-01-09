package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

var (
	HELP bool
	SRC_DIR  string
	Separator string
	HMAC_LIST map[string]string
)

func init()  {
	HMAC_LIST = make(map[string]string, 1024)
	Separator = fmt.Sprintf("%c", os.PathSeparator)
	flag.BoolVar(&HELP,"help",false,"This help.")
	flag.StringVar(&SRC_DIR,"src", "","The directory path to be searched.")
}

func hmacCalc(filename string) string {
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("%s open failed! (%s)\n",filename, err.Error())
		return ""
	}
	h := hmac.New(sha256.New, []byte("test"))
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

func matchFile(filename string, hmacvalue string) bool {
	file, b := HMAC_LIST[hmacvalue]
	if b == true {
		log.Printf("find match between %s <-> %s\n", file, filename)
		return true
	}
	HMAC_LIST[hmacvalue] = filename
	return false
}

func searchedFile(srcdir string) {
	files, _ := ioutil.ReadDir(srcdir)
	for _, file := range files {
		if file.IsDir() {
			searchedFile(srcdir + Separator + file.Name())
		} else {
			filename := srcdir + Separator + file.Name()
			hmac := hmacCalc(filename)
			if hmac == "" {
				continue
			}
			if matchFile(filename, hmac) == true {
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
	if HELP || SRC_DIR == ""{
		flag.Usage()
		os.Exit(-1)
	}

	searchedFile(SRC_DIR)
}
