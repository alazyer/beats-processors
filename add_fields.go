package main

import (
	"encoding/json"
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/processors"
)

type addFields struct {
	fields common.MapStr
}

// FieldsKey is the default target key for the add_fields processor.
const FieldsKey = "fields"

// CreateAddFields constructs an add_fields processor from config.
func CreateAddFields(c *common.Config) (processors.Processor, error) {
	config := struct {
		Fields common.MapStr `config:"fields" validate:"required"`
		Target *string       `config:"target"`
	}{}
	err := c.Unpack(&config)
	if err != nil {
		return nil, fmt.Errorf("fail to unpack the add_fields configuration: %s", err)
	}

	return makeFieldsProcessor(
		optTarget(config.Target, FieldsKey),
		config.Fields,
	), nil
}

// NewAddFields creates a new processor adding the given fields to events.
// Set `shared` true if there is the chance of labels being changed/modified by
// subsequent processors.
func NewAddFields(fields common.MapStr) processors.Processor {
	return &addFields{fields: fields}
}

func (af *addFields) Run(event *beat.Event) (*beat.Event, error) {
	fields := af.fields

	event.Fields.DeepUpdate(fields)
	return event, nil
}

func (af *addFields) String() string {
	s, _ := json.Marshal(af.fields)
	return fmt.Sprintf("add_fields=%s", s)
}

func optTarget(opt *string, def string) string {
	if opt == nil {
		return def
	}
	return *opt
}

func makeFieldsProcessor(target string, fields common.MapStr) processors.Processor {
	if target != "" {
		fields = common.MapStr{
			target: fields,
		}
	}

	return NewAddFields(fields)
}
