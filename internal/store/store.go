package store

import (
	"aviatus/vividora/internal/config"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var kv *KeyValueStore
var mutex sync.RWMutex

type Store map[string]string
type KeyValueStore struct {
	store Store
}

func NewKeyValueStore() (*KeyValueStore, error) {
	files, err := ioutil.ReadDir(config.StoragePath)
	if err != nil {
		return nil, err
	}

	size := len(files)
	return &KeyValueStore{
		store: make(Store, size),
	}, nil
}

func Set(key string, value string, isRestore bool) error {
	mutex.Lock()
	defer mutex.Unlock()
	_, exists := kv.store[key]
	kv.store[key] = value

	if !isRestore {
		err := UpdateStorePersist(key, value, exists)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
	}
	return nil
}

func Get(key string) (string, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	value, ok := kv.store[key]
	return value, ok
}

func Delete(key string) error {
	mutex.Lock()
	defer mutex.Unlock()
	delete(kv.store, key)
	err := DeleteStorePersist(key)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func getStore() map[string]string {
	mutex.RLock()
	defer mutex.RUnlock()
	return kv.store
}

func RestoreFromSnapshot(snapshotPath string) error {
	file, err := os.Open(snapshotPath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer file.Close()

	// Create a decoder
	decoder := gob.NewDecoder(file)

	err = decoder.Decode(&kv.store)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func TakeSnapshot() error {
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02-15:04:05")

	// Create a buffer to hold the encoded binary data
	buffer := new(bytes.Buffer)

	// Create an encoder
	encoder := gob.NewEncoder(buffer)

	// Encode the map
	err := encoder.Encode(getStore())
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	// Get the binary data from the buffer
	binaryData := buffer.Bytes()

	// Create the directory path if it doesn't exist
	err = os.MkdirAll(config.SnapshotPath, 0755)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	file, err := os.Create(config.SnapshotPath + "snapshot-" + formattedTime + ".vdb")
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	defer file.Close()

	// Encode and write the map to the file
	err = binary.Write(file, binary.LittleEndian, binaryData)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	currentTime = time.Now()
	formattedTime = currentTime.Format("2006-01-02.15:04:05")
	fmt.Println("Successfull Snapshot at: ", formattedTime)
	return nil
}

func UpdateStorePersist(key string, value string, exist bool) error {
	hash := sha256.New()
	hash.Write([]byte(key))
	hashSum := hash.Sum(nil)
	hashString := hex.EncodeToString(hashSum)
	path := config.StoragePath + hashString

	switch exist {
	case true:
		err := ioutil.WriteFile(path, []byte(key+" "+value), 0644)
		if err != nil {
			log.Fatal(err)
		}
	case false:
		file, err := os.Create(path)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
		defer file.Close()

		// Encode and write the map to the file
		err = binary.Write(file, binary.LittleEndian, []byte(key+" "+value))
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}
	}
	return nil
}

func DeleteStorePersist(key string) error {
	hash := sha256.New()
	hash.Write([]byte(key))
	hashSum := hash.Sum(nil)
	hashString := hex.EncodeToString(hashSum)
	path := config.StoragePath + hashString

	err := os.Remove(path)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}
	return nil
}

func RestoreFromStorage() error {
	fmt.Println("Restoring from storage is started...")
	files, err := ioutil.ReadDir(config.StoragePath)
	if err != nil {
		log.Fatal(err)
		return err
	}

	for _, f := range files {
		file, err := os.Open(config.StoragePath + f.Name())
		if err != nil {
			fmt.Println("Error File:", err)
			return err
		}

		stat, err := file.Stat()
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		fileSize := stat.Size()
		buffer := make([]byte, fileSize)

		_, err = file.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return err
		}

		str := strings.Split(string(buffer), " ")
		Set(str[0], str[1], true)
		file.Close()
	}
	fmt.Println("Restoring from storage is successfully finished...")
	return nil
}

func StartStore() error {
	var err error
	kv, err = NewKeyValueStore()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	err = os.MkdirAll(config.StoragePath, 0755)
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	err = RestoreFromStorage()
	if err != nil {
		fmt.Println("Error:", err)
		return err
	}

	return nil
}
