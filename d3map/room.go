package d3map

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ShiryaevAnton/d3mapping/regRepo"
)

const (
	equipment = "Equipment"
)

type Device struct {
	name        string
	deviceTypes string
	systemType  string
}

type Room struct {
	name    string
	number  int
	devices []Device
}

func (d Room) GetName() string {
	return d.name
}

func (d Room) GetNumber() int {
	return d.number
}

func (d Room) GetDevices() []Device {
	return d.devices
}

func (d Device) GetName() string {
	return d.name
}

func (d Device) GetDeiveType() string {
	return d.deviceTypes
}

func (d Device) GetSystemType() string {
	return d.systemType
}

func NewRoom(roomName string, roomNumber int, listOfDevice string, listOfDevicesType string, systemType string) Room {

	var room Room

	room.name = strings.ReplaceAll(roomName, " ", "_")

	temp := strings.ReplaceAll(listOfDevice, " ", "_")
	tempListOfDevices := strings.Split(temp, "|")
	tempListOfDevicesType := strings.Split(listOfDevicesType, "")

	for i, device := range tempListOfDevices {
		if i < len(tempListOfDevicesType) {
			room.devices = append(room.devices, Device{name: device, deviceTypes: tempListOfDevicesType[i], systemType: systemType})
		}
	}

	room.number = roomNumber

	return room
}

func (d Room) String() string {

	var listOfDeivces string

	for _, v := range d.devices {
		listOfDeivces += "{" + v.systemType + ": " + v.name + " " + "type:" + v.deviceTypes + "}, "
	}

	return fmt.Sprintln("Room" + strconv.Itoa(d.number) + "-" + d.name + ": " + listOfDeivces)

}

func GetRooms(simplString string) ([]Room, error) {

	var rooms []Room

	roomNames, err := getRoomNames(simplString)
	if err != nil {
		return nil, err
	}

	for _, roomName := range roomNames {
		room, _ := getRoom(roomName, simplString)
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (d Room) GetLightString() string {

	return d.getDeviceString("light")
}

func (d Room) GetShadeString() string {

	return d.getDeviceString("shade")
}

func (d Room) getDeviceString(systemType string) string {

	var res string

	for _, device := range d.devices {
		if device.systemType == systemType {
			res += device.name + "|"
		}
	}

	return res
}

func getRoom(roomName string, simplString string) (Room, error) {

	var room Room

	room.name = strings.ReplaceAll(roomName, " ", "_")

	signals, err := regRepo.GetSignals(room.name, simplString)
	if err != nil {
		return room, err
	}

	for _, signal := range signals {

		device, err := getDevice(signal, room.name, simplString)
		if err != nil {
			return room, err
		}

		room.devices = append(room.devices, device)
	}

	return room, nil
}

func getRoomNames(simplString string) ([]string, error) {

	var roomNames []string

	roomNameString, err := regRepo.GetRoomNameString(simplString)
	if err != nil {
		return nil, err
	}

	i := 1

	for {
		pName, err := regRepo.GetName(strconv.Itoa(i), roomNameString)
		if err != nil {
			return nil, err
		}
		if pName == "" {
			break
		}
		roomName := strings.TrimLeft(pName, "P"+strconv.Itoa(i)+"=")
		roomName = strings.ReplaceAll(roomName, "\r\n", "")

		i++
		if strings.Contains(roomName, equipment) {
			continue
		}

		roomNames = append(roomNames, roomName)
	}

	return roomNames, nil
}

func getDevice(signal string, roomName string, simplString string) (Device, error) {

	var device Device

	name := strings.TrimLeft(signal, "["+roomName+"][")
	name = strings.TrimRight(name, "][Name")
	name = strings.ReplaceAll(name, "_", " ")
	device.name = name

	match, err := regRepo.IsLight(signal, simplString)
	if err != nil {
		return device, err
	}
	if match {
		device.systemType = "light"
	}
	match, err = regRepo.IsShade(signal, simplString)
	if err != nil {
		return device, err
	}
	if match {
		device.systemType = "shade"
	}

	return device, nil
}
