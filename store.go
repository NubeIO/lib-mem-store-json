package storage

import (
	"encoding/json"
	"fmt"
	"os"
)

type store map[string]interface{}

type StorageHandler interface {
	RegisterData(id string, dataPtr any) error
	WriteToDisk() error
	ReadFromDisk() error
}

type storageInst struct {
	StorageMap store
	Filename   string
	RawData    map[string]json.RawMessage
}

// NewStorage Create a new StorageHandler with the absolute path filename as the disk backup location
func NewStorage(filename string) StorageHandler {
	return &storageInst{
		StorageMap: make(store),
		Filename:   filename,
		RawData:    make(map[string]json.RawMessage),
	}
}

// RegisterStruct keeps a pointer to any data in dataPtr under the id
func (inst *storageInst) RegisterData(id string, dataPtr any) error {
	inst.StorageMap[id] = dataPtr
	data, ok := inst.RawData[id]
	if !ok {
		return fmt.Errorf("DB no data existing for key: %s. Most likely a new database", id)
	}
	err := json.Unmarshal(data, dataPtr)
	if err != nil {
		return fmt.Errorf("DB error unmarshalling node: %s, %v, %w", id, dataPtr, err)
	}
	delete(inst.RawData, id)
	return nil
}

// WriteToDisk Write all storage components to disk
func (inst *storageInst) WriteToDisk() error {
	raw, err := json.Marshal(inst.StorageMap)
	if err != nil {
		return fmt.Errorf("db error writting to file: %w", err)
	}

	// Create and write temporary file
	tempFile, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		return fmt.Errorf("db error writting to file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	if _, err := tempFile.Write(raw); err != nil {
		tempFile.Close()
		return fmt.Errorf("db error writting to file: %w", err)
	}

	return os.Rename(tempFile.Name(), inst.Filename)
}

// ReadFromDisk Read in storage from disk
func (inst *storageInst) ReadFromDisk() error {
	fileContent, err := os.ReadFile(inst.Filename)

	if err != nil {
		if os.IsNotExist(err) {
			err = os.WriteFile(inst.Filename, []byte{}, 0644)
			if err != nil {
				return fmt.Errorf("db error creating file: %w", err)
			}
			return nil
		}
		return err
	}

	if len(fileContent) == 0 {
		return nil
	}
	err = json.Unmarshal(fileContent, &inst.RawData)
	if err != nil {
		return fmt.Errorf("db error unmarshalling raw file: %w", err)
	}
	return nil
}
