package main

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/processors"
)

type copyFields struct {
	config copyFieldsConfig
}

type copyFieldsConfig struct {
	Fields []sourceTarget `config:"fields"`
}

type sourceTarget struct {
	// Default info is Source field is missing in event.
	Default string `config:"default"`

	// Source.
	Source string `config:"source"`

	// Target.
	Target string `config:"target"`
}

func (r *copyFields) String() string {
	return fmt.Sprintf("copy fields config => %#v", r.config)
}

func CreateCopyFields(c *common.Config) (processors.Processor, error) {
	config := copyFieldsConfig{}
	err := c.Unpack(&config)
	if err != nil {
		return nil, err
	}

	return &copyFields{
		config: config,
	}, nil
}

func (f *copyFields) Run(event *beat.Event) (*beat.Event, error) {
	for _, field := range f.config.Fields {
		f.copyField(field.Source, field.Target, field.Default, event.Fields)
	}

	return event, nil
}

func (f *copyFields) copyField(source, target, defaultStr string, fields common.MapStr) {
	value, err := fields.GetValue(source)
	if err != nil && err != common.ErrKeyNotFound {
		return
	}

	if value == nil {
		// ErrKeyNotFound occured
		fields.Put(target, defaultStr)
	} else {
		fields.Put(target, value.(string))
	}
}
