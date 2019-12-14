package main

import (
	"fmt"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/libbeat/processors"
)

type flattenFields struct {
	Fields []string
}

func (r *flattenFields) String() string {
	return fmt.Sprintf("flatten fields => %#v", r.Fields)
}

func CreateFlattenFields(c *common.Config) (processors.Processor, error) {
	logp.Info("CreateFlattenFields processor")
	config := struct {
		Fields []string `config:"fields"`
	}{}
	err := c.Unpack(&config)
	if err != nil {
		return nil, fmt.Errorf("fail to unpack the flatten_fields configuration: %s", err)
	}

	f := &flattenFields{Fields: config.Fields}
	return f, nil
}

func (f *flattenFields) Run(event *beat.Event) (*beat.Event, error) {
	for _, field := range f.Fields {
		source, err := event.Fields.GetValue(field)
		if err != nil && err != common.ErrKeyNotFound {
			// error is ignored
			logp.Err("FlattenFields processor flatten field with err: %+v", err)
			continue
		}

		switch sub := source.(type) {
		case common.MapStr:
			event.Fields.Delete(field)
			event.Fields.Put(field, sub.Flatten())
		case map[string]interface{}:
			tmp := common.MapStr(sub)
			event.Fields.Delete(field)
			event.Fields.Put(field, tmp.Flatten())
		default:
			continue
		}
	}

	return event, nil
}
