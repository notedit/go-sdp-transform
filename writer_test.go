package sdptransform

import (
	"fmt"
	"testing"
)

func TestWrite(t *testing.T) {

	session, err := Parse([]byte(sdpStr))
	if err != nil {
		t.Error(err)
	}

	fmt.Println(session)

	ret := Write(session)

	fmt.Println(ret)

}
