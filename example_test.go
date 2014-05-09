// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package config_test

import (
	"fmt"

	"code.minty.io/config"
)

func ExampleString() {
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
	fmt.Println(host, ok)
	// Output:
	// google.com true
}

func Examplerequired_String() {
	// for a `config.json` file like:
	/*
		{
			"host": "google.com",
			"links": {
				"google": "https://google.com"
			}
		}
	*/
	host := config.Required.String("host")
	fmt.Println(host)
	// Output:
	// google.com true
	//
	// panics when not found within `config.json`
}

func ExampleGroupString() {
	// for a `config.json` file like:
	/*
		{
			"host": "google.com",
			"links": {
				"google": "https://google.com"
			}
		}
	*/
	groupHost, ok := config.GroupString("links", "google")
	fmt.Println(groupHost, ok)
	// Output:
	// https://google.com true
}

func Examplerequired_GroupString() {
	// for a `config.json` file like:
	/*
		{
			"host": "google.com",
			"links": {
				"google": "https://google.com"
			}
		}
	*/
	groupHost := config.Required.GroupString("links", "google")
	fmt.Println(groupHost, ok)
	// Output:
	// https://google.com true
	//
	// panics when not found within `config.json`
}
