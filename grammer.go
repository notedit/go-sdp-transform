package sdptransform

import (
	"regexp"
)

type Rule struct {
	Name       string
	Push       string
	Reg        *regexp.Regexp
	Names      []string
	Types      []rune
	Format     string
	FormatFunc interface{}
}

var RulesMap map[byte][]*Rule = map[byte][]*Rule{
	'v': []*Rule{
		&Rule{
			Name:   "version",
			Push:   "",
			Reg:    regexp.MustCompile("^(\\d*)$"),
			Names:  []string{},
			Types:  []rune{'d'},
			Format: "%d",
		},
	},
	'o': []*Rule{
		&Rule{
			Name:   "origin",
			Push:   "",
			Reg:    regexp.MustCompile("^(\\S*) (\\d*) (\\d*) (\\S*) IP(\\d) (\\S*)"),
			Names:  []string{"username", "sessionId", "sessionVersion", "netType", "ipVer", "address"},
			Types:  []rune{'s', 'd', 'd', 's', 'd', 's'},
			Format: "%s %d %d %s IP%d %s",
		},
	},
	's': []*Rule{
		&Rule{
			Name:   "name",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	'i': []*Rule{
		&Rule{
			Name:   "description",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	'u': []*Rule{
		&Rule{
			Name:   "uri",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	'e': []*Rule{
		&Rule{
			Name:   "email",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	'p': []*Rule{
		&Rule{
			Name:   "phone",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	'z': []*Rule{
		&Rule{
			Name:   "timezones",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	'r': []*Rule{
		&Rule{
			Name:   "repeats",
			Push:   "",
			Reg:    regexp.MustCompile("(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "%s",
		},
	},
	't': []*Rule{
		&Rule{
			Name:   "timing",
			Push:   "",
			Reg:    regexp.MustCompile("^(\\d*) (\\d*)"),
			Names:  []string{"start", "stop"},
			Types:  []rune{'d', 'd'},
			Format: "%d %d",
		},
	},
	'c': []*Rule{
		&Rule{
			Name:   "connection",
			Push:   "",
			Reg:    regexp.MustCompile("^IN IP(\\d) (\\S*)"),
			Names:  []string{"version", "ip"},
			Types:  []rune{'d', 's'},
			Format: "IN IP%d %s",
		},
	},
	'b': []*Rule{
		&Rule{
			Name:   "",
			Push:   "bandwidth",
			Reg:    regexp.MustCompile("^(TIAS|AS|CT|RR|RS):(\\d*)"),
			Names:  []string{"type", "limit"},
			Types:  []rune{'s', 'd'},
			Format: "%s:%d",
		},
	},
	'm': []*Rule{ // m=video 51744 RTP/AVP 126 97 98 34 31
		&Rule{
			Name:   "",
			Push:   "",
			Reg:    regexp.MustCompile("^(\\w*) (\\d*) ([\\w\\/]*)(?: (.*))?"),
			Names:  []string{"type", "port", "protocal", "payloads"},
			Types:  []rune{'s', 'd', 's', 's'},
			Format: "%s %d %s %s",
		},
	},
	'a': []*Rule{ // a=rtpmap:110 opus/48000/2
		&Rule{
			Name:       "",
			Push:       "rtp",
			Reg:        regexp.MustCompile("^rtpmap:(\\d*) ([\\w\\-\\.]*)(?:\\s*\\/(\\d*)(?:\\s*\\/(\\S*))?)?"),
			Names:      []string{"playload", "codec", "rate", "encoding"},
			Types:      []rune{'d', 's', 'd', 's'},
			Format:     "",
			FormatFunc: nil,
		},
		// a=fmtp:108 profile-level-id=24;object=23;bitrate=64000
		// a=fmtp:111 minptime=10; useinbandfec=1
		&Rule{
			Name:   "",
			Push:   "fmtp",
			Reg:    regexp.MustCompile("^fmtp:(\\d*) ([\\S| ]*)"),
			Names:  []string{"payload", "config"},
			Types:  []rune{'d', 's'},
			Format: "fmtp:%d %s",
		},
		// a=control:streamid=0
		&Rule{
			Name:   "control",
			Push:   "",
			Reg:    regexp.MustCompile("^control:(.*)"),
			Names:  []string{},
			Types:  []rune{'s'},
			Format: "controle:%s",
		},
		// a=rtcp:65179 IN IP4 193.84.77.194
		&Rule{
			Name:       "rtcp",
			Push:       "",
			Reg:        regexp.MustCompile("^rtcp:(\\d*)(?: (\\S*) IP(\\d) (\\S*))?"),
			Names:      []string{"port", "netType", "ipVer", "address"},
			Types:      []rune{'d', 's', 'd', 's'},
			Format:     "",
			FormatFunc: nil,
		},
		// a=rtcp-fb:98 trr-int 100
		&Rule{
			Name:   "",
			Push:   "rtcpFbTrrInt",
			Reg:    regexp.MustCompile("^rtcp-fb:(\\*|\\d*) trr-int (\\d*)"),
			Names:  []string{"payload", "value"},
			Types:  []rune{'s', 'd'},
			Format: "rtcp-fb:%s trr-int %d",
		},
		// a=rtcp-fb:98 nack rpsi
		&Rule{
			Name:       "",
			Push:       "rtcpFb",
			Reg:        regexp.MustCompile("^rtcp-fb:(\\*|\\d*) ([\\w\\-_]*)(?: ([\\w\\-_]*))?"),
			Names:      []string{"payload", "type", "subtype"},
			Types:      []rune{'s', 's', 's'},
			Format:     "",
			FormatFunc: nil,
		},
		// a=extmap:2 urn:ietf:params:rtp-hdrext:toffset
		// a=extmap:1/recvonly URI-gps-string
		&Rule{
			Name:       "",
			Push:       "ext",
			Reg:        regexp.MustCompile("^extmap:(\\d+)(?:\\/(\\w+))? (\\S*)(?: (\\S*))?"),
			Names:      []string{"value", "direction", "uri", "config"},
			Types:      []rune{'d', 's', 's', 's'},
			Format:     "",
			FormatFunc: nil,
		},
		// a=crypto:1 AES_CM_128_HMAC_SHA1_80 inline:PS1uQCVeeCFCanVmcjkpPywjNWhcYD0mXXtxaVBR|2^20|1:32
	},
}
