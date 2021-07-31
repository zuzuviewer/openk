package command

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
)

var (
	// file of data source
	sourceFile string
	// output,default is empty,output to stdout
	dest string
	// is output prettified json data,if it's true,only execute prettified data and only accept json format data
	indent bool
	// is output string format json data,if it's true,only output string format json data and only accept json format data
	outStr bool
)

func init() {
	convertCmd.Flags().StringVarP(&sourceFile, "file", "f", "", "source file (required)")
	convertCmd.MarkFlagRequired("file")
	convertCmd.Flags().StringVarP(&dest, "write", "w", "", "dest to output")
	convertCmd.Flags().BoolVarP(&indent, "indent", "i", false, "output prettified json")
	convertCmd.Flags().BoolVarP(&outStr, "string", "s", false, "output string format json")
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert json to yaml or yaml to json",
	Long: "Convert distinguish file type by file name extension\n" +
		"'*.json' file will be treat as json file and convert to yaml data,'*.yml' and '*.yaml' will be treat as yaml file and convert to json data\n" +
		"if file name doesn't have these extension name,Convert will try convert json to yaml first,if error occurred,then try convert yaml to json\n" +
		"the result output to stdout default,otherwise output to the dest if execute with -w.\n" +
		"prettified json output with -i or --indent,and now only accept json data and won't convert data to yaml.\n" +
		"output string format json data with -s or --string.\n" +
		"flag 'i'(indent) and 's'(string) only accept json data.\n",
	Example: "openk convert -f test.json\nopenk convert -f test.yml -w test.json\n" +
		"openk convert -i -f json.txt\nopenk convert -i -f json.txt -w t.json\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		return convert()
	},
}

func convert() error {
	data, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
	}
	if indent || outStr {
		ok := json.Valid(data)
		if !ok {
			return fmt.Errorf("invaild json input data")
		}
		// only do prettified json data
		if indent {
			data, err = indentJson(data)
			if err != nil {
				return err
			}
			return output(data)
		}
		// only output string format json data
		buf := bytes.NewBuffer(make([]byte, 0))
		if err = json.Compact(buf, data); err != nil {
			return err
		}
		return output(buf.Bytes())
	}

	if strings.HasSuffix(sourceFile, ".json") {
		if data, err = yaml.JSONToYAML(data); err != nil {
			return err
		}
	} else if strings.HasSuffix(sourceFile, ".yml") ||
		strings.HasSuffix(sourceFile, ".yaml") {
		if data, err = yaml.YAMLToJSON(data); err != nil {
			return err
		}
		// pretty json
		if data, err = indentJson(data); err != nil {
			return err
		}
	} else {
		// avoid yaml.YAMLToJson use null data
		var yamlData []byte
		// try json to yaml,if error occured,try yaml to json
		if yamlData, err = yaml.JSONToYAML(data); err != nil {
			firstErr := err
			if data, err = yaml.YAMLToJSON(data); err != nil {
				return fmt.Errorf("try convert json to yaml error:%v\ntry convert yaml to json error:%v", firstErr, err)
			} else {
				if data, err = indentJson(data); err != nil {
					return err
				}
			}
		} else {
			data = yamlData
		}
	}
	if data == nil {
		return errors.New("input data is nil")
	}
	return output(data)
}

// indentJson pretty json data
func indentJson(data []byte) ([]byte, error) {
	var str bytes.Buffer
	if err := json.Indent(&str, data, "", "  "); err != nil {
		return nil, err
	}
	return str.Bytes(), nil
}

// output write data to stdout or dest file
func output(data []byte) error {
	var (
		out io.Writer
		err error
	)
	if len(dest) == 0 {
		out = os.Stdout
	} else {
		if out, err = os.OpenFile(dest, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0666); err != nil {
			return err
		}
	}
	_, err = out.Write(data)
	return err
}
