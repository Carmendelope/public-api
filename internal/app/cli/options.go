/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package cli

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

const OptionsPath = "options"

type Options struct{
	// TODO Refactor and store the options in a single map that is serialized.
	Options map[string]string `json:"options"`
}

func NewOptions() * Options{
	return &Options{}
}

func (o * Options) getPath() string {
	path := filepath.Join(resolvePath(DefaultPath), OptionsPath)
	log.Debug().Str("path", path).Msg("Options directory")
	return path
}

func (o * Options) createIfNotExists(path string) {
	_ = os.MkdirAll(path, 0700)
}

// Set the value of a key as a persistent option.
func (o * Options) Set(key string, value string){

	if key == "" {
		log.Fatal().Msg("key must not be empty")
	}

	if value == "" {
		log.Fatal().Msg("value must not be empty")
	}

	basePath := o.getPath()
	o.createIfNotExists(basePath)
	targetPath := filepath.Join(basePath, key)
	err := ioutil.WriteFile(targetPath, []byte(value), 0600)
	if err != nil{
		log.Fatal().Err(err).Str("key", key).Msg("cannot write option")
	}
}

// Get the value of a previously store key
func (o * Options) Get(key string) string {

	if key == "" {
		log.Fatal().Msg("key must not be empty")
	}

	targetPath := filepath.Join(o.getPath(), key)
	log.Debug().Str("targetPath", targetPath).Msg("Get")
	value, err := ioutil.ReadFile(targetPath)
	if err != nil{
		log.Warn().Err(err).Str("key", key).Msg("cannot read option")
		return ""
	}
	log.Debug().Str("value", string(value)).Msg("Get")
	return string(value)
}

// Delete a given key
func (o * Options) Delete(key string){
	targetPath := filepath.Join(o.getPath(), key)
	_ = os.Remove(targetPath)
}

// List the available options
func (o * Options) List(){
	targetPath := o.getPath()
	_ = filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir(){
			retrieved := o.Get(filepath.Base(info.Name()))
			fmt.Printf("Key: %s Value: %s\n", info.Name(), retrieved)
		}
		return nil
	})
}

// Resolve the effective value of a parameter.
func (o * Options) Resolve(key string, paramValue string) string {
	log.Info().Str("key", key).Str("paramValue", paramValue).Msg("resolving option")
	stored := o.Get(key)
	if stored == "" {
		return paramValue
	}
	if stored != "" && paramValue == "" {
		log.Info().Str(key, stored).Msg("using stored option")
		return stored
	}
	return paramValue
}

func (o * Options) ResolveAsInt(key string, paramValue int) int {
	toStr := ""
	if paramValue >= 0 {
		toStr = strconv.Itoa(paramValue)
	}
	value, err := strconv.Atoi(o.Resolve(key, toStr))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot convert value to int")
	}
	return value
}
