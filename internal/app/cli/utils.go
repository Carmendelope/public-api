/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cli

import (
	"encoding/base64"
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os/user"
	"path/filepath"
	"strings"
)

func GetLabels(rawLabels string) map[string]string {
	labels := make(map[string]string, 0)

	split := strings.Split(rawLabels, ";")
	for _, l := range split {
		ls := strings.Split(l, ":")
		if len(ls) != 2 {
			log.Fatal().Str("label", l).Msg("malformed label, expecting key:value")
		}
		labels[ls[0]] = ls[1]
	}
	return labels
}


// GetPath resolves a given path by adding support for relative paths.
func GetPath(path string) string {
	if strings.HasPrefix(path, "~") {
		usr, _ := user.Current()
		return strings.Replace(path, "~", usr.HomeDir, 1)
	}
	if strings.HasPrefix(path, "../") {
		abs, _ := filepath.Abs("../")
		return strings.Replace(path, "..", abs, 1)
	}
	if strings.HasPrefix(path, ".") {
		abs, _ := filepath.Abs("./")
		return strings.Replace(path, ".", abs, 1)
	}
	return path
}

// PhotoToBase64 reads a image an convert the content to a base64 string
func PhotoToBase64(path string) (string, derrors.Error) {
	// if there is no path -> empty image
	if path == ""  {
		return "", nil
	}

	convertedPath := GetPath(path)
	content, err := ioutil.ReadFile(convertedPath)
	if err != nil {
		return "", derrors.AsError(err, "cannot read descriptor")
	}
	// convert the buffer bytes to base64 string - use buf.Bytes() for new image
	imgBase64Str := base64.StdEncoding.EncodeToString(content)

	return imgBase64Str, nil
}