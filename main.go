package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

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
var colomnRoomNumber string
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
	flag.StringVar(&signalLightName, "snl", "LIGHTING_LOAD", "Signal name for lights: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightOn, "sol", "ON", "Suffix for light ON: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightOnFb, "sofl", "ON_FB", "Suffix for light ON_FB: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightOff, "sfl", "OFF", "Suffix for light OFF: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightOffFb, "sffl", "OFF_FB", "Suffix for light OFF_FB: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightLevel, "sll", "ON", "Suffix for light LEVEL: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightLevelFb, "slfl", "ON_FB", "Suffix for light LEVEL_FB: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightRaise, "srl", "RAISE", "Suffix for light RAISE: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")
	flag.StringVar(&suffixLightDim, "srfl", "DIM", "Suffix for light DIM: Prefix_RoomNumber_SignalName_SignalNumber_Suffix")

	flag.StringVar(&sheetName, "sn", "PROJECTCONFIG", "Name of sheet")
	flag.StringVar(&colomnRoom, "colr", "B", "Room's colomn")
	flag.StringVar(&colomnRoomNumber, "colrn", "A", "Room's number")
	flag.StringVar(&colomnLigth, "coll", "L", "Light's colomn")
	flag.StringVar(&colomnShade, "cols", "O", "Shade's colomn")
	flag.IntVar(&startRow, "sr", 58, "Start row")

	flag.Parse()
}

func main() {

	f, err := excelize.OpenFile(configPath)

	if err != nil {
		panic(err)
	}

	d3List := make([]d3map.D3List, 8)

	fmt.Println("\nConfig file: " + strings.ReplaceAll(configPath, "./", "") + " contains:\n")

	for i := startRow; true; i++ {

		room := f.GetCellValue(sheetName, colomnRoom+strconv.Itoa(i))
		if room == "" {
			break
		}
		light := f.GetCellValue(sheetName, colomnLigth+strconv.Itoa(i))
		shade := f.GetCellValue(sheetName, colomnShade+strconv.Itoa(i))
		roomNumber := f.GetCellValue(sheetName, colomnRoomNumber+strconv.Itoa(i))

		newd3List := d3map.NewD3List(room, roomNumber, light, shade)

		fmt.Println(newd3List)

		d3List = append(d3List, *newd3List)
	}

	fmt.Println("---------------------------------------------------")

	// test := d3map.NewD3Map(roomPrefix, 1, signalLightName, 1, suffixLightRaise, "Games", "Accents", "Raise")

	// file, err := os.ReadFile(simplPath)
	// if err != nil {
	// 	panic(err)
	// }

	// temp := []byte(test.Replace(string(file)))
	// err = os.WriteFile("./test1.smw", temp, 0666)

}
