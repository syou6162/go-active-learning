package model

import (
	"encoding/json"
)

type LabelType int

func (lt *LabelType) MarshalBinary() ([]byte, error) {
	return json.Marshal(lt)
}

func (lt *LabelType) UnmarshalBinary(data []byte) error {
	if err := json.Unmarshal(data, &lt); err != nil {
		return err
	}
	return nil
}

const (
	POSITIVE  LabelType = 1
	NEGATIVE  LabelType = -1
	UNLABELED LabelType = 0
)
