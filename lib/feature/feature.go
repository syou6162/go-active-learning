package feature

import (
	"encoding/json"
)

type FeatureVector []string

func (fv *FeatureVector) MarshalBinary() ([]byte, error) {
	json, err := json.Marshal(fv)
	if err != nil {
		return nil, err
	}
	return []byte(json), nil
}

func (fv *FeatureVector) UnmarshalBinary(data []byte) error {
	err := json.Unmarshal(data, fv)
	if err != nil {
		return err
	}
	return nil
}
