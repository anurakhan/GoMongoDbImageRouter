package models

import (
	"encoding/hex"
)

type FileId struct {
	Id string
}

func (fileId *FileId) ToHex() string {
	ret, err := hex.DecodeString(fileId.Id)
	if err != nil {
		panic(err)
	}
	return string(ret)
}
