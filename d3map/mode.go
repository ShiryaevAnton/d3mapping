package d3map

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ShiryaevAnton/d3mapping/config"
)

func Comparing(configPath string, c config.Config) error {

	f, err := excelize.OpenFile(configPath)
	if err != nil {
		return err
	}

	fTemp, err := excelize.OpenFile("temp/temp.xlsx")
	if err != nil {
		return err
	}

	for i := 0; true; i++ {

		index := fTemp.GetCellValue("Sheet1", "B"+strconv.Itoa(i+1))
		if index == "" {
			break
		}

		for j := 1; true; j++ {
			tempIndex := fTemp.GetCellValue("Sheet1", "D"+strconv.Itoa(j))

			if tempIndex == "" {
				break
			}

			if tempIndex == index {
				name := fTemp.GetCellValue("Sheet1", "C"+strconv.Itoa(j))
				light := fTemp.GetCellValue("Sheet1", "E"+strconv.Itoa(j))
				shade := fTemp.GetCellValue("Sheet1", "F"+strconv.Itoa(j))
				f.SetCellStr(c.SheetConfig.SheetName, c.SheetConfig.ColomnRoom+strconv.Itoa(i+c.SheetConfig.StartRow), name)
				f.SetCellStr(c.SheetConfig.SheetName, c.SheetConfig.ColomnLight+strconv.Itoa(i+c.SheetConfig.StartRow), light)
				f.SetCellStr(c.SheetConfig.SheetName, c.SheetConfig.ColomnShade+strconv.Itoa(i+c.SheetConfig.StartRow), shade)

			}

		}
	}

	err = f.Save()
	if err != nil {
		return err
	}

	return nil
}

func Searching(configPath string, simplPath string, c config.Config) error {
	f, err := excelize.OpenFile(configPath)
	if err != nil {
		return err
	}

	fSimpl, err := os.ReadFile(simplPath)
	if err != nil {
		return err
	}

	simplPathNew := strings.ReplaceAll(simplPath, ".smw", "") + "_UPDATE.smw"

	var rooms []Room

	for i := c.SheetConfig.StartRow; true; i++ {

		room := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnRoom+strconv.Itoa(i))
		if room == "" {
			break
		}
		light := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnLight+strconv.Itoa(i))
		lightType := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnLightType+strconv.Itoa(i))
		roomLight := NewRoom(room, i-c.SheetConfig.StartRow+1, light, lightType, "light")
		rooms = append(rooms, roomLight)

		shade := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnShade+strconv.Itoa(i))
		shadeType := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnShadeType+strconv.Itoa(i))
		roomShade := NewRoom(room, i-c.SheetConfig.StartRow+1, shade, shadeType, "shade")
		rooms = append(rooms, roomShade)

	}

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

					signalMap := NewSignalMap(
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

					resultString = replaceSignal(resultString, signalMap)
				}
			}
		}
	}

	resultByte := []byte(resultString)

	if err := os.WriteFile(simplPathNew, resultByte, 0666); err != nil {
		return err
	}

	log.Println("File: " + simplPathNew + " is created")

	return nil

}

func GetNames(configPath string, d3Path string, c config.Config) error {

	f, err := excelize.OpenFile(configPath)
	if err != nil {
		return err
	}

	fD3Simpl, err := os.ReadFile(d3Path)
	if err != nil {
		return err
	}

	d3string := string(fD3Simpl)

	_, err = os.Stat("temp")
	if os.IsNotExist(err) {
		err := os.Mkdir("temp\\", 0666)
		if err != nil {
			return err
		}
	}

	fTemp := excelize.NewFile()

	rooms, err := GetRooms(d3string)
	if err != nil {
		return err
	}

	for i := 0; true; i++ {

		room := f.GetCellValue(c.SheetConfig.SheetName, c.SheetConfig.ColomnRoom+strconv.Itoa(i+c.SheetConfig.StartRow))
		if room == "" {
			break
		}
		fTemp.SetCellStr("Sheet1", "A"+strconv.Itoa(i+1), room)
		fTemp.SetCellInt("Sheet1", "B"+strconv.Itoa(i+1), i+1)

	}

	for i, room := range rooms {
		name := strings.ReplaceAll(room.GetName(), "_", " ")
		fTemp.SetCellStr("Sheet1", "C"+strconv.Itoa(i+1), name)
		fTemp.SetCellStr("Sheet1", "E"+strconv.Itoa(i+1), room.GetLightString())
		fTemp.SetCellStr("Sheet1", "F"+strconv.Itoa(i+1), room.GetShadeString())
	}

	if err := fTemp.SaveAs("temp/temp.xlsx"); err != nil {
		return err
	}

	return nil
}

func replaceSignal(resultString string, signalMap *SignalMap) string {

	var isSuccess bool
	var err error
	resultString, isSuccess, err = signalMap.Replace(resultString)
	if err != nil {
		log.Fatalf("Searching error: %v", err)
	}
	if isSuccess {
		log.Printf("SUCCEDED - %v", signalMap)
	} else {
		log.Printf("FAILED - %v", signalMap)
	}

	return resultString
}
