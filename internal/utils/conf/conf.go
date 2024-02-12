package conf

import (
	"gopkg.in/ini.v1"
)

func LoadConf(path string) error {
	cfg, err := ini.Load(path)
	if err != nil {
		return err
	}

	config := new(config)

	config.cfg = cfg
	config.path = path

	if err := cfg.MapTo(config); err != nil {
		return err
	}

	Config = config

	return nil
}

func (c *config) GetSection(section string) *ini.Section {
	return c.cfg.Section(section)
}

func (c *config) GetKey(section, key string) *ini.Key {
	return c.cfg.Section(section).Key(key)
}

func (c *config) SetKey(section, key, value string) error {
	c.cfg.Section(section).Key(key).SetValue(value)
	return c.SaveTo()
}

func (c *config) SaveTo() error {
	return c.cfg.SaveTo(c.path)
}
