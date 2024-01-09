package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func OpenFile(filename string) ([]byte, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	defer fd.Close()

	var body []byte
	for {
		var tmp [128]byte
		cnt, err := fd.Read(tmp[:])
		if err != nil {
			if io.EOF == err {
				break
			}
			return nil, err
		}
		body = append(body, tmp[:cnt]...)
	}

	return body, nil
}

func SaveFile(filename string, body []byte) error {

	var fd *os.File
	var err error

	for {
		fd, err = os.Create(filename)
		if err == nil {
			break
		}
		if os.ErrExist != err {
			return err
		}
		err = os.Remove(filename)
		if err != nil {
			return err
		}
	}

	defer fd.Close()

	cnt := 0
	for {
		num, err := fd.Write(body[cnt:])
		if err != nil {
			return err
		}
		cnt += num
		if cnt == len(body) {
			break
		}
	}

	return nil
}

func Format(input string) string {

	input = strings.ToLower(input)

	inputbody := []byte(input)

	output := make([]byte, 0)

	for _, v := range inputbody {

		if v == ' ' {
			output = append(output, '-')
		} else if v == '.' {
			continue
		} else if v == '(' || v == ')' {
			continue
		} else {
			output = append(output, v)
		}
	}

	return string(output)
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage : <input.md> ")
		return
	}

	input := os.Args[1]
	output := os.Args[1]

	body, err := OpenFile(input)
	if err != nil {
		log.Println(err.Error())
		return
	}

	log.Println("read file success!")

	linelist := strings.Split(string(body[:]), "\r\n")

	var outbody string

	for idx, line := range linelist {

		enter := "\r\n"
		if idx+1 == len(linelist) {
			enter = ""
		}

		if len(line) == 0 {
			outbody += fmt.Sprintf(enter)
			continue
		}

		begin := strings.Index(line, "(#")
		if begin == -1 {
			outbody += fmt.Sprintf("%s%s", line, enter)
			continue
		}

		begin += 2
		end := strings.Index(line[begin:], ")")
		if end == -1 {
			outbody += fmt.Sprintf("%s%s", line, enter)
			continue
		}
		end += begin

		if end == begin {
			outbody += fmt.Sprintf("%s%s", line, enter)

			log.Println(line)
			continue
		}

		fmt.Printf("%s\r\n", line)
		fmt.Println(begin, end)

		outbody += fmt.Sprintf("%s%s%s%s", line[:begin], Format(line[begin:end]), line[end:], enter)
	}

	err = SaveFile(output, []byte(outbody))
	if err != nil {
		log.Println(err.Error())
		return
	}
}
