package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var (
	HELP bool

	DIR       string
	Separator string
)

func init() {
	Separator = fmt.Sprintf("%c", os.PathSeparator)

	flag.BoolVar(&HELP, "help", false, "This help.")

	flag.StringVar(&DIR, "dir", "."+Separator, "directory path file rename.")
}

func RenameFile(filename, newfilename string) error {
	return os.Rename(filename, newfilename)
}

func FindFile(dir string) {
	files, _ := ioutil.ReadDir(dir)
	for _, file := range files {
		if file.IsDir() {
			FindFile(dir + Separator + file.Name())
		}
		filename := file.Name()
		newfilename := strings.ReplaceAll(filename, " ", "")
		if filename != newfilename {
			err := RenameFile(dir+Separator+filename, dir+Separator+newfilename)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Printf("change file name [%s] to [%s] success\n", filename, newfilename)
			}
		}
	}
}

func main() {
	flag.Parse()
	if HELP {
		flag.Usage()
		os.Exit(-1)
	}

	FindFile(DIR)
}
