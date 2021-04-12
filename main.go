package main

import (
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ShiryaevAnton/d3mapping/config"
	"github.com/ShiryaevAnton/d3mapping/d3map"
	"github.com/pelletier/go-toml"
)

var simplPath string
var configPath string

var c config.Config

func init() {

	fConfig, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	toml.Unmarshal(fConfig, &c)

	flag.StringVar(&configPath, "cp", "", "Path to config file")
	flag.StringVar(&simplPath, "sp", "", "Path to simpl file")

	flag.Parse()
}

func main() {

	logName := "log" + strconv.FormatInt(time.Now().UTC().Unix(), 10) + ".txt"

	fLog, err := os.OpenFile("log\\"+logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer fLog.Close()

	wrt := io.MultiWriter(os.Stdout, fLog)
	log.SetOutput(wrt)

	f, err := excelize.OpenFile(configPath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	fSimpl, err := os.ReadFile(simplPath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	var signalLightDimmerMap = make(map[string]string)
	var signalLightMap = make(map[string]string)

	signalLightMap[c.KeypanelLight.On] = c.CoreLight.On
	signalLightMap[c.KeypanelLight.OnFb] = c.CoreLight.OnFb
	signalLightMap[c.KeypanelLight.Off] = c.CoreLight.Off
	signalLightMap[c.KeypanelLight.OffFb] = c.CoreLight.OffFb

	signalLightDimmerMap[c.KeypanelLight.Raise] = c.CoreLight.Raise
	signalLightDimmerMap[c.KeypanelLight.Dim] = c.CoreLight.Dim
	signalLightDimmerMap[c.KeypanelLight.Level] = c.CoreLight.Level
	signalLightDimmerMap[c.KeypanelLight.LevelFb] = c.CoreLight.LevelFb

	simplPathNew := strings.ReplaceAll(simplPath, ".smw", "") + "_UPDATE.smw"

	var d3List []d3map.D3List

	log.Println("\nConfig file: " + strings.ReplaceAll(configPath, "./", "") + " contains:\n")

	for i := c.SheetConfig.StartRow; true; i++ {

		room := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnRoom+strconv.Itoa(i))
		if room == "" {
			break
		}
		light := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnLight+strconv.Itoa(i))
		lightType := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnLightType+strconv.Itoa(i))
		shade := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnShade+strconv.Itoa(i))
		shadeType := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnShadeType+strconv.Itoa(i))

		newd3List := d3map.NewD3List(room, i-c.SheetConfig.StartRow+1, light, lightType, shade, shadeType)

		log.Println(newd3List)

		d3List = append(d3List, newd3List)
	}

	log.Println("Start to find and replace signals")

	resultString := string(fSimpl)

	for _, d3 := range d3List {
		for i, device := range d3.GetListOfLights() {
			if device.GetName() != "" {
				for k, v := range signalLightMap {
					d3mapping := d3map.NewD3Map(
						c.Prefix.RoomPrefix,
						d3.GetRoomNumber(),
						c.CoreLight.SignalName,
						i+1,
						v,
						d3.GetRoomName(),
						device.GetName(),
						k)

					resultString = Replace(resultString, d3mapping)
				}
				if device.GetType() == "D" {
					for k, v := range signalLightDimmerMap {
						d3mapping := d3map.NewD3Map(
							c.Prefix.RoomPrefix,
							d3.GetRoomNumber(),
							c.CoreLight.SignalName,
							i+1,
							v,
							d3.GetRoomName(),
							device.GetName(),
							k)

						resultString = Replace(resultString, d3mapping)
					}
				}
			}
		}
	}

	resultByte := []byte(resultString)

	if err := os.WriteFile(simplPathNew, resultByte, 0666); err != nil {
		log.Fatalf("error writting file: %v", err)
	}

	log.Println("File: " + simplPathNew + " is created")
}

func Replace(resultString string, d3mapping *d3map.D3map) string {

	var isSuccess bool
	var err error
	resultString, isSuccess, err = d3mapping.Replace(resultString)
	if err != nil {
		log.Fatalf("Find and replace error: %v", err)
	}
	if isSuccess {
		log.Printf("Signal SUCCESSFUL found and replaced: %v", d3mapping)
	} else {
		log.Printf("Signal DID NOT FIND: %v", d3mapping)
	}

	return resultString
}
