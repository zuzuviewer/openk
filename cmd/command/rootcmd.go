package command

import "github.com/spf13/cobra"

func init() {
	RootOpenkCmd.AddCommand(convertCmd)
}

var RootOpenkCmd = &cobra.Command{
	Use:   "openk",
	Short: "openk is for convert json to yaml and yaml to json",
	Long:  "openk support convert json to yaml and yaml to json,input from a file,output to stdout or dest file",
}
