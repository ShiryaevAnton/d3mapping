package d3map

import (
	"fmt"
	"strconv"
	"strings"
)

type D3List struct {
	roomName     string
	roomNumber   string
	listOfLights []string
	listOfShades []string
}

func NewD3List(roomName string, roomNumber string, listOfLights string, listOfShades string) *D3List {

	var d3List D3List

	d3List.roomName = strings.ReplaceAll(roomName, " ", "_")

	temp := strings.ReplaceAll(listOfLights, " ", "_")
	d3List.listOfLights = strings.Split(temp, "|")

	temp = strings.ReplaceAll(listOfShades, " ", "_")
	d3List.listOfShades = strings.Split(temp, "|")
	roomNumberInit, _ := strconv.Atoi(roomNumber)

	if roomNumberInit < 10 {
		d3List.roomNumber = "0" + roomNumber
	} else {
		d3List.roomNumber = roomNumber
	}

	return &d3List
}

func (d *D3List) String() string {

	var listOfLights, listOfShades string

	for _, v := range d.listOfLights {
		listOfLights += v + " "
	}

	for _, v := range d.listOfShades {
		listOfShades += v + " "
	}

	return fmt.Sprintln("-----> Room" + d.roomNumber + ":" + d.roomName + " Lights:" + listOfLights + "Shades:" + listOfShades + "<-----")

}
