package parser

import (
	"bytes"
	"encoding/base64"
	"os"
	"text/template"

	"github.com/bborbe/teamvault-utils"
	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type Parser interface {
	Parse(content []byte) ([]byte, error)
}

type configParser struct {
	teamvaultConnector teamvault.Connector
}

func New(
	teamvaultConnector teamvault.Connector,
) *configParser {
	c := new(configParser)
	c.teamvaultConnector = teamvaultConnector
	return c
}

func (c *configParser) Parse(content []byte) ([]byte, error) {
	t, err := template.New("config").Funcs(c.createFuncMap()).Parse(string(content))
	if err != nil {
		glog.V(2).Infof("parse config failed: %v", err)
		return nil, err
	}
	b := &bytes.Buffer{}
	if err := t.Execute(b, nil); err != nil {
		glog.V(2).Infof("execute template failed: %v", err)
		return nil, err
	}
	return b.Bytes(), nil
}

func (c *configParser) createFuncMap() template.FuncMap {
	return template.FuncMap{
		"teamvaultUser": func(val interface{}) (interface{}, error) {
			glog.V(4).Infof("get teamvault value for %v", val)
			if val == nil {
				return "", nil
			}
			key := teamvault.Key(val.(string))
			user, err := c.teamvaultConnector.User(key)
			if err != nil {
				glog.V(2).Infof("get user from teamvault for key %v failed: %v", key, err)
				return "", err
			}
			glog.V(4).Infof("return value %s", user.String())
			return user.String(), nil
		},
		"teamvaultPassword": func(val interface{}) (interface{}, error) {
			glog.V(4).Infof("get teamvault value for %v", val)
			if val == nil {
				return "", nil
			}
			key := teamvault.Key(val.(string))
			pass, err := c.teamvaultConnector.Password(key)
			if err != nil {
				glog.V(2).Infof("get password from teamvault for key %v failed: %v", key, err)
				return "", err
			}
			glog.V(4).Infof("return value %s", pass.String())
			return pass.String(), nil
		},
		"teamvaultHtpasswd": func(val interface{}) (interface{}, error) {
			glog.V(4).Infof("get teamvault value for %v", val)
			if val == nil {
				return "", nil
			}
			htpasswd := teamvault.Htpasswd{
				Connector: c.teamvaultConnector,
			}
			content, err := htpasswd.Generate(teamvault.Key(val.(string)))
			if err != nil {
				return "", errors.Wrap(err, "generate htpasswd failed")
			}
			glog.V(4).Infof("return value %s", string(content))
			return string(content), nil
		},
		"teamvaultUrl": func(val interface{}) (interface{}, error) {
			glog.V(4).Infof("get teamvault value for %v", val)
			if val == nil {
				return "", nil
			}
			key := teamvault.Key(val.(string))
			pass, err := c.teamvaultConnector.Url(key)
			if err != nil {
				glog.V(2).Infof("get url from teamvault for key %v failed: %v", key, err)
				return "", err
			}
			glog.V(4).Infof("return value %s", pass.String())
			return pass.String(), nil
		},
		"teamvaultFile": func(val interface{}) (interface{}, error) {
			glog.V(4).Infof("get teamvault value for %v", val)
			if val == nil {
				return "", nil
			}
			key := teamvault.Key(val.(string))
			file, err := c.teamvaultConnector.File(key)
			if err != nil {
				glog.V(2).Infof("get file from teamvault for key %v failed: %v", key, err)
				return "", err
			}
			glog.V(4).Infof("return value %s", file.String())
			content, err := file.Content()
			if err != nil {
				return "", err
			}
			return string(content), nil
		},
		"teamvaultFileBase64": func(val interface{}) (interface{}, error) {
			glog.V(4).Infof("get teamvault value for %v", val)
			if val == nil {
				return "", nil
			}
			key := teamvault.Key(val.(string))
			file, err := c.teamvaultConnector.File(key)
			if err != nil {
				glog.V(2).Infof("get file from teamvault for key %v failed: %v", key, err)
				return "", err
			}
			glog.V(4).Infof("return value %s", file.String())
			content, err := file.Content()
			if err != nil {
				return "", err
			}
			return base64.StdEncoding.EncodeToString(content), nil
		},
		"env": func(val interface{}) (interface{}, error) {
			glog.V(4).Infof("get env value for %v", val)
			if val == nil {
				return "", nil
			}
			value := os.Getenv(val.(string))
			glog.V(4).Infof("return value %s", value)
			return value, nil
		},
		"base64": func(val interface{}) (interface{}, error) {
			glog.V(4).Infof("base64 value %v", val)
			if val == nil {
				return "", nil
			}
			return base64.StdEncoding.EncodeToString([]byte(val.(string))), nil
		},
	}
}
