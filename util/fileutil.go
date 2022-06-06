//
// Copyright 2021 Tim Miller.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package util

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
)

// I don't like this here.  TODO: Move elsewhere.
func DownloadFile(fileurl string, destinationdir string) (*string, error) {
	resp, err := http.Get(fileurl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	u, err := url.Parse(fileurl)
	if err != nil {
		return nil, err
	}
	filename := path.Base(u.Path)
	fullpath := path.Join(destinationdir, filename)

	err = os.MkdirAll(destinationdir, os.ModePerm)
	if err != nil {
		return nil, err
	}

	out, err := os.Create(fullpath)
	if err != nil {
		return nil, err
	}

	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return &fullpath, nil
}

func ReadFile(filepath string) (*string, error) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	str := string(b)
	return &str, nil
}
