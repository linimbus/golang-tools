package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"time"
)

var (
	Output string
	Help   bool
)

type Interface struct {
	Name  string   `json:name`
	Index int      `json:index`
	Flag  string   `json:flag`
	Mac   string   `json:mac`
	Mtu   int      `json:mtu`
	IP    []string `json:ip`
}

func init() {
	flag.StringVar(&Output, "output", "ipinfo.json", "ip lookup output filename")
	flag.BoolVar(&Help, "help", false, "usage help")
}

func main() {
	flag.Parse()
	if Help {
		flag.Usage()
		return
	}

	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime | log.Lmicroseconds)

	var latest []byte
	var count int

	for {
		if count > 0 {
			time.Sleep(time.Minute)
		}
		count++

		ifaces, err := net.Interfaces()
		if err != nil {
			log.Println(err.Error())
			continue
		}

		lookup := make([]Interface, 0)

		for _, iface := range ifaces {
			var info Interface

			info.Name = iface.Name
			info.Index = iface.Index
			info.Mac = iface.HardwareAddr.String()
			info.Mtu = iface.MTU
			info.Flag = iface.Flags.String()

			address, err := iface.Addrs()
			if err != nil {
				log.Println(err.Error())
			} else {
				for _, ip := range address {
					info.IP = append(info.IP, ip.String())
				}
			}
			lookup = append(lookup, info)
		}

		body, err := json.MarshalIndent(lookup, "", "\t")
		if err != nil {
			log.Println(err.Error())
			continue
		}

		if bytes.Equal(body, latest) {
			continue
		}

		latest = body
		err = ioutil.WriteFile(Output, body, 0644)
		if err != nil {
			log.Println(err.Error())
			continue
		}

		log.Println("ip lookup success")
	}
}
