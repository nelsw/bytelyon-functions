package pw

import (
	"bytelyon-functions/pkg/service/s3"
	"encoding/json"

	"github.com/playwright-community/playwright-go"
)

func GetStorageState() (*playwright.OptionalStorageState, error) {
	var storageState playwright.OptionalStorageState
	if data, err := s3.New().Get("pw/storage-state/_.json"); err != nil {
		return nil, err
	} else if err = json.Unmarshal(data, &storageState); err != nil {
		return nil, err
	}
	return &storageState, nil
}

func SaveStorageState(state *playwright.StorageState) error {
	if b, err := json.MarshalIndent(state, "", "\t"); err != nil {
		return err
	} else if err = s3.New().Put("pw/storage-state/_.json", b); err != nil {
		return err
	}
	return nil
}
