package ts3

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mitchellh/mapstructure"
)

var (
	// encoder performs white space and special character encoding
	// as required by the ServerQuery protocol.
	encoder = strings.NewReplacer(
		`\`, `\\`,
		`/`, `\/`,
		` `, `\s`,
		`|`, `\p`,
		"\a", `\a`,
		"\b", `\b`,
		"\f", `\f`,
		"\n", `\n`,
		"\r", `\r`,
		"\t", `\t`,
		"\v", `\v`,
	)

	// decoder performs white space and special character decoding
	// as required by the ServerQuery protocol.
	decoder = strings.NewReplacer(
		`\\`, "\\",
		`\/`, "/",
		`\s`, " ",
		`\p`, "|",
		`\a`, "\a",
		`\b`, "\b",
		`\f`, "\f",
		`\n`, "\n",
		`\r`, "\r",
		`\t`, "\t",
		`\v`, "\v",
	)
)

// Decode returns a decoded version of str.
func Decode(str string) string {
	return decoder.Replace(str)
}

// DecodeResponse decodes a response into a struct.
func DecodeResponse(lines []string, v interface{}) error {
	if len(lines) != 1 {
		return NewInvalidResponseError("too many lines", lines)
	}

	input := make(map[string]interface{})
	value := reflect.ValueOf(v)
	var slice reflect.Value
	var elemType reflect.Type
	if value.Kind() == reflect.Ptr {
		slice = value.Elem()
		if slice.Kind() == reflect.Slice {
			elemType = slice.Type().Elem()
		}
	}

	for _, part := range strings.Split(lines[0], "|") {
		for _, val := range strings.Split(part, " ") {
			parts := strings.SplitN(val, "=", 2)
			// TODO(steve): support groups
			key := Decode(parts[0])
			if len(parts) == 2 {
				v := Decode(parts[1])
				if i, err := strconv.Atoi(v); err != nil {
					input[key] = v
				} else {
					input[key] = i
				}
			} else {
				input[key] = ""
			}
		}

		if elemType != nil {
			// Expecting a slice
			if err := decodeSlice(elemType, slice, input); err != nil {
				return err
			}

			// Reset the input map
			input = make(map[string]interface{})
		}
	}

	if elemType != nil {
		// Expecting a slice, already decoded
		return nil
	}

	return decodeMap(input, v)
}

// decodeMap decodes input into r.
func decodeMap(d map[string]interface{}, r interface{}) error {
	cfg := &mapstructure.DecoderConfig{
		WeaklyTypedInput: true,
		TagName:          "ms",
		Result:           r,
	}
	dec, err := mapstructure.NewDecoder(cfg)
	if err != nil {
		return err
	}
	return dec.Decode(d)
}

// decodeSlice decodes input into slice.
func decodeSlice(elemType reflect.Type, slice reflect.Value, input map[string]interface{}) error {
	var v reflect.Value
	if elemType.Kind() == reflect.Ptr {
		v = reflect.New(elemType.Elem())
	} else {
		v = reflect.New(elemType)
	}

	if !v.CanInterface() {
		return fmt.Errorf("can't interface %#v", v)
	}

	if err := decodeMap(input, v.Interface()); err != nil {
		return err
	}

	if elemType.Kind() == reflect.Struct {
		v = v.Elem()
	}
	slice.Set(reflect.Append(slice, v))

	return nil
}
