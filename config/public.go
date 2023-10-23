package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/ComixHe/proxy/utils"
	"github.com/go-playground/validator"
	"go.uber.org/zap"
)

const (
	globalConfigLocation = "/etc/deepin-proxy/config.json"
)

func New(opts ...Opt) (ret *GlobalConfig, err error) {
	c := &GlobalConfig{}
	for i := range opts {
		c, err = opts[i](c)
		if err != nil {
			return
		}
	}

	if c.log == nil {
		c.log = zap.NewNop().Sugar()
	}

	for _, v := range c.Configs {
		err = v.Check()
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}

	ret = c
	return
}

func WithContent(raw []byte) Opt {
	return func(c *GlobalConfig) (ret *GlobalConfig, err error) {
		c.raw = raw
		ret = c
		return
	}
}

func WithLogger(log *zap.SugaredLogger) Opt {
	return func(c *GlobalConfig) (ret *GlobalConfig, err error) {
		c.log = log
		ret = c
		return
	}
}

func LoadConfig() (content []byte, err error) {
	log := utils.GetLogger()
	c, err := os.ReadFile(globalConfigLocation)
	conf := &GlobalConfig{}
	if errors.Is(err, os.ErrNotExist) {
		c, e := json.MarshalIndent(conf, "", "\t")
		if e != nil {
			log.Warn(e)
			return nil, e
		}

		e = os.MkdirAll(path.Dir(globalConfigLocation), os.ModePerm)
		if e != nil {
			return nil, e
		}

		file, e := os.OpenFile(globalConfigLocation, os.O_WRONLY|os.O_CREATE, os.ModePerm)
		if e != nil {
			log.Warn(e)
			return nil, e
		}
		defer file.Close()

		_, e = file.Write(c)
		if e != nil {
			log.Warn(e)
			return nil, e
		}
		content = c
		return
	} else if err != nil {
		log.Warn(err)
		return nil, err
	}

	content = c

	return
}

func (c *UserConfig) Check() (err error) {

	var validator = validator.New()
	err = validator.Struct(c)
	if err != nil {
		err = fmt.Errorf("validator: %w", err)
		return
	}

	if c.TProxies == nil {
		c.TProxies = map[string]*TProxy{}
	}

	for name := range c.TProxies {
		tp := c.TProxies[name]
		if tp.Name == "" {
			tp.Name = name
		}
		if tp.DNSHijack != nil && tp.DNSHijack.IP == nil {
			addr := IPv4LocalhostStr
			tp.DNSHijack.IP = &addr
		}
	}

	return
}

func (c *GlobalConfig) WriteToFile() (err error) {
	content, err := json.MarshalIndent(c, "", "\t")
	if err != nil {
		c.log.Warn("write failed.", "error", err)
	}

	err = os.WriteFile(globalConfigLocation, content, os.ModePerm)
	return
}
