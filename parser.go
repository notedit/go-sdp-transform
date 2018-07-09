package sdptransform

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strconv"

	"github.com/Jeffail/gabs"
)

const maxLineSize = 1024

var keyvalueRegex = regexp.MustCompile("^\\s*([^= ]+)(?:\\s*=\\s*([^ ]+))?$")
var validLineRegex = regexp.MustCompile("^([a-z])=(.*)")

func Parse(sdp []byte) (session *gabs.Container, err error) {

	bufferReader := bufio.NewReader(bytes.NewReader(sdp))
	session = gabs.New()
	location := session

	session.Array("media")

	for {
		if line, err := bufferReader.ReadSlice('\n'); err != nil {

			if len(line) > 0 && line[len(line)-1] == '\r' {
				line = line[:len(line)-1]
			}
			if !validLineRegex.Match(line) {
				continue
			}

			lineType := line[0]
			content := line[2:]

			fmt.Println("line type ", lineType, "content ", content)

			if lineType == byte('m') {
				m := gabs.New()
				m.Array("rtp")
				m.Array("fmtp")

				session.ArrayAppend(m.Data(), "media")
				index, _ := session.ArrayCount("media")
				location, _ = session.ArrayElement(index, "media")

				fmt.Println("add media ", session.String())
			}

			if _, ok := rulesMap[lineType]; !ok {
				fmt.Println("sdp can not find type ", lineType)
				continue
			}

			rules := rulesMap[lineType]

			for _, rule := range rules {
				if rule.Reg.Match(content) {
					parseReg(rule, location, content)
				}
			}

		}
	}
	fmt.Println("parsed session", session.String())
	return
}

func ParseParams(str []byte) map[string]string {

	bufferReader := bufio.NewReader(bytes.NewReader(str))
	ret := map[string]string{}

	for {
		if param, err := bufferReader.ReadSlice(';'); err != nil {
			param = bytes.TrimSpace(param)
			if len(param) == 0 {
				continue
			}
			keyValue := bytes.SplitN(param, []byte{'='}, 2)

			if len(keyValue) == 2 {
				ret[string(keyValue[0])] = string(keyValue[1])
			} else if len(keyValue) == 1 {
				ret[string(keyValue[0])] = ""
			}
		}
	}

	return ret

}

func ParsePayloads(str []byte) []int {

	keyValue := bytes.Split(str, []byte{' '})
	payloads := []int{}

	for _, value := range keyValue {
		value = bytes.TrimSpace(value)
		if len(value) == 0 {
			continue
		}
		intValue, _ := strconv.Atoi(string(value))
		payloads = append(payloads, intValue)
	}

	return payloads
}

func ParseImageAttributes(str []byte) []map[string]int {

	bufferReader := bufio.NewReader(bytes.NewReader(str))
	ret := []map[string]int{}

	for {
		if param, err := bufferReader.ReadSlice(' '); err != nil {
			param = bytes.TrimSpace(param)
			if len(param) == 0 {
				continue
			}
			if len(param) < 5 {
				continue
			}

			keyValues := bytes.Split(param[1:len(param)-2], []byte{','})

			retMap := map[string]int{}

			for _, keyValue := range keyValues {
				_keyValue := bytes.SplitN(keyValue, []byte{'='}, 2)
				if len(_keyValue) != 2 {
					continue
				}

				value, err := strconv.Atoi(string(_keyValue[1]))
				if err != nil {
					continue
				}

				retMap[string(_keyValue[0])] = value

			}

			ret = append(ret, retMap)
		}
	}

	return ret

}

func ParseSimulcastStreamList(str []byte) [][]map[string]interface{} {

	ret := [][]map[string]interface{}{}

	streams := bytes.Split(str, []byte{';'})

	for _, stream := range streams {

		streamFormat := []map[string]interface{}{}

		formats := bytes.Split(stream, []byte{','})

		for _, format := range formats {
			var scid interface{}
			paused := false

			if format[0] != '~' {
				scid = toType(string(format), 'd')
			} else {
				scid = toType(string(format[1:len(format)-1]), 'd')
				paused = true
			}

			streamFormat = append(streamFormat, map[string]interface{}{
				"scid":   scid,
				"paused": paused,
			})

		}
		ret = append(ret, streamFormat)
	}

	return ret
}

func parseReg(rule *Rule, location *gabs.Container, content []byte) {

	needsBlank := len(rule.Name) != 0 && len(rule.Names) != 0

	if len(rule.Push) != 0 {
		if !location.Exists(rule.Push) {
			location.Array(rule.Push)
		}
	} else if needsBlank {
		if !location.Exists(rule.Name) {
			location.Set(nil, rule.Name)
		}
	}

	match := rule.Reg.FindAll(content, -1)

	object := gabs.New()
	var keyLocation *gabs.Container

	if len(rule.Push) != 0 {
		keyLocation = object
	} else {
		if needsBlank {
			keyLocation = location.Path(rule.Name)
		} else {
			keyLocation = location
		}
	}

	attachProperties(match, keyLocation, rule.Names, rule.Name, rule.Types)

	if len(rule.Push) != 0 {
		location.ArrayAppend(keyLocation.Data(), rule.Push)
		fmt.Println("array append ", location.String())
	}
}

func attachProperties(match [][]byte, location *gabs.Container, names []string, rawName string, types []rune) {

	if len(rawName) != 0 && len(names) == 0 {
		location.Set(toType(string(match[1]), types[0]), rawName)
	} else {
		for i := 0; i < len(names); i++ {
			if len(match) > i+1 && match[i+1] != nil {
				location.Set(toType(string(match[i+1]), types[i]), names[i])
			}
		}
	}
}

func toType(str string, otype rune) interface{} {
	switch otype {
	case 's':
		return str
	case 'd':
		if number, err := strconv.Atoi(str); err != nil {
			return 0
		} else {
			return number
		}
	case 'f':
		if flo, err := strconv.ParseFloat(str, 64); err != nil {
			return 0.0
		} else {
			return flo
		}
	}
	return nil
}
