package golang

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/alex-held/devctl/pkg/devctlpath"
	"github.com/alex-held/devctl/pkg/system"
	"github.com/spf13/afero"
)

type Config struct {
	InstallPath string `yaml:"install_path"`
	Fs          afero.Fs
}

// CreateConfig creates the default plugin config
func CreateConfig(devctlPath string) map[string]string {
	goSDKInstallPath := devctlpath.NewPather(devctlpath.WithConfigRootFn(func() string {
		return devctlPath
	})).SDK("go")

	return map[string]string{
		goSDKInstallPathKey: goSDKInstallPath,
	}
}

const goSDKInstallPathKey = "goSDKInstallPath"

func New(cfgMap map[string]string, args []string) (err error) {
	args = args[0:]

	cfg := &Config{InstallPath: cfgMap[goSDKInstallPathKey]}

	if len(args) == 0 {
		usage()
		return fmt.Errorf("must atleast provide one argument")
	}

	subcmd := args[0]
	args = args[1:]

	switch subcmd {
	case "list":
		if err = validateArgsForSubcommand(subcmd, args, 0); err != nil {
			return err
		}
		return handleList(cfg)
	case "current":
		if err = validateArgsForSubcommand(subcmd, args, 0); err != nil {
			return err
		}
		return handleCurrent(cfg)
	case "install":
		if err = validateArgsForSubcommand(subcmd, args, 1); err != nil {
			return err
		}
		return handleInstall(args[0], cfg)
	case "use":
		if err = validateArgsForSubcommand(subcmd, args, 1); err != nil {
			return err
		}
		return handleUse(args[0], cfg)
	default:
		return fmt.Errorf("unknown subcommand '%s'; expected on of 'list, current, install, use'", subcmd)
	}
}

func formatGoArchiveArtifactName(ri system.RuntimeInfo, version string) string {
	return ri.Format("go%s.[os]-[arch].tar.gz", version)
}

func dlArchive(version string, fs afero.Fs) (archive *bytes.Buffer, err error) {
	ri := system.OSRuntimeInfoGetter{}
	artifactName := formatGoArchiveArtifactName(ri.Get(), strings.TrimPrefix(version, "v"))

	// file, err := afero.TempFile(fs, "", "devctl-plugin-golang-*")
	// if err != nil {
	// 	return nil, err
	// }

	dlUri := ri.Get().Format("%s/dl/%s", "https://golang.org", artifactName)

	buf := &bytes.Buffer{}
	err = download(context.Background(), dlUri, buf)
	if err != nil {
		return buf, fmt.Errorf("failed downloading go sdk %v from the remote server %s; err=%v", version, "https://golang.org", err)
	}

	return buf, nil
}

func download(ctx context.Context, url string, outWriter io.Writer) (err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(outWriter, resp.Body)
	return err
}

func handleInstall(version string, config *Config) (err error) {

	installPath := path.Join(config.InstallPath, version)

	archive, err := dlArchive(version, config.Fs)
	if err != nil {
		return err
	}

	err = config.Fs.MkdirAll(installPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to Extract go sdk %s; dest=%s; archive=%s;err=%v\n", version, installPath, "*Bytes.Buffer", err)
	}
	err = unTarGzip(archive, installPath, unarchiveRenamer(), config.Fs)
	if err != nil {
		return fmt.Errorf("failed to Extract go sdk %s; dest=%s; archive=%s;err=%v\n", version, installPath, "*Bytes.Buffer", err)
	}
	return nil
}

func handleList(config *Config) (err error) {
	installPath := path.Join(config.InstallPath)
	ds, err := os.ReadDir(installPath)

	if err != nil {
		return err
	}

	output := strings.Builder{}

	for _, d := range ds {
		if d.IsDir() {
			output.WriteString("v" + d.Name() + "\n")
		}
	}

	fmt.Print(output.String())
	return nil
}

func handleUse(version string, config *Config) (err error) {
	installPath := path.Join(config.InstallPath, version)
	fmt.Printf("[sdk/go] installing go sdk version %s into %s", version, installPath)

	return err
}

func handleCurrent(config *Config) (err error) {
	installPath := path.Join(config.InstallPath, "current")
	link, err := os.Readlink(installPath)
	if err != nil {
		return err
	}
	currentDir := path.Base(link)
	currentVersion := "v" + currentDir
	fmt.Println(currentVersion)
	return nil
}

func usage() {
	fmt.Printf("USAGE")
}

func validateArgsForSubcommand(subcmd string, args []string, expected int) error {
	if len(args) != expected {
		return fmt.Errorf("provided wrong number of argument for subcommand '%s'; expected=%d; provided=%d", subcmd, expected, len(args))
	}
	return nil
}

func unTarGzip(buf *bytes.Buffer, target string, renamer Renamer, fs afero.Fs) error {
	gr, _ := gzip.NewReader(buf)
	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()

		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		filename := header.Name
		if renamer != nil {
			filename = renamer(filename)
		}

		p := filepath.Join(target, filename)
		fi := header.FileInfo()

		if fi.IsDir() {
			if e := fs.MkdirAll(p, fi.Mode()); e != nil {
				return e
			}
			continue
		}
		file, err := fs.OpenFile(p, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fi.Mode())
		if err != nil {
			return err
		}

		_, err = io.Copy(file, tr)
		file.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

type Renamer func(p string) string

func unarchiveRenamer() Renamer {
	return func(p string) string {
		parts := strings.Split(p, string(filepath.Separator))
		parts = parts[1:]
		newPath := strings.Join(parts, string(filepath.Separator))
		return newPath
	}
}
