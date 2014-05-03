// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config_test

import "bitbucket.org/juztin/config"

func main() {
	// for a `config.json` file like:
	/*
		{
			"host": "google.com",
			"links": {
				"google": "https://google.com"
			}
		}
	*/
	host, ok := config.String("host")
	if !ok {
		host = "localhost"
	}

	// panics when not found within `config.json`
	host = config.Required.String("host")

	groupHost, ok := config.GroupString("links", "google")
	if !ok {
		groupHost = "localhost"
	}
}
