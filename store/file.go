package store

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
)

func StoreFile(filePath string, b []byte) error {
	err := ioutil.WriteFile(filePath, b, 0644)
	if err != nil {
		log.Printf("WriteFile failed err:%v\n", err)
		return fmt.Errorf("StoreFile failed")
	}
	return nil
}

func MkdirAll(filePath string) error{
	if _, err := os.Stat(path.Dir(filePath)); os.IsNotExist(err) {
		if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
			log.Fatalf("MkdirAll failed err:%v", err)
			return err
		}
	}
	return nil
}
