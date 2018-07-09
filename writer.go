package sdptransform

import (
	"fmt"
	"strings"

	"github.com/Jeffail/gabs"
)

var outerOrder = []byte{'v', 'o', 's', 'i', 'u', 'e', 'p', 'c', 'b', 't', 'r', 'z', 'a'}
var innerOrder = []byte{'i', 'c', 'b', 'a'}

func Write(session *gabs.Container) string {

	if !session.Exists("version") {
		session.Set(0, "version")
	}

	if !session.Exists("name") {
		session.Set("-", "name")
	}

	if !session.Exists("media") {
		session.Set([]interface{}{}, "media")
	}

	sdp := []string{}

	mLines, _ := session.ArrayCount("media")

	for i := 0; i < mLines; i++ {
		mLine, _ := session.ArrayElement(i, "media")
		if !mLine.Exists("payloads") {
			mLine.Set("", "payloads")
		}
	}

	for _, outType := range outerOrder {
		for _, rule := range rulesMap[outType] {
			if len(rule.Name) != 0 && session.Exists(rule.Name) && session.Path(rule.Name) != nil {
				lineStr := makeLine(outType, rule, session)
				sdp = append(sdp, lineStr)
			} else if len(rule.Push) > 0 && session.Exists(rule.Push) {
				count, err := session.ArrayCount(rule.Push)
				if err != nil {
					continue
				}

				for i := 0; i < count; i++ {
					el, _ := session.ArrayElement(i, rule.Push)
					lineStr := makeLine(outType, rule, el)
					sdp = append(sdp, lineStr)
				}

			}
		}
	}

	for i := 0; i < mLines; i++ {
		mLine, _ := session.ArrayElement(i, "media")
		lineStr := makeLine('m', rulesMap['m'][0], mLine)

		sdp = append(sdp, lineStr)

		for _, inType := range innerOrder {
			for _, rule := range rulesMap[inType] {
				if len(rule.Name) > 0 && session.Exists(rule.Name) && session.Path(rule.Name) != nil {
					lineStr := makeLine(inType, rule, mLine)
					sdp = append(sdp, lineStr)
				} else if len(rule.Push) > 0 && mLine.Exists(rule.Push) {
					count, err := mLine.ArrayCount(rule.Push)
					if err != nil {
						continue
					}

					for i := 0; i < count; i++ {
						el, _ := mLine.ArrayElement(i, rule.Push)
						lineStr := makeLine(inType, rule, el)
						sdp = append(sdp, lineStr)
					}
				}
			}
		}
	}

	sdpStr := strings.Join(sdp, "\r\n") + "\r\n"

	return sdpStr
}

func makeLine(otype byte, rule *Rule, location *gabs.Container) string {

	var format string
	if len(rule.Format) == 0 {
		if rule.FormatFunc != nil {
			var container *gabs.Container
			if len(rule.Push) != 0 {
				container = location
			} else {
				container = location.Path(rule.Name)
			}
			format = rule.FormatFunc(container)
		}
	} else {
		format = rule.Format
	}

	args := []interface{}{}

	if len(rule.Names) > 0 {

		for _, name := range rule.Names {
			if len(rule.Name) > 0 && location.Exists(rule.Name) && location.Exists(rule.Name, name) {
				fmt.Println("append", location.Search(rule.Name, name).Data())
				args = append(args, location.Search(rule.Name, name).Data())
			} else if location.Exists(name) {
				fmt.Println("append", location.Path(name).Data())
				args = append(args, location.Path(name).Data())
			} else {
				args = append(args, "")
			}
		}
	} else if location.Exists(rule.Name) {
		args = append(args, location.Path(rule.Name).Data())
	}

	fmt.Println(format, args)
	str := fmt.Sprintf(format, args)

	return str
}
