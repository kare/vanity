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

func usage() {
	fmt.Fprintf(os.Stderr, "usage: vanity -d domain -c vanity.conf [-p 80]\n")
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetPrefix("vanity: ")
	log.SetFlags(0)

	var (
		domain   = flag.String("d", "", "http domain name")
		port     = flag.Int("p", 80, "http server port")
		confFile = flag.String("c", "", "configuration file")
	)

	if *domain == "" || *confFile == "" {
		usage()
		os.Exit(2)
	}

	file, err := os.Open(*confFile)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("error closing configuration file: %v", err)
		}
	}()
	conf, err := readConfig(file)
	if err != nil {
		log.Fatal(err)
	}
	server := newServer(*domain, conf)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", *port), server))
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
