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
)

func init() {
	convertCmd.Flags().StringVarP(&sourceFile, "file", "f", "", "source file (required)")
	convertCmd.MarkFlagRequired("file")
	convertCmd.Flags().StringVarP(&dest, "write", "w", "", "dest to output")
}

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "convert json to yaml or yaml to json",
	Long: "Convert distinguish file type by file name extension\n" +
		"'*.json' file will be treat as json file and convert to yaml data,'*.yml' and '*.yaml' will be treat as yaml file and convert to json data\n" +
		"if file name doesn't have these extension name,Convert will try convert json to yaml first,if error occurred,then try convert yaml to json\n" +
		"the result output to stdout default,otherwise output to the dest if execute with -w",
	Example: "openk convert -f test.json\nopenk convert -f test.yml -w test.json",
	RunE: func(cmd *cobra.Command, args []string) error {
		return convert()
	},
}

func convert() error {
	data, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		return err
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
