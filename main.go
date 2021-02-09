package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

var Options = struct {
	IncludeHeader bool // default to false
	Pretty        bool
}{}

type JWT struct {
	Header  map[string]interface{} `json:"header"`
	Payload map[string]interface{} `json:"payload"`
}

func decode(input []byte) (*JWT, error) {
	parts := bytes.Split(input, []byte("."))
	if len(parts) != 3 {
		return nil, errors.New("invalid number of parts")
	}

	var jwt JWT

	header, err := base64.RawURLEncoding.DecodeString(string(parts[0]))
	if err != nil {
		return nil, errors.Wrap(err, "header was not base64 encoded")
	}

	if err := json.Unmarshal(header, &jwt.Header); err != nil {
		return nil, errors.Wrap(err, "header was not JSON encoded")
	}

	payload, err := base64.RawURLEncoding.DecodeString(string(parts[1]))
	if err != nil {
		return nil, errors.Wrap(err, "payload was not base64 encoded")
	}

	if err := json.Unmarshal(payload, &jwt.Payload); err != nil {
		return nil, errors.Wrap(err, "payload was not JSON encoded")
	}

	return &jwt, nil
}

func run() error {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return errors.Wrap(err, "could not read from stdin")
	}

	jwt, err := decode(input)
	if err != nil {
		return errors.Wrap(err, "could not decode jwt")
	}

	encoder := json.NewEncoder(os.Stdout)

	if Options.Pretty {
		encoder.SetIndent("", "  ")
	}

	if Options.IncludeHeader {
		if err := encoder.Encode(jwt); err != nil {
			return errors.Wrap(err, "could not encode jwt to json")
		}
	} else {
		if err := encoder.Encode(jwt.Payload); err != nil {
			return errors.Wrap(err, "could not encode jwt to json")
		}
	}

	return nil
}

func main() {
	flag.BoolVar(&Options.IncludeHeader, "header", false, "include header in output")
	flag.BoolVar(&Options.Pretty, "pretty", false, "prettify output")
	flag.Parse()

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "jwt: %s\n", err)
		os.Exit(1)
	}
}
