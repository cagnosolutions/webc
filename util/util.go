package util

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
)

func GetUser(r *http.Request) string {
	return base64.StdEncoding.EncodeToString([]byte(r.RemoteAddr + r.UserAgent()))
}

func UUID4() string {
	u := make([]byte, 16)
	if _, err := rand.Read(u[:16]); err != nil {
		log.Println(err)
	}
	u[8] = (u[8] | 0x80) & 0xbf
	u[6] = (u[6] | 0x40) & 0x4f
	return fmt.Sprintf("%x-%x-%x-%x-%x", u[:4], u[4:6], u[6:8], u[8:10], u[10:])
}

func SliceString(s string, b byte) []string {
	var ss []string
	if len(s) > 1 {
		i := bytes.IndexByte([]byte(s), b)
		for ; i > -1; i = bytes.IndexByte([]byte(s), b) {
			ss = append(ss, s[0:i])
			s = s[i+1 : len(s)]
		}
		ss = append(ss, s[0:len(s)])
	} else {
		ss = append(ss, s)
	}
	return ss
}

/*
func SliceString(s string, b byte) []string {
	var ss []string
	if len(s) > 1 {
		i := bytes.IndexByte([]byte(s), b)
		for ; i > -1; i = bytes.IndexByte([]byte(s), b) {
			if i > 0 {
				ss = append(ss, s[0:i])
			}
			s = s[i+1 : len(s)]
		}
		ss = append(ss, s[0:len(s)])
	} else {
		ss = append(ss, s)
	}
	return ss
}
*/
