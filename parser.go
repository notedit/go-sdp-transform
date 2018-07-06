package sdptransform

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"

	"github.com/Jeffail/gabs"
	json "github.com/bitly/go-simplejson"
)

const maxLineSize = 1024
const validLineStr = "^([a-z])=(.*)"

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func Parse(sdp []byte) (obj *json.Json, err error) {

	validLineRegex := regexp.MustCompile(validLineStr)
	bufferReader := bufio.NewReader(bytes.NewReader(sdp))
	container := gabs.New()
	medias, _ := container.Array("media")

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

				medias.ArrayAppend(m)

				//todo location the media
			}
		}

	}

}
