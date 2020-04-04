package fcw

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func BasicFirefox(userdir string, private bool, args ...string) (UI, error) {
	if !private {
		os.MkdirAll(userdir, os.ModePerm)
	} else {
		add := true
		for _, arg := range args {
			if arg == "--private-window" {
				add = false
			}
		}
		if add {
			args = append(args, "--private-window")
			args = append(args, "about:blank")
		}
	}
	userdir, err := filepath.Abs(userdir)
	if err != nil {
		return nil, err
	}
	return NewFirefox("", userdir, 800, 600, args...)
}

func ExtendedFirefox(userdir string, private bool, extensiondirs []string, args ...string) (UI, error) {
	var extensionArgs []string
	for _, extension := range extensiondirs {
		if _, err := os.Stat(extension); err != nil {
			log.Println("extension load warning,", err)
		}
	}
	args = append(args, extensionArgs...)
	return BasicFirefox(userdir, private, args...)
}

func SecureExtendedFirefox(userdir string, private bool, extensionxpis, extensionhashes []string, args ...string) (UI, error) {
	var extensionArgs []string
	if len(extensionxpis) != len(extensionhashes) {
		return nil, fmt.Errorf("hash list is different from extension XPI list")
	}
	for index, extension := range extensionxpis {
		if _, err := os.Stat(userdir + "/extensions/" + extension); err != nil {
			return nil, err
		}
		if bytes, err := ioutil.ReadFile(userdir + "/extensions/" + extension); err == nil {
			hash := sha256.Sum256(bytes)
			hexed := hex.EncodeToString(hash[:])
			if extensionhashes[index] != hexed {
				return nil, fmt.Errorf("hash mismatch error on extension %s \n'%s' \n!= \n'%s'", userdir+"/extensions/"+extension, hexed, extensionhashes[index])
			}
			log.Printf("hash match on extension %s \n'%s' \n == \n'%s'", userdir+"/extensions/"+extension, hexed, extensionhashes[index])
		} else {
			return nil, fmt.Errorf("hash calculation error on extension %s %s", userdir+"/extensions/"+extension, err.Error())
		}
	}
	args = append(args, extensionArgs...)
	return BasicFirefox(userdir, private, args...)
}

var FIREFOX, ERROR = BasicFirefox("basic", true, "--headless")

func Run() error {
	if ERROR != nil {
		return ERROR
	}
	defer FIREFOX.Close()
	<-FIREFOX.Done()
	return nil
}
