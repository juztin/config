package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type required struct {
}

var (
	ConfigFile = "config.json"
	Required   = *new(required)
	cfg        map[string]interface{}
)

func Load() error {
	// get|read configuration from file
	p, c, err := getConfig()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to load configuration file: %s, from: %s\n%v", ConfigFile, p, err))
	}

	var j interface{}
	if err := json.Unmarshal(c, &j); err != nil {
		return errors.New(fmt.Sprintf("Failed to read configuration file: %s, from: %s\n%v", ConfigFile, p, err))
	}

	cfg = j.(map[string]interface{})
	return nil
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
	e = errors.New(fmt.Sprintf("Failed to find a configuration file: %s", ConfigFile))

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

func colVal(key string, col map[string]interface{}) (interface{}, bool) {
	if v, ok := col[key]; ok {
		return v, true
	}
	return nil, false
}

func Bool(key string) (bool, bool) {
	return colBool(key, cfg)
}

func String(key string) (string, bool) {
	return colString(key, cfg)
}

func Int(key string) (int, bool) {
	return colInt(key, cfg)
}

func Val(key string) (interface{}, bool) {
	return colVal(key, cfg)
}

func GroupBool(group, key string) (v bool, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colBool(key, col)
		}
	}
	return
}

func GroupString(group, key string) (v string, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colString(key, col)
		}
	}
	return
}

func GroupInt(group, key string) (v int, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colInt(key, col)
		}
	}
	return
}

func GroupVal(group, key string) (v interface{}, ok bool) {
	if m, exists := cfg[group]; exists {
		if col, isMap := m.(map[string]interface{}); isMap {
			v, ok = colVal(key, col)
		}
	}
	return
}

// root
func setCfgInt(key string, val *int) {
	if v, ok := Int(key); ok {
		*val = v
	}
}
func setCfgBool(key string, val *bool) {
	if v, ok := Bool(key); ok {
		*val = v
	}
}
func setCfgString(key string, val *string) {
	if v, ok := String(key); ok {
		*val = v
	}
}

// group
func setCfgGroupInt(group, key string, val *int) {
	if v, ok := GroupInt(group, key); ok {
		*val = v
	}
}
func setCfgGroupBool(group, key string, val *bool) {
	if v, ok := GroupBool(group, key); ok {
		*val = v
	}
}
func setCfgGroupString(group, key string, val *string) {
	if v, ok := GroupString(group, key); ok {
		*val = v
	}
}

// required
func (r required) Bool(key string) bool {
	b, ok := Bool(key)
	if !ok {
		panic(fmt.Sprintf("Failed to retrieve '%s' bool from config", key))
	}
	return b
}

func (r required) String(key string) string {
	s, ok := String(key)
	if !ok {
		panic(fmt.Sprintf("Failed to retrieve '%s' string from config", key))
	}
	return s
}

func (r required) Int(key string) int {
	i, ok := Int(key)
	if !ok {
		panic(fmt.Sprintf("Failed to retrieve '%s' int from config", key))
	}
	return i
}

func (r required) Val(key string) interface{} {
	o, ok := Val(key)
	if !ok {
		panic(fmt.Sprintf("Failed to retrieve '%s' value from config", key))
	}
	return o
}

func (r required) GroupBool(group, key string) bool {
	b, ok := GroupBool(group, key)
	if !ok {
		panic(fmt.Sprintf("Failed to retrieve '%s' group bool from config", key))
	}
	return b
}

func (r required) GroupString(group, key string) string {
	s, ok := GroupString(group, key)
	if !ok {
		panic(fmt.Sprintf("Failed to retrieve '%s' group string from config", key))
	}
	return s
}

func (r required) GroupInt(group, key string) int {
	i, ok := GroupInt(group, key)
	if !ok {
		panic(fmt.Sprintf("Failed to retrieve '%s' group int from config", key))
	}
	return i
}

func (r required) GroupVal(group, key string) interface{} {
	o, ok := GroupVal(group, key)
	if !ok {
		panic(fmt.Sprintf("Failed to retrieve '%s' group value from config", key))
	}
	return o
}
