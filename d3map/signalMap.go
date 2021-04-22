package d3map

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ShiryaevAnton/d3mapping/regRepo"
)

const (
	regIONumber = "[IO]([0-9]+)="

	regTitle = "PrNm="
)

type SignalMap struct {
	keyPanel string
	keyCore  string
}

func NewSignalMap(prefixCore string, roomNumber int, signalNameCore string,
	signalNumberCore int, suffixCore string, roomNamePanel string, deviceNamePanel string, suffixPanel string) *SignalMap {

	var signalMap SignalMap

	if roomNumber < 10 {
		signalMap.keyCore = prefixCore + "_" + "0" + strconv.Itoa(roomNumber) + "_" + signalNameCore
	} else {
		signalMap.keyCore = prefixCore + "_" + strconv.Itoa(roomNumber) + "_" + signalNameCore
	}

	if signalNumberCore < 0 {
		signalMap.keyCore += "_" + suffixCore
	} else if signalNumberCore > 0 && signalNumberCore < 10 {
		signalMap.keyCore += "_" + "0" + strconv.Itoa(signalNumberCore) + "_" + suffixCore
	} else {
		signalMap.keyCore += "_" + strconv.Itoa(signalNumberCore) + "_" + suffixCore
	}

	signalMap.keyPanel = "\\[" + roomNamePanel + "\\]" + "\\[" + deviceNamePanel + "\\]" + suffixPanel

	return &signalMap
}

func (d *SignalMap) String() string {

	return fmt.Sprintln(d.keyCore + " --------> " + strings.ReplaceAll(d.keyPanel, "\\", ""))
}

func (d *SignalMap) Replace(simplString string) (string, bool, error) {

	numberCore, err := regRepo.GetNumber(d.keyCore, simplString)
	if err != nil {
		return "", false, err
	}
	if numberCore == "" {
		return simplString, false, nil
	}

	numberPanel, err := regRepo.GetNumber(d.keyPanel, simplString)
	if err != nil {
		return "", false, err
	}
	if numberPanel == "" {
		return simplString, false, nil
	}

	listOfMatchIO, err := regRepo.GetListOfIO(regIONumber, numberCore, simplString)
	if err != nil {
		return "", false, err
	}

	for _, matchIO := range listOfMatchIO {

		IOName, err := regRepo.GetIOName(matchIO)
		if err != nil {
			return "", false, err
		}
		simplString, err = regRepo.Replace(matchIO, IOName, numberPanel, simplString)
		if err != nil {
			return "", false, err
		}
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
