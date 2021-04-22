package regRepo

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
	switchs              = "switch"
	lightSuffix          = "On_fb"
	shadeSuffix          = "Upper_Limit_Set"
	regRoomName          = "Cmn1=Generate Room Name Signals.\\n"
	regDeviceName        = "Cmn1=Generate Device Name Signals.\\n"
)

func GetSignals(roomName string, simplString string) ([]string, error) {

	r, err := regexp.Compile("\\[" + roomName + "\\]" + "\\[" + "(.+)" + "\\]" + "Name")
	if err != nil {
		return nil, err
	}

	rawSignals := r.FindAllString(simplString, -1)

	var signals []string

	for _, signal := range rawSignals {

		if strings.Contains(signal, allRoomLigths) || strings.Contains(signal, dimmer) || strings.Contains(signal, switchs) {
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

func GetRoomNameString(simplString string) (string, error) {
	return getNameString(regRoomName, simplString)
}

func GetDeviceNameString(simplString string) (string, error) {
	return getNameString(regDeviceName, simplString)
}

func IsLight(signal string, simplString string) (bool, error) {

	light := strings.ReplaceAll(signal, "Name", lightSuffix)
	light = strings.ReplaceAll(light, "]", "\\]")
	light = strings.ReplaceAll(light, "[", "\\[")

	return regexp.MatchString(light, simplString)

}

func IsShade(signal string, simplString string) (bool, error) {

	shade := strings.ReplaceAll(signal, "Name", shadeSuffix)
	shade = strings.ReplaceAll(shade, "]", "\\]")
	shade = strings.ReplaceAll(shade, "[", "\\[")

	return regexp.MatchString(shade, simplString)
}

func Replace(matchIO string, IOName string, numberPanel string, simplString string) (string, error) {
	r, err := regexp.Compile(matchIO)
	if err != nil {
		return "", err
	}
	return r.ReplaceAllString(simplString, IOName+"="+numberPanel+"\r\n"), nil
}

func getNameString(reg string, simplString string) (string, error) {

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
