package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path"
	"reflect"

	"gopkg.in/yaml.v2"
)

func main() {
	app := path.Base(os.Args[0])
	if len(os.Args) < 2 {
		fmt.Printf("*** Usage: %s <command> [args]\n", app)
		return
	}
	var err error
	switch cmd := os.Args[1]; cmd {
	case "enc":
		err = encode(app, cmd, os.Args[2:])
	default:
		fmt.Println("*** Invalid command:", cmd)
	}
	if err != nil {
		fmt.Println("***", err)
		os.Exit(1)
	}
}

func encode(app, cmd string, args []string) error {
	cmdline := flag.NewFlagSet(app, flag.ExitOnError)
	fVerbose := cmdline.Bool("verbose", false, "Enable verbose debugging mode.")
	cmdline.Parse(args)

	for _, p := range cmdline.Args() {
		ext := path.Ext(p)
		dst := p[:len(p)-len(ext)] + ".enc" + ext

		if *fVerbose {
			fmt.Println("enc:", p, "->", dst)
		}

		f, err := os.Open(p)
		if err != nil {
			return err
		}

		var d interface{}
		err = yaml.NewDecoder(f).Decode(&d)
		if err != nil {
			return err
		}

		e, err := encfile(d)
		if err != nil {
			return err
		}

		r, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}

		defer r.Close()

		err = yaml.NewEncoder(r).Encode(e)
		if err != nil {
			return err
		}
	}

	return nil
}

func encfile(src interface{}) (interface{}, error) {
	return encvalue(reflect.ValueOf(src))
}

func encvalue(src reflect.Value) (interface{}, error) {
	src = reflect.Indirect(src)
	switch k := src.Kind(); k {
	case reflect.Interface:
		return encvalue(reflect.ValueOf(src.Interface()))
	case reflect.Bool:
		return base64enc(fmt.Sprint(src.Bool())), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return base64enc(fmt.Sprint(src.Int())), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return base64enc(fmt.Sprint(src.Uint())), nil
	case reflect.Float32, reflect.Float64:
		return base64enc(fmt.Sprint(src.Float())), nil
	case reflect.String:
		return base64enc(src.String()), nil
	case reflect.Array, reflect.Slice:
		v := make([]interface{}, src.Len())
		for i := 0; i < src.Len(); i++ {
			e, err := encvalue(src.Index(i))
			if err != nil {
				return nil, err
			}
			v[i] = e
		}
		return v, nil
	case reflect.Map:
		v := make(map[interface{}]interface{})
		r := src.MapRange()
		for r.Next() {
			k := reflect.Indirect(r.Key())
			e, err := encvalue(r.Value())
			if err != nil {
				return nil, err
			}
			v[k.Interface()] = e
		}
		return v, nil
	default:
		return nil, fmt.Errorf("Invalid type: (%v) %v", k, src)
	}
}

func base64enc(s string) string {
	return base64.StdEncoding.EncodeToString([]byte(s))
}
