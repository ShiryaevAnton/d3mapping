package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ShiryaevAnton/d3mapping/config"
	"github.com/ShiryaevAnton/d3mapping/d3map"
	"github.com/pelletier/go-toml"
)

const (
	search        = "s"
	pull          = "p"
	comp          = "c"
	compAndSearch = "cs"
)

func main() {

	var simplPath string
	var configPath string
	var d3Path string
	var mode string

	var c config.Config

	fConfig, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	toml.Unmarshal(fConfig, &c)

	flag.StringVar(&configPath, "cp", "", "Path to config file")
	flag.StringVar(&simplPath, "sp", "", "Path to simpl file")
	flag.StringVar(&d3Path, "dp", "", "Path to d3 simpl program")
	flag.StringVar(&mode, "m", "", "Mode: s - searching and replace signal. p - get rooms amd signals name from d3")

	flag.Parse()

	logName := "log" + strconv.FormatInt(time.Now().UTC().Unix(), 10) + ".txt"

	_, err = os.Stat("log")
	if os.IsNotExist(err) {
		err := os.Mkdir("log\\", 0666)
		if err != nil {
			log.Fatal("Folder does not exist.")
		}
	}

	fLog, err := os.OpenFile("log\\"+logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer fLog.Close()

	wrt := io.MultiWriter(os.Stdout, fLog)
	log.SetOutput(wrt)

	if mode == search {
		if err := d3map.Searching(configPath, simplPath, c); err != nil {
			log.Fatalln(err)
		}
	}
	if mode == pull {
		if err := d3map.GetNames(configPath, d3Path, c); err != nil {
			log.Fatalln(err)
		}
	}

	if mode == comp {
		if err := d3map.Comparing(configPath, c); err != nil {
			log.Fatalln(err)
		}
	}

	if mode == compAndSearch {
		if err := d3map.Comparing(configPath, c); err != nil {
			log.Fatalln(err)
		}
		if err := d3map.Searching(configPath, simplPath, c); err != nil {
			log.Fatalln(err)
		}
	}

}
