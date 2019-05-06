package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"

	yaml "gopkg.in/yaml.v2"
)

type transformer func(string) (string, error)

func main() {
	app := path.Base(os.Args[0])
	if len(os.Args) < 2 {
		fmt.Printf("*** Usage: %s <command> [args]\n", app)
		return
	}
	var err error
	switch cmd := os.Args[1]; cmd {
	case "enc":
		err = runCmd(app, cmd, os.Args[2:], base64enc, os.Stdout)
	case "dec":
		err = runCmd(app, cmd, os.Args[2:], base64dec, os.Stdout)
	default:
		fmt.Println("*** Invalid command:", cmd)
	}
	if err != nil {
		fmt.Println("***", err)
		os.Exit(1)
	}
}

func runCmd(app, cmd string, args []string, xform transformer, dst io.Writer) error {
	cmdline := flag.NewFlagSet(app, flag.ExitOnError)
	fVerbose := cmdline.Bool("verbose", false, "Enable verbose debugging mode.")
	cmdline.Parse(args)

	args = cmdline.Args()
	if len(args) < 1 {
		err := xformFile(app, cmd, args, xform, os.Stdin, dst)
		if err != nil {
			return err
		}
	} else {
		for i, p := range args {
			if *fVerbose {
				fmt.Println("enc:", p)
			}
			if i > 0 {
				fmt.Fprintln(dst, "---")
			}

			f, err := os.Open(p)
			if err != nil {
				return err
			}

			defer f.Close()

			err = xformFile(app, cmd, args, xform, f, dst)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func xformFile(app, cmd string, args []string, xform transformer, src io.Reader, dst io.Writer) error {
	var values interface{}

	err := yaml.NewDecoder(src).Decode(&values)
	if err != nil {
		return err
	}

	encoded, err := procFile(xform, values)
	if err != nil {
		return err
	}

	err = yaml.NewEncoder(dst).Encode(encoded)
	if err != nil {
		return err
	}

	return nil
}

func procFile(xform transformer, src interface{}) (interface{}, error) {
	return procValue(xform, reflect.ValueOf(src))
}

func procValue(xform transformer, src reflect.Value) (interface{}, error) {
	src = reflect.Indirect(src)
	switch k := src.Kind(); k {
	case reflect.Interface:
		return procValue(xform, reflect.ValueOf(src.Interface()))
	case reflect.Bool:
		return xform(fmt.Sprint(src.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return xform(fmt.Sprint(src.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return xform(fmt.Sprint(src.Uint()))
	case reflect.Float32, reflect.Float64:
		return xform(fmt.Sprint(src.Float()))
	case reflect.String:
		return xform(src.String())
	case reflect.Array, reflect.Slice:
		v := make([]interface{}, src.Len())
		for i := 0; i < src.Len(); i++ {
			e, err := procValue(xform, src.Index(i))
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
			e, err := procValue(xform, r.Value())
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

func base64enc(s string) (string, error) {
	return base64.StdEncoding.EncodeToString([]byte(s)), nil
}

func base64dec(s string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
