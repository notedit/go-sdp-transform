package sdptransform

import (
	"testing"
)

var sdpStr = `a=ice-ufrag:F7gI
a=ice-pwd:x9cml/YzichV2+XlhiMu8
a=fingerprint:sha-1 42:89:c5:c6:55:9d:6e:c8:e8:83:55:2a:39:f9:b6:eb:e9:a3:a9:e7
`

func TestParse(t *testing.T) {

	session, err := Parse([]byte(sdpStr))

	if err != nil {

		t.Error(err)
	}

	t.Log(session)
}
