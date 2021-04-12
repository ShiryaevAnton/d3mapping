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

	return fmt.Sprintln(d.keyCore + " --------> " + strings.ReplaceAll(d.keyPanel, "\\", ""))
}

func (d *D3map) Replace(simplString string) (string, bool, error) {

	r, err := regexp.Compile(regSignalPrefix + d.keyPanel + ".\\n")
	if err != nil {
		return "", false, err
	}
	matchString := r.FindString(simplString)

	r, err = regexp.Compile(regSignalNumber)
	if err != nil {
		return "", false, err
	}
	d.numberPanel = r.FindString(matchString)
	if d.numberPanel == "" {
		return simplString, false, nil
	}

	r, err = regexp.Compile(regSignalPrefix + d.keyCore + ".\\n")
	if err != nil {
		return "", false, err
	}
	matchString = r.FindString(simplString)
	r, err = regexp.Compile(regSignalNumber)
	if err != nil {
		return "", false, err
	}
	d.numberCore = r.FindString(matchString)
	if d.numberCore == "" {
		return simplString, false, nil
	}

	r, err = regexp.Compile(regIONumber + d.numberCore + ".\\n")
	if err != nil {
		return "", false, err
	}

	listOfMatchIO := r.FindAllString(simplString, -1)

	for _, matchIO := range listOfMatchIO {

		r, err = regexp.Compile(regIO)
		if err != nil {
			return "", false, err
		}
		IOName := r.FindString(matchIO)
		r, err = regexp.Compile(matchIO)
		if err != nil {
			return "", false, err
		}
		simplString = r.ReplaceAllString(simplString, IOName+"="+d.numberPanel+"\n")
	}

	return simplString, true, nil
}
