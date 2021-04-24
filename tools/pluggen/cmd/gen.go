package cmd

import (
	"fmt"
	"os"
	. "os/exec"
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var (
	PkgFlag    string
	PluginFlag string
	OutFlag    string
)

func NewGenCmd() (cmd *cobra.Command) {

	cmd = &cobra.Command{
		Use:     "gen",
		Example: "pluggen gen --out $PWD/testdata/myplugin.so --plugin plugins/testdata/myplugin --pkg myprojectname",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Printf("[tools/pluggen] %v\n", os.Args)
			gen, err := NewGenerator()
			if err != nil {
				return err
			}
			err = gen.ValidateArgs()
			if err != nil {
				return err
			}
			err = gen.Generate()
			return err
		},
	}

	// "$PWD/testdata/$GENERATE_PLUGIN_NAME.so"
	cmd.Flags().StringVarP(&OutFlag, "out", "o", "", "output filename")
	cmd.Flags().StringVarP(&PluginFlag, "plugin", "p", "", "directory of the plugin to re-generate")
	cmd.Flags().StringVar(&PkgFlag, "pkg", "", "root directory of the package to generate into")

	return cmd
}

func (g *generator) Generate() error {
	goBuildCommand := Command("go", "build", "-buildmode=plugin", "-o", g.OutputPath, g.PluginPath)
	goBuildCommand.Dir = g.PkgRoot
	err := goBuildCommand.Run()
	return err
}

func (g *generator) ValidateArgs() error {
	if OutFlag == "" || filepath.Ext(OutFlag) != ".so" {
		return fmt.Errorf("--out must be set")
	}
	if PluginFlag == "" {
		return fmt.Errorf("--plugin must be set")
	}
	if g.PkgRoot == "" {
		return fmt.Errorf("--pkg must be set")
	}

	return nil
}

type Option struct {
	Out        string
	PluginPath string
	PkgRoot    string
}

type generator struct {
	PluginPath string
	OutputPath string
	PkgRoot    string
}

func NewGenerator() (g *generator, err error) {
	if PkgFlag != "" {
		PkgFlag = resolvePkgRoot(os.ExpandEnv(OutFlag), os.ExpandEnv(PkgFlag))
	}

	pluginPath := os.ExpandEnv(PluginFlag)
	if !path.IsAbs(pluginPath) {
		pluginPath = path.Join(PkgFlag, pluginPath)
	}

	outputPath := os.ExpandEnv(strings.Replace(OutFlag, "./", "$PWD/", 1))
	if !path.IsAbs(outputPath) {
		outputPath = path.Join(PkgFlag, outputPath)
	}

	g = &generator{
		PluginPath: pluginPath,
		OutputPath: outputPath,
		PkgRoot:    PkgFlag,
	}
	fmt.Printf("%+v\n", *g)
	return g, nil
}

func resolvePkgRoot(outfile, root string) string {
	wd := outfile
	if !path.IsAbs(outfile) {
		wd, _ = os.Getwd()
	}
	for path.Base(wd) != root {
		wd = path.Dir(wd)
	}
	return wd
}

func isFsRoot(p string) bool {
	return p == path.Dir("/") || p == path.Dir(".")
}
