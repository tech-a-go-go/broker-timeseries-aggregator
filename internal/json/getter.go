package json

import (
	"github.com/valyala/bytebufferpool"
)

func GetStringInJson(jsonBytes []byte, param []byte, startOffset int, buf []byte) (stringValue []byte, endIndex int, ok bool) {
	jsonLen := len(jsonBytes)
	if startOffset >= jsonLen-1 {
		return
	}
	// table
	paramLen := len(param)
	for i := startOffset; i < jsonLen; i++ {
		// まず " と ":" を見つける
		if i-paramLen-3 >= 0 && jsonBytes[i-paramLen-3] == '"' /* param word */ && jsonBytes[i-2] == '"' && jsonBytes[i-1] == ':' && jsonBytes[i] == '"' {
			matched := true
			for k := 0; k < paramLen; k++ {
				if jsonBytes[i-paramLen-2+k] != param[k] {
					matched = false
					break
				}
			}
			if !matched {
				continue
			}
			for j := i + 1; j < jsonLen; j++ {
				if jsonBytes[j] != '"' {
					buf = append(buf, jsonBytes[j])
					//buf.WriteByte(jsonBytes[j])
				} else {
					if len(buf) > 0 {
						ok = true
						stringValue = buf
						endIndex = j
						break
					}
					break
				}
			}
			return
		}
	}
	return
}

func GetNumberInJson(jsonBytes []byte, param []byte, startOffset int, buf *bytebufferpool.ByteBuffer) (numberValue []byte, endIndex int, ok bool) {
	jsonLen := len(jsonBytes)
	if startOffset >= jsonLen-1 {
		return
	}
	// table
	paramLen := len(param)
	for i := startOffset; i < jsonLen; i++ {
		// まず " と ": を見つける
		if i-paramLen-2 >= 0 && jsonBytes[i-paramLen-2] == '"' /* param word */ && jsonBytes[i-1] == '"' && jsonBytes[i] == ':' {
			matched := true
			for k := 0; k < paramLen; k++ {
				if jsonBytes[i-paramLen-1+k] != param[k] {
					matched = false
					break
				}
			}
			if !matched {
				continue
			}
			numBytes := bytebufferpool.Get()
			for j := i + 1; j < jsonLen; j++ {
				if ('0' <= jsonBytes[j] && jsonBytes[j] <= '9') || jsonBytes[j] == '.' || jsonBytes[j] == '+' || jsonBytes[j] == '-' || jsonBytes[j] == 'e' || jsonBytes[j] == 'E' {
					numBytes.WriteByte(jsonBytes[j])
				} else {
					if numBytes.Len() > 0 {
						ok = true
						numberValue = numBytes.B
						endIndex = j
						break
					}
					break
				}
			}
			bytebufferpool.Put(numBytes)
			return
		}
	}
	return
}
