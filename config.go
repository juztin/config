// Copyright 2013 Justin Wilson. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package config provides a very basic JSON config file reader.
// By default it reads the file `config.json` during init. The name
// of the file can be overriden, and the file can also be re-loaded
// when necessary.
// All values are stored in memory and can be looked up, or overriden
// to a different value. Changes are not persisted.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type required struct{}

// ConfigFile is the name of the file to read configuration from.
var ConfigFile = "config.json"

// Required is has the same methods as this package, but panics when the keys are don't exist.
var Required = *new(required)

var (
	cfg    map[string]interface{}
	loaded = false
)

func init() {
	Load()
}

// Load reads the `config.json` file, within the current executing directory, and
// loads it's data.
// This is called on import automatically so there is usually no need to call this directory.
func Load() error {
	if loaded {
		return nil
	}

	// get|read configuration from file
	p, c, err := getConfig()
	if err != nil {
		return errors.New(fmt.Sprintf("failed to load configuration file: %s, from: %s\n%v", ConfigFile, p, err))
	}

	var j interface{}
	if err := json.Unmarshal(c, &j); err != nil {
		return errors.New(fmt.Sprintf("failed to read configuration file: %s, from: %s\n%v", ConfigFile, p, err))
	}

	cfg = j.(map[string]interface{})
	loaded = true
	return nil
}

// Reload will force the reloading/reading of the config file
func Reload() error {
	loaded = false
	return Load()
}

func getConfig() (p string, c []byte, e error) {
	p = filepath.Dir(os.Args[0])
	f := filepath.Join(p, ConfigFile)

	// if a config file exists within the executables path
	if _, err := os.Stat(f); err == nil {
		c, e = ioutil.ReadFile(f)
		return
	}

	// if a config file exists within the current working dir
	if p, e = os.Getwd(); e == nil {
		f = filepath.Join(p, ConfigFile)
		if _, e = os.Stat(f); e == nil {
			c, e = ioutil.ReadFile(f)
			return
		}
	}

	// no configuration was found
	p = ""
	e = errors.New(fmt.Sprintf("failed to find a configuration file: %s", ConfigFile))

	return
}

// accessors
func colBool(key string, col map[string]interface{}) (bool, bool) {
	if v, ok := col[key]; ok {
		b, ok := v.(bool)
		return b, ok
	}
	return false, false
}

func colString(key string, col map[string]interface{}) (string, bool) {
	if v, ok := col[key]; ok {
		s, ok := v.(string)
		return s, ok
	}
	return *new(string), false
}

func colInt(key string, col map[string]interface{}) (int, bool) {
	if v, ok := col[key]; ok {
		switch v.(type) {
		case int:
			return v.(int), true
		case float64:
			return int(v.(float64)), true
		}
	}
	return -1, false
}

func colFloat64(key string, col map[string]interface{}) (float64, bool) {
	if v, ok := col[key]; ok {
		switch v.(type) {
		case float64:
			return v.(float64), true
		case int:
			return float64(v.(int)), true
		}
	}
	return -1.0, false
}

func colVal(key string, col map[string]interface{}) (interface{}, bool) {
	if v, ok := col[key]; ok {
		return v, true
	}
	return nil, false
}

// Bool returns the boolean value for the `key` within the root level.
// The value, or default value, is returned along with boolean of wether the key was found.
func Bool(key string) (bool, bool) {
	return colBool(key, cfg)
}

// String returns the string value for the `key` within the root level.
// The value, or default value, is returned along with boolean of wether the key was found.
func String(key string) (string, bool) {
	return colString(key, cfg)
}

// Int returns the int value for the `key` within the root level.
// The value, or default value, is returned along with boolean of wether the key was found.
func Int(key string) (int, bool) {
	return colInt(key, cfg)
}

// Float64 returns the float64 value for the `key` within the root level.
// The value, or default value, is returned along with boolean of wether the key was found.
func Float64(key string) (float64, bool) {
	return colFloat64(key, cfg)
}

// Val returns the value, as an interface{}, for the `key` within the root level.
// The value, or nil, is returned along with boolean of wether the key was found.
func Val(key string) (interface{}, bool) {
	return colVal(key, cfg)
}

// GroupBool returns the boolean value for the `key` within the group level.
// The boolean, or false, is returned along with boolean of wether the key was found.
func GroupBool(group, key string) (v bool, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colBool(key, col)
		}
	}
	return
}

// GroupBool returns the boolean value for the `key` within the group level.
// The string, or empty string, is returned along with boolean of wether the key was found.
func GroupString(group, key string) (v string, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colString(key, col)
		}
	}
	return
}

// GroupBool returns the boolean value for the `key` within the group level
// The int, or 0, is returned along with boolean of wether the key was found.
func GroupInt(group, key string) (v int, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colInt(key, col)
		}
	}
	return
}

// GroupBool returns the boolean value for the `key` within the group level
// The float64, or 0, is returned along with boolean of wether the key was found.
func GroupFloat64(group, key string) (v float64, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colFloat64(key, col)
		}
	}
	return
}

// GroupVal returns the value, as an interface{}, for the `key` within the group level
// The value, or nil, is returned along with boolean of wether the key was found.
func GroupVal(group, key string) (v interface{}, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colVal(key, col)
		}
	}
	return
}

// SetInt sets the value for given key, within the root.
func SetInt(key string, val int) {
	if _, ok := Int(key); ok {
		cfg[key] = val
	}
}

// SetFloat64 sets the value for given key, within the root.
func SetFloat64(key string, val float64) {
	if _, ok := Float64(key); ok {
		cfg[key] = val
	}
}

// SetBool sets the value for given key, within the root.
func SetBool(key string, val bool) {
	if _, ok := Bool(key); ok {
		cfg[key] = val
	}
}

// SetString sets the value for given key, within the root.
func SetString(key string, val string) {
	if _, ok := String(key); ok {
		cfg[key] = val
	}
}

// SetGroupInt sets the value for given key, within the group.
func SetGroupInt(group, key string, val int) {
	if _, ok := GroupInt(group, key); ok {
		cfg[group].(map[string]interface{})[key] = val
	}
}

// SetGroupFloat64 sets the value for given key, within the group.
func SetGroupFloat64(group, key string, val float64) {
	if _, ok := GroupFloat64(group, key); ok {
		cfg[group].(map[string]interface{})[key] = val
	}
}

// SetGroupBool sets the value for given key, within the group.
func SetGroupBool(group, key string, val bool) {
	if _, ok := GroupBool(group, key); ok {
		cfg[group].(map[string]interface{})[key] = val
	}
}

// SetGroupString sets the value for given key, within the group.
func SetGroupString(group, key string, val string) {
	if _, ok := GroupString(group, key); ok {
		cfg[group].(map[string]interface{})[key] = val
	}
}

// Bool returns the boolean value, within the root, and panics when not found.
func (r required) Bool(key string) bool {
	b, ok := Bool(key)
	if !ok {
		panic(fmt.Sprintf("failed to retrieve '%s' bool from config", key))
	}
	return b
}

// String returns the string, within the root, and panics when not found.
func (r required) String(key string) string {
	s, ok := String(key)
	if !ok {
		panic(fmt.Sprintf("failed to retrieve '%s' string from config", key))
	}
	return s
}

// Int returns the int, within the root, and panics when not found.
func (r required) Int(key string) int {
	i, ok := Int(key)
	if !ok {
		panic(fmt.Sprintf("failed to retrieve '%s' int from config", key))
	}
	return i
}

// Float64 returns the float64, within the root, and panics when not found.
func (r required) Float64(key string) float64 {
	f, ok := Float64(key)
	if !ok {
		panic(fmt.Sprintf("failed to retrieve '%s' float64 from config", key))
	}
	return f
}

// Val returns the interface{} value, within the root, and panics when not found.
func (r required) Val(key string) interface{} {
	o, ok := Val(key)
	if !ok {
		panic(fmt.Sprintf("failed to retrieve '%s' value from config", key))
	}
	return o
}

// GroupBool returns the boolean, within the group, and panics when not found.
func (r required) GroupBool(group, key string) bool {
	b, ok := GroupBool(group, key)
	if !ok {
		panic(fmt.Sprintf("failed to retrieve '%s' group bool from config", key))
	}
	return b
}

// GroupString returns the string, within the group, and panics when not found.
func (r required) GroupString(group, key string) string {
	s, ok := GroupString(group, key)
	if !ok {
		panic(fmt.Sprintf("failed to retrieve '%s' group string from config", key))
	}
	return s
}

// GroupInt returns the int, within the group, and panics when not found.
func (r required) GroupInt(group, key string) int {
	i, ok := GroupInt(group, key)
	if !ok {
		panic(fmt.Sprintf("failed to retrieve '%s' group int from config", key))
	}
	return i
}

// GroupFlaot64 returns the float64, within the group, and panics when not found.
func (r required) GroupFloat64(group, key string) float64 {
	f, ok := GroupFloat64(group, key)
	if !ok {
		panic(fmt.Sprintf("failed to retrieve '%s' group int from config", key))
	}
	return f
}

// GroupVal returns the interface{} value, within the group, and panics when not found.
func (r required) GroupVal(group, key string) interface{} {
	o, ok := GroupVal(group, key)
	if !ok {
		panic(fmt.Sprintf("failed to retrieve '%s' group value from config", key))
	}
	return o
}
