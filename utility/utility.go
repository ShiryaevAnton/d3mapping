package utility

import (
	"regexp"
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
	allRoomLigths        = "All_Room_Lights"
	dimmer               = "dimmer"
)

func GetSignals(roomName string, simplString string) ([]string, error) {

	r, err := regexp.Compile("\\[" + roomName + "\\]" + "\\[" + "(.+)" + "\\]" + "Name")
	if err != nil {
		return nil, err
	}

	rawSignals := r.FindAllString(simplString, -1)

	var signals []string

	for _, signal := range rawSignals {

		if strings.Contains(signal, allRoomLigths) || strings.Contains(signal, dimmer) {
			continue
		}
		signals = append(signals, signal)
	}

	return signals, nil
}

func GetListOfIO(prefix string, root string, simplString string) ([]string, error) {

	r, err := regexp.Compile(prefix + root + ".\\n")
	if err != nil {
		return nil, err
	}
	return r.FindAllString(simplString, -1), nil
}

func GetIOName(matchIO string) (string, error) {

	r, err := regexp.Compile(regIO)
	if err != nil {
		return "", err
	}
	return r.FindString(matchIO), nil
}

func GetNumber(root string, simplString string) (string, error) {

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

func GetName(IONumber string, simplString string) (string, error) {

	r, err := regexp.Compile("P" + IONumber + "=" + regAnyLetter + regAnySymbolAnyTime + "\\n")
	if err != nil {
		return "", err
	}

	return r.FindString(simplString), nil
}

func GetNameString(reg string, simplString string) (string, error) {

	var nameString string
	modif := ""

	for {
		r, err := regexp.Compile(reg + modif)
		if err != nil {
			return "", err
		}
		nameString = r.FindString(simplString)
		if strings.Contains(nameString, "]") {
			break
		}
		modif = modif + regAnySymbolAnyTimeN
	}

	return nameString, nil
}
