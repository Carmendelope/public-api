/*
 * Copyright (C)  2018 Nalej - All Rights Reserved
 */

package options

import (
	"fmt"
	"github.com/nalej/public-api/internal/app/output"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

// DefaultPath to store and retrieve credentials
const DefaultPath = "~/.nalej/"

// OptionsPath with the path inside DefaultPath where options are stored.
const OptionsPath = "options"

// APIAddressKey with the name of the key that points to the API address
const APIAddressKey= "nalejAddress"
// APIAddressPrefix with the prefix for API address.
const APIAddressPrefix = "api."
// LoginAddressKey with the name of the key that points to the Login API address
const LoginAddressKey = "loginAddress"
// LoginAddressPrefix with the prefix for the login API address.
const LoginAddressPrefix = "login."

type Options struct {
	// TODO Refactor and store the options in a single map that is serialized.
	Options map[string]string `json:"options"`
}

func NewOptions() *Options {
	return &Options{}
}

func (o *Options) getPath() string {
	path := filepath.Join(GetPath(DefaultPath), OptionsPath)
	log.Debug().Str("path", path).Msg("Options directory")
	return path
}

func (o *Options) createIfNotExists(path string) {
	_ = os.MkdirAll(path, 0700)
}

// Set the value of a key as a persistent option.
func (o *Options) Set(key string, value string) {

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
	if err != nil {
		log.Fatal().Err(err).Str("key", key).Msg("cannot write option")
	}
}

// Get the value of a previously store key
func (o *Options) Get(key string) string {

	if key == "" {
		log.Fatal().Msg("key must not be empty")
	}

	targetPath := filepath.Join(o.getPath(), key)
	log.Debug().Str("targetPath", targetPath).Msg("Get")
	value, err := ioutil.ReadFile(targetPath)
	if err != nil {
		log.Debug().Err(err).Str("key", key).Msg("cannot read option")
		return ""
	}
	log.Debug().Str("value", string(value)).Msg("Get")
	return string(value)
}

// Delete a given key
func (o *Options) Delete(key string) {
	targetPath := filepath.Join(o.getPath(), key)
	_ = os.Remove(targetPath)
}

// List the available options
func (o *Options) List() {
	header := []string{"KEY", "VALUE"}
	values := make([][]string, 0)
	targetPath := o.getPath()
	_ = filepath.Walk(targetPath, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			retrieved := o.Get(filepath.Base(info.Name()))
			values = append(values, []string{info.Name(), retrieved})
		}
		return nil
	})
	output.PrintFromValues(header, values)
}

// Resolve the effective value of a parameter as string.
func (o *Options) Resolve(key string, paramValue string) string {
	log.Debug().Str("key", key).Str("paramValue", paramValue).Msg("resolving option")
	stored := o.Get(key)
	if stored == "" {
		return paramValue
	}
	if stored != "" && paramValue == "" {
		log.Debug().Str(key, stored).Msg("using stored option")
		return stored
	}
	return paramValue
}

// ResolveAsInt resolves the value of an option as int.
func (o *Options) ResolveAsInt(key string, paramValue int) int {
	toStr := ""
	if paramValue > 0 {
		toStr = strconv.Itoa(paramValue)
	}

	res := o.Resolve(key, toStr)
	if res == "" {
		return 0
	}
	value, err := strconv.Atoi(res)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot convert value to int")
	}
	return value
}

// UpdatePlatformAddress updates both the api endpoints for the login and the public api.
func (o *Options) UpdatePlatformAddress(newBaseAddress string) []string {
	if strings.HasPrefix(newBaseAddress, APIAddressPrefix) || strings.HasPrefix(newBaseAddress, LoginAddressPrefix){
		log.Fatal().Msg("expecting new base address without login. or api. prefixes")
	}
	apiAddress := fmt.Sprintf("%s%s",APIAddressPrefix, newBaseAddress)
	loginAddress := fmt.Sprintf("%s%s",LoginAddressPrefix, newBaseAddress)

	o.Set(APIAddressKey, apiAddress)
	o.Set(LoginAddressKey, loginAddress)

	return []string{apiAddress, loginAddress}
}


// GetPath resolves a given path by adding support for relative paths.
func GetPath(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, _ := user.Current()
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}
	if strings.HasPrefix(path, "../"){
		abs, _ := filepath.Abs("../")
		return strings.Replace(path, "..", abs, 1)
	}
	if strings.HasPrefix(path, "."){
		abs, _ := filepath.Abs("./")
		return strings.Replace(path, ".", abs, 1)
	}
	return path
}
