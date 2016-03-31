package cmd

import (
	"github.com/spf13/cobra"
	"fmt"
	"github.com/forgerock/frconfig/idmconfig"
)



// policyCmd represents the policy command
var idmCmd = &cobra.Command{
	Use:   "idm",
	Short: "Get config objects from OpenIDM",
	Long: `Get Config from OpenIDM`,
	Run: func(cmd *cobra.Command, args []string) {
		idm, err := idmconfig.GetOpenIDMConnection()
		if err != nil {
			msg := fmt.Sprintf("Can't create OpenIDM connection: %v", err)
			fatal(msg)
		}

		// get output format
		format,_ := cmd.Flags().GetString("output")

		s,err := idm.GetConfig(format)
		if err != nil {
			msg := fmt.Sprintf("Cant get config err=%v", err)
			fatal(msg)
		}
		//str,err := idm.ExportPolicies(o, realm)
		fmt.Println(s)
	},
}


func init() {
	getCmd.AddCommand(idmCmd)

}
