package room

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	regAnyLetter         = "[A-Za-z]"
	regAnySymbolAnyTime  = "(.+)"
	regAnySymbolAnyTimeN = "(.+)\\n"
	regRoomName          = "Cmn1=Generate Room Name Signals.\\n"
	regDeviceName        = "Cmn1=Generate Device Name Signals.\\n"
	allRoomLigths        = "All_Room_Lights"
	dimmer               = "dimmer"
	equipment            = "Equipment"
)

type device struct {
	name        string
	deviceTypes string
}

type Room struct {
	name   string
	number string
	lights device
	shades device
}

func GetRoomName(simplString string) ([]string, error) {

	var roomNames []string

	roomNameString, err := getNameString(regRoomName, simplString)
	if err != nil {
		return nil, err
	}

	i := 1

	for {
		pName, err := getName(strconv.Itoa(i), roomNameString)
		if err != nil {
			return nil, err
		}
		if pName == "" {
			break
		}
		roomName := strings.TrimLeft(pName, "P"+strconv.Itoa(i)+"=")
		roomName = strings.TrimRight(roomName, "\n")

		i++
		if strings.Contains(roomName, equipment) {
			continue
		}

		roomNames = append(roomNames, roomName)
	}

	return roomNames, nil
}

func GetRooms(roomName string, simplString string) (*Room, error) {

	var room *Room

	// deviceNameString, err := getNameString(regDeviceName, simplString)
	// if err != nil {
	// 	return nil, err
	// }

	//[Bedroom_1][Shade_1]Name

	//r, _ := regexp.Compile("\\[Bedroom_1\\]\\[Shade_1\\]Name")
	r, err := regexp.Compile("\\[" + roomName + "\\]" + "\\[" + "(.+)" + "\\]" + "Name")
	if err != nil {
		return nil, err
	}

	rawDevices := r.FindAllString(simplString, -1)

	var devices []string

	for _, device := range rawDevices {

		if strings.Contains(device, allRoomLigths) || strings.Contains(device, dimmer) {
			continue
		}
		devices = append(devices, device)
	}

	fmt.Println(devices)

	return room, nil
}

func (r *Room) GetName() string {
	return r.name
}

func (r *Room) GetNumber() string {
	return r.number
}

func (r *Room) GetLightType() string {
	return r.lights.deviceTypes
}

func (r *Room) GetShadeType() string {
	return r.shades.deviceTypes
}

func (r *Room) GetLightName() string {
	return r.lights.name
}

func (r *Room) GetShadeName() string {
	return r.shades.name
}

func getName(IONumber string, simplString string) (string, error) {

	r, err := regexp.Compile("P" + IONumber + "=" + regAnyLetter + regAnySymbolAnyTime + "\\n")
	if err != nil {
		return "", err
	}

	return r.FindString(simplString), nil
}

func getNameString(reg string, simplString string) (string, error) {

	var nameString string
	modif := ""

	for {
		r, err := regexp.Compile(reg + modif)
		if err != nil {
			return "", err
		}
		nameString = r.FindString(simplString)
		if strings.Contains(nameString, "]") {
			break
		}
		modif = modif + regAnySymbolAnyTimeN
	}

	return nameString, nil
}
