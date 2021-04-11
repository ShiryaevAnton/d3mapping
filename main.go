package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ShiryaevAnton/d3mapping/d3map"
)

const (
	panelSuffixLightOn      = "On"
	panelSuffixLightOnFb    = "On_fb"
	panelSuffixLightOff     = "Off"
	panelSuffixLightOffFb   = "Off_fb"
	panelSuffixLightRaise   = "Raise"
	panelSuffixLightDim     = "Lower"
	panelSuffixLightLevel   = "Level"
	panelSuffixLightLevelFb = "Current_Level"
)

var sheetName string
var colomnRoom string

//var colomnRoomNumber string
var colomnLigth string
var colomnShade string
var startRow int
var configPath string
var simplPath string
var roomPrefix string
var signalLightName string
var suffixLightOn string
var suffixLightOff string
var suffixLightOnFb string
var suffixLightOffFb string
var suffixLightLevel string
var suffixLightLevelFb string
var suffixLightRaise string
var suffixLightDim string

func init() {
	flag.StringVar(&configPath, "cp", "", "Path to config file")
	flag.StringVar(&simplPath, "sp", "", "Path to simpl file")

	flag.StringVar(&roomPrefix, "rp", "ROOM", "Prefix for signals: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&signalLightName, "snl", "PL_LIGHT", "Signal name for lights: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightOn, "sol", "ON", "Suffix for light ON: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightOnFb, "sofl", "ON_FB", "Suffix for light ON_FB: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightOff, "sfl", "OFF", "Suffix for light OFF: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightOffFb, "sffl", "OFF_FB", "Suffix for light OFF_FB: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightLevel, "sll", "LVL", "Suffix for light LEVEL: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightLevelFb, "slfl", "LVL_FB", "Suffix for light LEVEL_FB: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightRaise, "srl", "RAISE", "Suffix for light RAISE: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightDim, "srfl", "DIM", "Suffix for light DIM: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")

	flag.StringVar(&sheetName, "sn", "PROJECTCONFIG", "Name of sheet")
	flag.StringVar(&colomnRoom, "colr", "B", "Room's colomn")
	//flag.StringVar(&colomnRoomNumber, "colrn", "A", "Room's number")
	flag.StringVar(&colomnLigth, "coll", "L", "Light's colomn")
	flag.StringVar(&colomnShade, "cols", "O", "Shade's colomn")
	flag.IntVar(&startRow, "sr", 58, "Start row")

	flag.Parse()
}

func main() {

	logName := "log" + strconv.FormatInt(time.Now().UTC().Unix(), 10) + ".txt"

	fLog, err := os.OpenFile(logName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
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

	var signalSuffixMap = make(map[string]string)

	signalSuffixMap[panelSuffixLightOn] = suffixLightOn
	signalSuffixMap[panelSuffixLightOnFb] = suffixLightOnFb
	signalSuffixMap[panelSuffixLightOff] = suffixLightOff
	signalSuffixMap[panelSuffixLightOffFb] = suffixLightOffFb
	signalSuffixMap[panelSuffixLightRaise] = suffixLightRaise
	signalSuffixMap[panelSuffixLightDim] = suffixLightDim
	signalSuffixMap[panelSuffixLightLevel] = suffixLightLevel
	signalSuffixMap[panelSuffixLightLevelFb] = suffixLightLevelFb

	simplPathNew := strings.ReplaceAll(simplPath, ".smw", "") + "_UPDATE.smw"

	var d3List []d3map.D3List

	log.Println("\nConfig file: " + strings.ReplaceAll(configPath, "./", "") + " contains:\n")

	for i := startRow; true; i++ {

		room := f.GetCellValue(sheetName, colomnRoom+strconv.Itoa(i))
		if room == "" {
			break
		}
		light := f.GetCellValue(sheetName, colomnLigth+strconv.Itoa(i))
		shade := f.GetCellValue(sheetName, colomnShade+strconv.Itoa(i))

		newd3List := d3map.NewD3List(room, i-startRow+1, light, shade)

		fmt.Println(newd3List)

		d3List = append(d3List, newd3List)
	}

	log.Println("Start to find and replace signals")

	resultString := string(fSimpl)

	for _, d3 := range d3List {
		for i, device := range d3.GetListOfLights() {
			if device != "" {
				for k, v := range signalSuffixMap {
					d3mapping := d3map.NewD3Map(roomPrefix, d3.GetRoomNumber(), signalLightName, i+1,
						v, d3.GetRoomName(), device, k)

					resultString = d3mapping.Replace(resultString)
					log.Println(d3mapping)
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
