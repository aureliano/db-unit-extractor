package dataconv

import (
	"encoding/base64"
	"fmt"
)

type BlobConverter struct{}

func (BlobConverter) Convert(_ string, source interface{}) (interface{}, error) {
	bts, isArrByte := source.([]byte)
	if !isArrByte {
		return nil, fmt.Errorf("'%v' is not []byte", source)
	}

	return base64.StdEncoding.EncodeToString(bts), nil
}

func (BlobConverter) Handle(vl interface{}) bool {
	_, handled := vl.([]byte)
	return handled
}
