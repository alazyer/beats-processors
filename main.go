package main

import (
	"github.com/elastic/beats/libbeat/plugin"
	"github.com/elastic/beats/libbeat/processors"
)

var Bundle = plugin.Bundle(
	processors.Plugin("add_fields", CreateAddFields),
	processors.Plugin("copy_fields", CreateCopyFields),
	processors.Plugin("flatten_fields", CreateFlattenFields),
	processors.Plugin("replace_sub_fields", CreateReplaceSubFields),
)
