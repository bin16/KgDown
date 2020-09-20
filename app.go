package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

func init() {
	nameMap = make(map[string]string)
	conf := config{}
	toml.DecodeFile("config.toml", &conf)
	if conf.Names != nil {
		nameMap = conf.Names
	}
}

func main() {
	urls := []string{}
	a1 := os.Args[1]
	if a1 == "-i" {
		srcPath := os.Args[2]
		srcFile, err := ioutil.ReadFile(srcPath)
		if err != nil {
			log.Fatalf("Failed to read %s: %v", srcPath, err)
		}

		rows := bytes.Split(srcFile, []byte("\n"))
		for _, r := range rows {
			s := string(r)
			if strings.HasPrefix(s, "http") {
				urls = append(urls, s)
			}
		}
	} else {
		urls = append(urls, a1)
	}

	for i, r := range urls {
		fmt.Printf("%d - %s\n", i+1, r)
		for t := 0; t < 5; t++ {
			time.Sleep(time.Second * 2 << t)
			mp3Path, err := downloadMP3(r)
			if err != nil {
				fmt.Printf("Failed: %v\n", err)
			} else {
				fmt.Printf("Ok: %s\n", mp3Path)
				break
			}
		}
	}
}
