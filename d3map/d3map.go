package d3map

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	regSignalPrefix = "H=([0-9]+).\\nNm="
	regSignalNumber = "([0-9]+)"
	regIONumber     = "[IO]([0-9]+)="
	regIO           = "[IO]([0-9]+)"
)

type D3map struct {
	keyPanel    string
	keyCore     string
	numberPanel string
	numberCore  string
}

func NewD3Map(prefixCore string, roomNumber int, signalNameCore string,
	signalNumberCore int, suffixCore string, roomNamePanel string, deviceNamePanel string, suffixPanel string) *D3map {

	var d3Map D3map

	if roomNumber < 10 {
		d3Map.keyCore = prefixCore + "_" + "0" + strconv.Itoa(roomNumber) + "_" + signalNameCore
	} else {
		d3Map.keyCore = prefixCore + "_" + strconv.Itoa(roomNumber) + "_" + signalNameCore
	}

	if signalNumberCore < 10 {
		d3Map.keyCore += "_" + "0" + strconv.Itoa(signalNumberCore) + "_" + suffixCore
	} else {
		d3Map.keyCore += "_" + strconv.Itoa(signalNumberCore) + "_" + suffixCore
	}

	d3Map.keyPanel = "\\[" + roomNamePanel + "\\]" + "\\[" + deviceNamePanel + "\\]" + suffixPanel

	d3Map.numberCore = "0"
	d3Map.numberPanel = "0"

	return &d3Map
}

func (d *D3map) String() string {

	return fmt.Sprintln(d.keyCore + " --------> " + strings.ReplaceAll(d.keyPanel, "\\", "") + " " + d.numberCore + " --------> " + d.numberPanel)
}

func (d *D3map) Replace(simplString string) string {

	r, _ := regexp.Compile(regSignalPrefix + d.keyPanel + ".\\n")
	matchString := r.FindString(simplString)
	r, _ = regexp.Compile(regSignalNumber)
	d.numberPanel = r.FindString(matchString)

	r, _ = regexp.Compile(regSignalPrefix + d.keyCore + ".\\n")
	matchString = r.FindString(simplString)
	r, _ = regexp.Compile(regSignalNumber)
	d.numberCore = r.FindString(matchString)

	r, _ = regexp.Compile(regIONumber + d.numberCore + ".\\n")

	listOfMatchIO := r.FindAllString(simplString, -1)

	for _, matchIO := range listOfMatchIO {

		r, _ = regexp.Compile(regIO)
		IOName := r.FindString(matchIO)
		r, _ = regexp.Compile(matchIO)
		simplString = r.ReplaceAllString(simplString, IOName+"="+d.numberPanel+"\n")
	}

	return simplString
}
