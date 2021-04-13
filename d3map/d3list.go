package d3map

import (
	"fmt"
	"strconv"
	"strings"
)

type Device struct {
	name        string
	deviceTypes string
	systemType  string
}

type D3List struct {
	roomName   string
	roomNumber int
	devices    []Device
}

func NewD3List(roomName string, roomNumber int, listOfDevice string, listOfDevicesType string, systemType string) D3List {

	var d3List D3List

	d3List.roomName = strings.ReplaceAll(roomName, " ", "_")

	temp := strings.ReplaceAll(listOfDevice, " ", "_")
	tempListOfDevices := strings.Split(temp, "|")
	tempListOfDevicesType := strings.Split(listOfDevicesType, "")

	for i, device := range tempListOfDevices {
		if i < len(tempListOfDevicesType) {
			d3List.devices = append(d3List.devices, Device{name: device, deviceTypes: tempListOfDevicesType[i], systemType: systemType})
		}
	}

	d3List.roomNumber = roomNumber

	return d3List
}

func (d D3List) String() string {

	var listOfDeivces string

	for _, v := range d.devices {
		listOfDeivces += "{" + v.systemType + ": " + v.name + " " + "type:" + v.deviceTypes + "}, "
	}

	return fmt.Sprintln("-----> Room" + strconv.Itoa(d.roomNumber) + "-" + d.roomName + ": " + listOfDeivces + "<-----")

}

func (d D3List) GetRoomName() string {
	return d.roomName
}

func (d D3List) GetRoomNumber() int {
	return d.roomNumber
}

func (d D3List) GetDevices() []Device {
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
