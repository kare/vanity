package main // import "kkn.fi/cmd/vanity"

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	domainFlag = flag.String("d", "", "http domain name")
	portFlag   = flag.Int("p", 80, "http server port")
	confFlag   = flag.String("c", "", "configuration file")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: vanity -d domain -c vanity.conf [-p 80]\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetPrefix("vanity: ")
	log.SetFlags(0)

	if *domainFlag == "" || *confFlag == "" {
		usage()
		os.Exit(2)
	}

	c, err := os.Open(*confFlag)
	if err != nil {
		log.Fatal(err)
	}
	conf, err := readConfig(c)
	if err != nil {
		log.Fatal(err)
	}
	server := newServer(*domainFlag, conf)
	port := fmt.Sprintf(":%v", *portFlag)
	log.Fatal(http.ListenAndServe(port, server))
}

func readConfig(r io.Reader) (map[string]*packageConfig, error) {
	conf := make(map[string]*packageConfig)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		switch len(fields) {
		case 0:
			continue
		case 3:
			path := fields[0]
			pack := newPackage(path, fields[1], fields[2])
			conf[path] = pack
		default:
			return conf, errors.New("configuration error: " + scanner.Text())
		}
	}
	return conf, nil
}
