package room

import (
	"strconv"
	"strings"

	"github.com/ShiryaevAnton/d3mapping/utility"
)

const (
	regRoomName   = "Cmn1=Generate Room Name Signals.\\n"
	regDeviceName = "Cmn1=Generate Device Name Signals.\\n"
	equipment     = "Equipment"
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

	roomNameString, err := utility.GetNameString(regRoomName, simplString)
	if err != nil {
		return nil, err
	}

	i := 1

	for {
		pName, err := utility.GetName(strconv.Itoa(i), roomNameString)
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

	deviceNameString, err := utility.GetNameString(regDeviceName, simplString)
	if err != nil {
		return nil, err
	}

	signals, err := utility.GetSignals(roomName, simplString)
	if err != nil {
		return nil, err
	}

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
