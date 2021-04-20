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
	"github.com/ShiryaevAnton/d3mapping/room"
	"github.com/pelletier/go-toml"
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

	fD3Simpl, err := os.ReadFile(d3Path)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

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
		d3ListLight := d3map.NewD3List(room, i-c.SheetConfig.StartRow+1, light, lightType, "light")
		d3List = append(d3List, d3ListLight)

		shade := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnShade+strconv.Itoa(i))
		shadeType := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnShadeType+strconv.Itoa(i))
		d3ListShade := d3map.NewD3List(room, i-c.SheetConfig.StartRow+1, shade, shadeType, "shade")
		d3List = append(d3List, d3ListShade)

	}

	//log.Println(d3List)

	//log.Println("Start searching signals")

	resultString := string(fSimpl)
	d3string := string(fD3Simpl)
	compliteMap := make(map[string]bool)

	for _, d3 := range d3List {
		for i, device := range d3.GetDevices() {
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

					roomName := d3.GetRoomName()
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

					d3mapping := d3map.NewD3Map(
						prefix,
						d3.GetRoomNumber(),
						signalName+signal.CoreSignalModif,
						signalNumber,
						signal.CoreSuffix,
						roomName,
						deviceName+signal.PanelSignalModif,
						signal.PanelSuffix)

					if compliteMap[d3mapping.String()] {
						continue
					} else {
						compliteMap[d3mapping.String()] = true
					}

					resultString = Replace(resultString, d3mapping)
				}
			}
		}
	}

	// roomNames, err := room.GetRoomName(d3string)
	// if err != nil {
	// 	log.Fatalf("error pulling room names: %v", err)
	// }

	// for _, roomName := range roomNames {
	// 	room.GetRooms(roomName, d3string)
	// }

	room.GetRooms("Bedroom_1", d3string)

	resultByte := []byte(resultString)

	if err := os.WriteFile(simplPathNew, resultByte, 0666); err != nil {
		log.Fatalf("error writting file: %v", err)
	}

	//log.Println("File: " + simplPathNew + " is created")
}

func Replace(resultString string, d3mapping *d3map.D3map) string {

	var isSuccess bool
	var err error
	resultString, isSuccess, err = d3mapping.Replace(resultString)
	if err != nil {
		log.Fatalf("Searching error: %v", err)
	}
	if isSuccess {
		log.Printf("SUCCESS: %v", d3mapping)
	} else {
		log.Printf("FAIL: %v", d3mapping)
	}

	return resultString
}
