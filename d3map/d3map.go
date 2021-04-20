package d3map

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	regSignalPrefix      = "H=([0-9]+).\\nNm="
	regAnyNumberAnyTime  = "([0-9]+)"
	regIONumber          = "[IO]([0-9]+)="
	regIO                = "[IO]([0-9]+)"
	regONumber           = "O([0-9]+)="
	regTitle             = "PrNm="
	regAnyLetter         = "[A-Za-z]"
	regAnySymbolAnyTime  = "(.+)"
	regAnySymbolAnyTimeN = "(.+)\\n"
)

type D3map struct {
	keyPanel string
	keyCore  string
}

func NewD3Map(prefixCore string, roomNumber int, signalNameCore string,
	signalNumberCore int, suffixCore string, roomNamePanel string, deviceNamePanel string, suffixPanel string) *D3map {

	var d3Map D3map

	if roomNumber < 10 {
		d3Map.keyCore = prefixCore + "_" + "0" + strconv.Itoa(roomNumber) + "_" + signalNameCore
	} else {
		d3Map.keyCore = prefixCore + "_" + strconv.Itoa(roomNumber) + "_" + signalNameCore
	}

	if signalNumberCore < 0 {
		d3Map.keyCore += "_" + suffixCore
	} else if signalNumberCore > 0 && signalNumberCore < 10 {
		d3Map.keyCore += "_" + "0" + strconv.Itoa(signalNumberCore) + "_" + suffixCore
	} else {
		d3Map.keyCore += "_" + strconv.Itoa(signalNumberCore) + "_" + suffixCore
	}

	d3Map.keyPanel = "\\[" + roomNamePanel + "\\]" + "\\[" + deviceNamePanel + "\\]" + suffixPanel

	return &d3Map
}

func (d *D3map) String() string {

	return fmt.Sprintln(d.keyCore + " --------> " + strings.ReplaceAll(d.keyPanel, "\\", ""))
}

func (d *D3map) Replace(simplString string) (string, bool, error) {

	numberCore, err := getNumber(d.keyCore, simplString)
	if err != nil {
		return "", false, err
	}
	if numberCore == "" {
		return simplString, false, nil
	}

	numberPanel, err := getNumber(d.keyPanel, simplString)
	if err != nil {
		return "", false, err
	}
	if numberPanel == "" {
		return simplString, false, nil
	}

	listOfMatchIO, err := getListOfIO(regIONumber, numberCore, simplString)
	if err != nil {
		return "", false, err
	}

	for _, matchIO := range listOfMatchIO {

		IOName, err := getIOName(matchIO)
		if err != nil {
			return "", false, err
		}
		r, err := regexp.Compile(matchIO)
		if err != nil {
			return "", false, err
		}
		simplString = r.ReplaceAllString(simplString, IOName+"="+numberPanel+"\n")
	}

	return simplString, true, nil
}

func ReplaceTitle(simplString string, originalTitle string, replaceTitle string) (string, error) {

	r, err := regexp.Compile(regTitle + originalTitle + ".\\n")
	if err != nil {
		return "", err
	}

	simplString = r.ReplaceAllString(simplString, replaceTitle)

	return simplString, nil
}

func getListOfIO(prefix string, root string, simplString string) ([]string, error) {

	r, err := regexp.Compile(prefix + root + ".\\n")
	if err != nil {
		return nil, err
	}
	return r.FindAllString(simplString, -1), nil
}

func getIOName(matchIO string) (string, error) {

	r, err := regexp.Compile(regIO)
	if err != nil {
		return "", err
	}
	return r.FindString(matchIO), nil
}

func getNumber(root string, simplString string) (string, error) {

	r, err := regexp.Compile(regSignalPrefix + root + ".\\n")
	if err != nil {
		return "", err
	}

	matchString := r.FindString(simplString)

	r, err = regexp.Compile(regAnyNumberAnyTime)
	if err != nil {
		return "", err
	}

	return r.FindString(matchString), nil
}
