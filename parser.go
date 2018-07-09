package sdptransform

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/Jeffail/gabs"
)

const maxLineSize = 1024

var keyvalueRegex = regexp.MustCompile("^\\s*([^= ]+)(?:\\s*=\\s*([^ ]+))?$")
var validLineRegex = regexp.MustCompile("^([a-z])=(.*)")

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func Parse(sdp []byte) (session *gabs.Container, err error) {

	bufferReader := bufio.NewReader(bytes.NewReader(sdp))
	session = gabs.New()
	location := session

	session.Array("media")

	for {
		if line, _, err := bufferReader.ReadLine(); err != nil {

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
