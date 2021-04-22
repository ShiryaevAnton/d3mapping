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

const (
	search = "s"
	pull   = "p"
	all    = "all"
)

var simplPath string
var configPath string
var d3Path string
var mode string

var c config.Config

func init() {

	fConfig, err := os.ReadFile("config.toml")
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	toml.Unmarshal(fConfig, &c)

	flag.StringVar(&configPath, "cp", "", "Path to config file")
	flag.StringVar(&simplPath, "sp", "", "Path to simpl file")
	flag.StringVar(&d3Path, "dp", "", "Path to d3 simpl program")
	flag.StringVar(&mode, "m", "all", "Mode")

	flag.Parse()
}

func main() {

	logName := "log" + strconv.FormatInt(time.Now().UTC().Unix(), 10) + ".txt"

	_, err := os.Stat("log")
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
		Searching()
	}
	if mode == pull {
		GetNames()
	}
	if mode == all {
		GetNames()
		Searching()
	}

}

func Searching() {
	f, err := excelize.OpenFile(configPath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	fSimpl, err := os.ReadFile(simplPath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	simplPathNew := strings.ReplaceAll(simplPath, ".smw", "") + "_UPDATE.smw"

	var rooms []d3map.Room

	//log.Println("\nConfig file: " + strings.ReplaceAll(configPath, "./", "") + " contains:\n")

	for i := c.SheetConfig.StartRow; true; i++ {

		room := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnRoom+strconv.Itoa(i))
		if room == "" {
			break
		}
		light := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnLight+strconv.Itoa(i))
		lightType := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnLightType+strconv.Itoa(i))
		roomLight := d3map.NewRoom(room, i-c.SheetConfig.StartRow+1, light, lightType, "light")
		rooms = append(rooms, roomLight)

		shade := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnShade+strconv.Itoa(i))
		shadeType := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnShadeType+strconv.Itoa(i))
		roomShade := d3map.NewRoom(room, i-c.SheetConfig.StartRow+1, shade, shadeType, "shade")
		rooms = append(rooms, roomShade)

	}

	//log.Println(rooms)

	//log.Println("Start searching signals")

	resultString := string(fSimpl)
	compliteMap := make(map[string]bool)
	for _, room := range rooms {
		for i, device := range room.GetDevices() {
			if device.GetName() != "" {
				for _, signal := range c.Signals {

					prefix := c.CoreSignal.Prefix
					if signal.OverrideRoomPrefix != "" {
						prefix = signal.OverrideRoomPrefix
					}

					signalName := c.CoreSignal.Name
					if signal.OverrideCoreName != "" {
						signalName = signal.OverrideCoreName
					}

					roomName := room.GetName()
					if signal.OverridePanelName != "" {
						roomName = signal.OverridePanelName
					}

					deviceName := device.GetName()
					if signal.OverrideDeviceName != "" {
						deviceName = signal.OverrideDeviceName
					}

					signalNumber := i + 1
					if signal.RoomLevelSignal < 0 {
						signalNumber = signal.RoomLevelSignal
					}

					if signal.SystemType != device.GetSystemType() {
						continue
					}

					if signal.DeviceType != "C" {
						if signal.DeviceType != device.GetDeiveType() {
							continue
						}
					}

					signalMap := d3map.NewSignalMap(
						prefix,
						room.GetNumber(),
						signalName+signal.CoreSignalModif,
						signalNumber,
						signal.CoreSuffix,
						roomName,
						deviceName+signal.PanelSignalModif,
						signal.PanelSuffix)

					if compliteMap[signalMap.String()] {
						continue
					} else {
						compliteMap[signalMap.String()] = true
					}

					resultString = Replace(resultString, signalMap)
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

func GetNames() {

	f, err := excelize.OpenFile(configPath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	fD3Simpl, err := os.ReadFile(d3Path)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	d3string := string(fD3Simpl)

	rooms, err := d3map.GetRooms(d3string)
	if err != nil {
		log.Fatalf("%v", err)
	}

	i := c.SheetConfig.StartRow

	for _, room := range rooms {
		name := strings.ReplaceAll(room.GetName(), "_", " ")
		f.SetCellStr(c.SheetConfig.SheetName, c.SheetConfig.ColomnRoom+strconv.Itoa(i), name)
		f.SetCellStr(c.SheetConfig.SheetName, c.SheetConfig.ColomnLight+strconv.Itoa(i), room.GetLightString())
		f.SetCellStr(c.SheetConfig.SheetName, c.SheetConfig.ColomnShade+strconv.Itoa(i), room.GetShadeString())
		i++
	}

	err = f.Save()
	if err != nil {
		log.Fatalf("error saving file: %v", err)
	}
}

func Replace(resultString string, signalMap *d3map.SignalMap) string {

	var isSuccess bool
	var err error
	resultString, isSuccess, err = signalMap.Replace(resultString)
	if err != nil {
		log.Fatalf("Searching error: %v", err)
	}
	if isSuccess {
		log.Printf("SUCCESS: %v", signalMap)
	} else {
		log.Printf("FAIL: %v", signalMap)
	}

	return resultString
}
