// Cross-platform configuration manager
package config

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
)

type Config struct {
	Base    string
	Entries map[string]ConfigEntry
}

type ConfigEntry struct {
	Key   string
	Value string
}

// Create a new Config manager
func NewConfig(base string) (config *Config) {
	config = &Config{}
	config.Base = base
	config.Entries = make(map[string]ConfigEntry)
	return
}

// List keys for a given config
func (c *Config) List(name string) (keys []string, err error) {
	configDir, err := c.configDirectory()
	if err != nil {
		return
	}
	dir := path.Join(configDir, name)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for _, file := range files {
		keys = append(keys, file.Name())
	}
	sort.Strings(keys)
	return
}

// Save a key/value pair for a config
func (c *Config) Save(name, key, value string) (err error) {
	return c.SaveGlobal(name, key, value)
}

func (c *Config) SaveGlobal(name, key, value string) (err error) {
	dir, err := c.configDirectory()
	if err != nil {
		return
	}
	filename := path.Join(dir, name, key)
	err = c.writeFile(filename, value)
	return
}

func (c *Config) SaveLocal(name, key, value string) (err error) {
	dir, err := c.localConfigDirectory()
	if err != nil {
		return
	}
	filename := path.Join(dir, name, key)
	err = c.writeFile(filename, value)
	return
}

// Load a value for a config key
func (c *Config) Load(name, key string) (body string, err error) {
	return c.LoadGlobal(name, key)
}

func (c *Config) LoadGlobal(name, key string) (body string, err error) {
	dir, err := c.configDirectory()
	if err != nil {
		return
	}
	filename := path.Join(dir, name, key)
	body, err = c.readFile(filename)
	return
}

// Load value from local config directory if it exists; otherwise from global
// config directory
func (c *Config) LoadLocalOrGlobal(name, key string) (body string, err error) {
	dir, err := c.localConfigDirectory()
	if err != nil {
		return
	}
	filename := path.Join(dir, name, key)
	body, err = c.readFile(filename)
	if err != nil {
		return c.LoadGlobal(name, key)
	}
	return
}

// Delete a config key/value pair
func (c *Config) Delete(name, key string) (err error) {
	return c.DeleteGlobal(name, key)
}

func (c *Config) DeleteGlobal(name, key string) (err error) {
	dir, err := c.configDirectory()
	if err != nil {
		return
	}
	filename := path.Join(dir, name, key)
	err = os.Remove(filename)
	return
}

// Delete a config key/value pair
func (c *Config) DeleteLocalOrGlobal(name, key string) (err error) {
	dir, err := c.localConfigDirectory()
	if err != nil {
		return
	}
	filename := path.Join(dir, name, key)
	err = os.Remove(filename)
	if err != nil {
		return c.DeleteGlobal(name, key)
	}
	return
}

func (c *Config) configDirectory() (configDir string, err error) {
	home, err := c.homeDirectory()
	if err != nil {
		return
	}
	configDir = path.Join(home, fmt.Sprintf(".%s", c.Base))
	return
}

func (c *Config) localConfigDirectory() (configDir string, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return
	}
	configDir = path.Join(cwd, fmt.Sprintf(".%s", c.Base))
	return
}

func (c *Config) writeFile(filename, body string) (err error) {
	dir := filepath.Dir(filename)
	if err = os.MkdirAll(dir, 0700); err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, []byte(body), 0600)
	if err != nil {
		return
	}
	return
}

func (c *Config) readFile(filename string) (body string, err error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	body = string(data)
	return
}

func (c *Config) homeDirectory() (string, error) {
	return homedir.Dir()
}
