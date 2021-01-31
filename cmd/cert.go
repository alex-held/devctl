package cmd

import (
	"os"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"

	"github.com/alex-held/dev-env/pkg/cli"
)

// NewCertCommand represents the cert command
func NewCertCommand() (cmd *cobra.Command) {

	cmd = &cobra.Command{
		Use:   "cert",
		Short: "Manage Certificates",
		RunE: func(c *cobra.Command, args []string) error {
			return c.Usage()
		},
	}

	cmd.AddCommand(NewCACommand())
	return cmd
}

func NewCACommand() (caCommand *cobra.Command) {

	caCommand = &cobra.Command{
		Use:              "ca",
		Short:            "Manages local Certificate Authorities",
		Example:          "dev-env cert ca install --key=ca.key --cert=ca.cer",
		Args:             cobra.RangeArgs(0, 11),
		TraverseChildren: true,
	}

	keyFlag := caCommand.PersistentFlags().String("key", "", "the CA's private key ")
	certFlag := caCommand.PersistentFlags().String("cert", "", "the CA's private key ")
	_, _ = keyFlag, certFlag

	caCommand.AddCommand(
		NewCAInstallCommand(),
	)
	return caCommand
}

func NewCAInstallCommand() *cobra.Command {

	return &cobra.Command{
		Use:         "install",
		Short:       "Manages local Certificate Authorities",
		Example:     "dev-env cert ca install --key=ca.key --cert=ca.cer",
		Args:        cobra.ExactArgs(0),
		Annotations: nil,
		Run: func(c *cobra.Command, args []string) {

			const caPath = "/var/ca"
			key := c.Flag("key").Value.String()
			cert := c.Flag("cert").Value.String()

			osFs := afero.NewOsFs()
			keyBytes, err := afero.ReadFile(osFs, key)
			err = afero.WriteFile(osFs, filepath.Join(caPath, filepath.Base(key)), keyBytes, os.ModePerm)
			if err != nil {
				cli.ExitWithError(1, err)
				return
			}

			certBytes, err := afero.ReadFile(osFs, cert)
			err = afero.WriteFile(osFs, filepath.Join(caPath, filepath.Base(cert)), certBytes, os.ModePerm)
			if err != nil {
				cli.ExitWithError(1, err)
				return
			}

		},
		TraverseChildren: true,
	}
}

func init() {
	rootCmd.AddCommand(NewCertCommand())
}
