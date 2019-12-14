package main

import (
	"fmt"
	"strings"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/processors"
)

type replaceSubFields struct {
	config replaceSubFieldsConfig
}

type replaceSubFieldsConfig struct {
	Fields []replaceSub `config:"fields"`
}

type replaceSub struct {
	// From to do replacement on. The default is message.
	Parent string `config:"parent"`

	// Old.
	Old string `config:"old"`

	// New.
	New string `config:"new"`
}

func (r *replaceSubFields) String() string {
	return fmt.Sprintf("replace sub config => %#v", r.config)
}

func CreateReplaceSubFields(c *common.Config) (processors.Processor, error) {
	logp.Info("CreateReplaceSubFields processor")
	config := replaceSubFieldsConfig{}
	err := c.Unpack(&config)
	if err != nil {
		return nil, err
	}

	return &replaceSubFields{
		config: config,
	}, nil
}

func (f *replaceSubFields) Run(event *beat.Event) (*beat.Event, error) {
	for _, field := range f.config.Fields {
		f.replaceSubField(field.Parent, field.Old, field.New, event.Fields)
	}

	return event, nil
}

func (f *replaceSubFields) replaceSubField(parent, old, new string, fields common.MapStr) {
	value, err := fields.GetValue(parent)
	if err != nil {
		return
	}

	originItems := value.(common.MapStr)
	targetItems := common.MapStr{}

	for key, value := range originItems {
		targetItems[strings.Replace(key, old, new, -1)] = value
	}

	fields.Put(parent, targetItems)
}
