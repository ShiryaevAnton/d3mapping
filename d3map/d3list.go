package d3map

import (
	"fmt"
	"strconv"
	"strings"
)

type D3List struct {
	roomName     string
	roomNumber   int
	listOfLights []string
	listOfShades []string
}

func NewD3List(roomName string, roomNumber int, listOfLights string, listOfShades string) D3List {

	var d3List D3List

	d3List.roomName = strings.ReplaceAll(roomName, " ", "_")

	temp := strings.ReplaceAll(listOfLights, " ", "_")
	d3List.listOfLights = strings.Split(temp, "|")

	temp = strings.ReplaceAll(listOfShades, " ", "_")
	d3List.listOfShades = strings.Split(temp, "|")

	d3List.roomNumber = roomNumber

	return d3List
}

func (d D3List) String() string {

	var listOfLights, listOfShades string

	for _, v := range d.listOfLights {
		listOfLights += v + " "
	}

	for _, v := range d.listOfShades {
		listOfShades += v + " "
	}

	return fmt.Sprintln("-----> Room" + strconv.Itoa(d.roomNumber) + ":" + d.roomName + " Lights:" + listOfLights + "Shades:" + listOfShades + "<-----")

}

func (d D3List) GetRoomName() string {
	return d.roomName
}

func (d D3List) GetRoomNumber() int {
	return d.roomNumber
}

func (d D3List) GetListOfLights() []string {
	return d.listOfLights
}

func (d D3List) GetListOfShade() []string {
	return d.listOfShades
}
