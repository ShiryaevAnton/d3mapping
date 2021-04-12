package d3map

import (
	"fmt"
	"strconv"
	"strings"
)

type Light struct {
	name  string
	types string
}

type Shade struct {
	name  string
	types string
}

type D3List struct {
	roomName   string
	roomNumber int
	lights     []Light
	shades     []Shade
}

func NewD3List(roomName string, roomNumber int, listOfLights string,
	listOfLightsType string, listOfShades string, listOfShadesType string) D3List {

	var d3List D3List

	d3List.roomName = strings.ReplaceAll(roomName, " ", "_")

	temp := strings.ReplaceAll(listOfLights, " ", "_")
	tempListOfLights := strings.Split(temp, "|")
	tempListOfLightsType := strings.Split(listOfLightsType, "")

	temp = strings.ReplaceAll(listOfShades, " ", "_")
	tempListOfShades := strings.Split(temp, "|")
	tempListOfShadesType := strings.Split(listOfShadesType, "")

	for i, light := range tempListOfLights {
		if i < len(tempListOfLightsType) {
			d3List.lights = append(d3List.lights, Light{name: light, types: tempListOfLightsType[i]})
		}
	}

	for i, shade := range tempListOfShades {
		if i < len(tempListOfShadesType) {
			d3List.shades = append(d3List.shades, Shade{name: shade, types: tempListOfShadesType[i]})
		}
	}

	d3List.roomNumber = roomNumber

	return d3List
}

func (d D3List) String() string {

	var listOfLights, listOfShades string

	for _, v := range d.lights {
		listOfLights += "{" + "Name:" + v.name + " " + "Type:" + v.types + "}, "
	}

	for _, v := range d.shades {
		listOfShades += "{" + "Name:" + v.name + " " + "Type:" + v.types + "}, "
	}

	return fmt.Sprintln("-----> Room" + strconv.Itoa(d.roomNumber) + ":" + d.roomName + " Lights:" + listOfLights + " " + "Shades:" + listOfShades + "<-----")

}

func (d D3List) GetRoomName() string {
	return d.roomName
}

func (d D3List) GetRoomNumber() int {
	return d.roomNumber
}

func (d D3List) GetListOfLights() []Light {
	return d.lights
}

func (d D3List) GetListOfShade() []Shade {
	return d.shades
}

func (l Light) GetName() string {
	return l.name
}

func (l Light) GetType() string {
	return l.types
}

func (l Shade) GetName() string {
	return l.name
}

func (l Shade) GetType() string {
	return l.types
}
