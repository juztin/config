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
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

//type Config map[string]interface{}
type Config struct {
	mu sync.Mutex
	m  map[string]interface{}
}

var cfg, _ = Read()

func ConfigFile() string {
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		return "config.json"
	}
	return fmt.Sprintf("config.%s.json", env)
}

func ReadFrom(b []byte) (Config, error) {
	var j interface{}
	err := json.Unmarshal(b, &j)
	if err != nil {
		return *new(Config), err
	}
	//return j.(map[string]interface{}), nil
	m := j.(map[string]interface{})
	return Config{sync.Mutex{}, m}, nil
}

func Read() (Config, error) {
	cfgFile := ConfigFile()
	// Grab the path for the the running executable.
	p := filepath.Dir(os.Args[0])
	f := filepath.Join(p, cfgFile)

	// If no config file was found, look within CWD.
	_, err := os.Stat(f)
	if err != nil {
		p, err = os.Getwd()
		f = filepath.Join(p, cfgFile)
		_, err = os.Stat(f)
	}

	var c Config
	if err != nil {
		return c, err
	}
	// Read the file bytes.
	data, err := ioutil.ReadFile(f)
	if err != nil {
		return c, err
	}
	// Load the configuration from the file.
	c, err = ReadFrom(data)
	if err != nil {
		err = fmt.Errorf("failed to read configuration file %s", f)
	}
	return c, err
}

func SetConfig(m map[string]interface{}) {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	cfg.m = m
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

func keys(m map[string]interface{}) []string {
	var keys []string
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func (c Config) Keys() []string {
	return keys(cfg.m)
}

func (c Config) GroupKeys(group string) []string {
	if m, exists := c.m[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			return keys(col)
		}
	}
	return nil
}

// Bool returns the boolean value for the `key` within the root level.
// The value, or default value, is returned along with boolean of wether the key was found.
func (c Config) Bool(key string) (bool, bool) {
	return colBool(key, c.m)
}

// String returns the string value for the `key` within the root level.
// The value, or default value, is returned along with boolean of wether the key was found.
func (c Config) String(key string) (string, bool) {
	return colString(key, c.m)
}

// Int returns the int value for the `key` within the root level.
// The value, or default value, is returned along with boolean of wether the key was found.
func (c Config) Int(key string) (int, bool) {
	return colInt(key, c.m)
}

// Float64 returns the float64 value for the `key` within the root level.
// The value, or default value, is returned along with boolean of wether the key was found.
func (c Config) Float64(key string) (float64, bool) {
	return colFloat64(key, c.m)
}

// Val returns the value, as an interface{}, for the `key` within the root level.
// The value, or nil, is returned along with boolean of wether the key was found.
func (c Config) Val(key string) (interface{}, bool) {
	return colVal(key, c.m)
}

// GroupBool returns the boolean value for the `key` within the group level.
// The boolean, or false, is returned along with boolean of wether the key was found.
func (c Config) GroupBool(group, key string) (v bool, ok bool) {
	if m, exists := c.m[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colBool(key, col)
		}
	}
	return
}

// GroupBool returns the boolean value for the `key` within the group level.
// The string, or empty string, is returned along with boolean of wether the key was found.
func (c Config) GroupString(group, key string) (v string, ok bool) {
	if m, exists := c.m[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colString(key, col)
		}
	}
	return
}

// GroupBool returns the boolean value for the `key` within the group level
// The int, or 0, is returned along with boolean of wether the key was found.
func (c Config) GroupInt(group, key string) (v int, ok bool) {
	if m, exists := c.m[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colInt(key, col)
		}
	}
	return
}

// GroupBool returns the boolean value for the `key` within the group level
// The float64, or 0, is returned along with boolean of wether the key was found.
func (c Config) GroupFloat64(group, key string) (v float64, ok bool) {
	if m, exists := c.m[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colFloat64(key, col)
		}
	}
	return
}

// GroupVal returns the value, as an interface{}, for the `key` within the group level
// The value, or nil, is returned along with boolean of wether the key was found.
func (c Config) GroupVal(group, key string) (v interface{}, ok bool) {
	if m, exists := c.m[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colVal(key, col)
		}
	}
	return
}

// Bool returns the boolean value, within the root, and exits when not found.
func (c Config) RequiredBool(key string) bool {
	b, ok := c.Bool(key)
	if !ok {
		log.Fatalf("failed to retrieve '%s' bool from config", key)
	}
	return b
}

// String returns the string, within the root, and exits when not found.
func (c Config) RequiredString(key string) string {
	s, ok := c.String(key)
	if !ok {
		log.Fatalf("failed to retrieve '%s' string from config", key)
	}
	return s
}

// Int returns the int, within the root, and exits when not found.
func (c Config) RequiredInt(key string) int {
	i, ok := c.Int(key)
	if !ok {
		log.Fatalf("failed to retrieve '%s' int from config", key)
	}
	return i
}

// Float64 returns the float64, within the root, and exits when not found.
func (c Config) RequiredFloat64(key string) float64 {
	f, ok := c.Float64(key)
	if !ok {
		log.Fatalf("failed to retrieve '%s' float64 from config", key)
	}
	return f
}

// Val returns the interface{} value, within the root, and exits when not found.
func (c Config) RequiredVal(key string) interface{} {
	o, ok := c.Val(key)
	if !ok {
		log.Fatalf("failed to retrieve '%s' value from config", key)
	}
	return o
}

// GroupBool returns the boolean, within the group, and exits when not found.
func (c Config) RequiredGroupBool(group, key string) bool {
	b, ok := c.GroupBool(group, key)
	if !ok {
		log.Fatalf("failed to retrieve '%s'.'%s' group bool from config", group, key)
	}
	return b
}

// GroupString returns the string, within the group, and exits when not found.
func (c Config) RequiredGroupString(group, key string) string {
	s, ok := c.GroupString(group, key)
	if !ok {
		log.Fatalf("failed to retrieve '%s'.'%s' group string from config", group, key)
	}
	return s
}

// GroupInt returns the int, within the group, and exits when not found.
func (c Config) RequiredGroupInt(group, key string) int {
	i, ok := c.GroupInt(group, key)
	if !ok {
		log.Fatalf("failed to retrieve '%s'.'%s' group int from config", group, key)
	}
	return i
}

// GroupFlaot64 returns the float64, within the group, and exits when not found.
func (c Config) RequiredGroupFloat64(group, key string) float64 {
	f, ok := c.GroupFloat64(group, key)
	if !ok {
		log.Fatalf("failed to retrieve '%s'.'%s' group int from config", group, key)
	}
	return f
}

// GroupVal returns the interface{} value, within the group, and exits when not found.
func (c Config) RequiredGroupVal(group, key string) interface{} {
	o, ok := c.GroupVal(group, key)
	if !ok {
		log.Fatalf("failed to retrieve '%s'.'%s' group value from config", group, key)
	}
	return o
}

func Keys() []string {
	return cfg.Keys()
}

func GroupKeys(group string) []string {
	return cfg.GroupKeys(group)
}

// Bool returns the boolean value for the `key` within the root level.
// The value, or default value, is returned along with boolean of wether the key was found.
func Bool(key string) (bool, bool) {
	return cfg.Bool(key)
}

// String returns the string value for the `key` within the root level.
// The value, or default value, is returned along with boolean of wether the key was found.
func String(key string) (string, bool) {
	return cfg.String(key)
}

// Int returns the int value for the `key` within the root level.
// The value, or default value, is returned along with boolean of wether the key was found.
func Int(key string) (int, bool) {
	return cfg.Int(key)
}

// Float64 returns the float64 value for the `key` within the root level.
// The value, or default value, is returned along with boolean of wether the key was found.
func Float64(key string) (float64, bool) {
	return cfg.Float64(key)
}

// Val returns the value, as an interface{}, for the `key` within the root level.
// The value, or nil, is returned along with boolean of wether the key was found.
func Val(key string) (interface{}, bool) {
	return cfg.Val(key)
}

// GroupBool returns the boolean value for the `key` within the group level.
// The boolean, or false, is returned along with boolean of wether the key was found.
func GroupBool(group, key string) (v bool, ok bool) {
	return cfg.GroupBool(group, key)
}

// GroupBool returns the boolean value for the `key` within the group level.
// The string, or empty string, is returned along with boolean of wether the key was found.
func GroupString(group, key string) (v string, ok bool) {
	return cfg.GroupString(group, key)
}

// GroupBool returns the boolean value for the `key` within the group level
// The int, or 0, is returned along with boolean of wether the key was found.
func GroupInt(group, key string) (v int, ok bool) {
	return cfg.GroupInt(group, key)
}

// GroupBool returns the boolean value for the `key` within the group level
// The float64, or 0, is returned along with boolean of wether the key was found.
func GroupFloat64(group, key string) (v float64, ok bool) {
	return cfg.GroupFloat64(group, key)
}

// GroupVal returns the value, as an interface{}, for the `key` within the group level
// The value, or nil, is returned along with boolean of wether the key was found.
func GroupVal(group, key string) (v interface{}, ok bool) {
	return cfg.GroupVal(group, key)
}

// Bool returns the boolean value, within the root, and exits when not found.
func RequiredBool(key string) bool {
	return cfg.RequiredBool(key)
}

// String returns the string, within the root, and exits when not found.
func RequiredString(key string) string {
	return cfg.RequiredString(key)
}

// Int returns the int, within the root, and exits when not found.
func RequiredInt(key string) int {
	return cfg.RequiredInt(key)
}

// Float64 returns the float64, within the root, and exits when not found.
func RequiredFloat64(key string) float64 {
	return cfg.RequiredFloat64(key)
}

// Val returns the interface{} value, within the root, and exits when not found.
func RequiredVal(key string) interface{} {
	return cfg.RequiredVal(key)
}

// GroupBool returns the boolean, within the group, and exits when not found.
func RequiredGroupBool(group, key string) bool {
	return cfg.RequiredGroupBool(group, key)
}

// GroupString returns the string, within the group, and exits when not found.
func RequiredGroupString(group, key string) string {
	return cfg.RequiredGroupString(group, key)
}

// GroupInt returns the int, within the group, and exits when not found.
func RequiredGroupInt(group, key string) int {
	return cfg.RequiredGroupInt(group, key)
}

// GroupFlaot64 returns the float64, within the group, and exits when not found.
func RequiredGroupFloat64(group, key string) float64 {
	return cfg.RequiredGroupFloat64(group, key)
}

// GroupVal returns the interface{} value, within the group, and exits when not found.
func RequiredGroupVal(group, key string) interface{} {
	return cfg.RequiredGroupVal(group, key)
}
