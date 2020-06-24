package cmd

import (
	"fmt"
	
	"github.com/spf13/cobra"
	
	config2 "github.com/alex-held/dev-env/internal/config"
)

// currentCmd represents the current command
var currentCmd = &cobra.Command{
	Use:   "current [setting]",
	Short: "Display the current configuration of the setting",
	Long: `A longer description that spans multiple lines and likely contains examples
                and usage of using your command. For example:`,
	Example:   "current java",
	ValidArgs: []string{"java", "javafx", "go", "kotlin"},
	Run: func(cmd *cobra.Command, args []string) {

		config := readOrCreateConfig()

		result := config.ListMatchingSdks(func(sdk config2.SDK) bool {
			return sdk.Name == args[0]
		})

		// THIS IS WRONG
		current := &result[0]
		path := current.Path
		if current == nil {
			fmt.Printf("Could not resolve SDK Path for %s\n\n", args[0])
			return
		}

		fmt.Println(path)
	},
}

func init() {
	rootCmd.AddCommand(currentCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// currentCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// currentCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
