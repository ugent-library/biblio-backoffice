package cmd

import (
	"google.golang.org/protobuf/encoding/protojson"
)

var (
	marshaller   = protojson.MarshalOptions{UseProtoNames: true}
	unmarshaller = protojson.UnmarshalOptions{}
)
