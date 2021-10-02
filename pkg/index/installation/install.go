package installation

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alex-held/devctl-kit/pkg/log"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"

	"github.com/alex-held/devctl/pkg/constants"
	"github.com/alex-held/devctl/pkg/env"
	"github.com/alex-held/devctl/pkg/index/download"
	"github.com/alex-held/devctl/pkg/index/pathutil"
	"github.com/alex-held/devctl/pkg/index/spec"
)

// InstallOpts specifies options for plugin installation operation.
type InstallOpts struct {
	ArchiveFileOverride string
}

type installOperation struct {
	pluginName string
	platform   spec.Platform

	installDir string
	binDir     string
}

// Plugin lifecycle errors
var (
	ErrIsAlreadyInstalled = errors.New("can't install, the newest version is already installed")
	ErrIsNotInstalled     = errors.New("plugin is not installed")
	ErrIsAlreadyUpgraded  = errors.New("can't upgrade, the newest version is already installed")
)

// Install will download and install a plugin. The operation tries
// to not get the plugin dir in a bad state if it fails during the process.
func Install(p env.Factory, plugin spec.Plugin, indexName string, opts InstallOpts) error {
	log.Infof("Looking for installed versions")
	_, err := Load(p.Fs(), p.Paths().PluginInstallReceiptPath(plugin.Name))
	if err == nil {
		return ErrIsAlreadyInstalled
	} else if !os.IsNotExist(err) {
		return errors.Wrap(err, "failed to look up plugin receipt")
	}

	// Find available installation candidate
	candidate, ok, err := GetMatchingPlatform(plugin.Spec.Platforms)
	if err != nil {
		return errors.Wrap(err, "failed trying to find a matching platform in plugin spec")
	}
	if !ok {
		return errors.Errorf("plugin %q does not offer installation for this platform", plugin.Name)
	}

	// The actual install should be the last action so that a failure during receipt
	// saving does not result in an installed plugin without receipt.
	log.Infof("Install plugin %s at version=%s", plugin.Name, plugin.Spec.Version)
	if err := install(p.Fs(), installOperation{
		pluginName: plugin.Name,
		platform:   candidate,

		binDir:     p.Paths().BinPath(),
		installDir: p.Paths().PluginVersionInstallPath(plugin.Name, plugin.Spec.Version),
	}, opts); err != nil {
		return errors.Wrap(err, "install failed")
	}

	log.Infof("Storing install receipt for plugin %s", plugin.Name)
	err = Store(p.Fs(), New(plugin, indexName, metav1.Now()), p.Paths().PluginInstallReceiptPath(plugin.Name))
	return errors.Wrap(err, "installation receipt could not be stored, uninstall may fail")
}

func install(fs afero.Fs, op installOperation, opts InstallOpts) error {
	// Download and extract
	log.Infof("Creating download staging directory")
	downloadStagingDir, err := ioutil.TempDir("", "krew-downloads")
	if err != nil {
		return errors.Wrapf(err, "could not create staging dir %q", downloadStagingDir)
	}
	log.Infof("Successfully created download staging directory %q", downloadStagingDir)
	defer func() {
		log.Infof("Deleting the download staging directory %s", downloadStagingDir)
		if err := os.RemoveAll(downloadStagingDir); err != nil {
			klog.Warningf("failed to clean up download staging directory: %s", err)
		}
	}()
	if err := downloadAndExtract(fs, downloadStagingDir, op.platform.URI, op.platform.Sha256, opts.ArchiveFileOverride); err != nil {
		return errors.Wrap(err, "failed to unpack into staging dir")
	}

	applyDefaults(&op.platform)
	if err := moveToInstallDir(fs, downloadStagingDir, op.installDir, op.platform.Files); err != nil {
		return errors.Wrap(err, "failed while moving files to the installation directory")
	}

	subPathAbs, err := filepath.Abs(op.installDir)
	if err != nil {
		return errors.Wrapf(err, "failed to get the absolute fullPath of %q", op.installDir)
	}
	fullPath := filepath.Join(op.installDir, filepath.FromSlash(op.platform.Bin))
	pathAbs, err := filepath.Abs(fullPath)
	if err != nil {
		return errors.Wrapf(err, "failed to get the absolute fullPath of %q", fullPath)
	}
	if _, ok := pathutil.IsSubPath(subPathAbs, pathAbs); !ok {
		return errors.Wrapf(err, "the fullPath %q does not extend the sub-fullPath %q", fullPath, op.installDir)
	}
	err = createOrUpdateLink(fs, op.binDir, fullPath, op.pluginName)
	return errors.Wrap(err, "failed to link installed plugin")
}

func applyDefaults(platform *spec.Platform) {
	if platform.Files == nil {
		platform.Files = []spec.FileOperation{{From: "*", To: "."}}
		log.Debugf("file operation not specified, assuming %v", platform.Files)
	}
}

// downloadAndExtract downloads the specified archive uri (or uses the provided overrideFile, if a non-empty value)
// while validating its checksum with the provided sha256sum, and extracts its contents to extractDir that must be.
// created.
func downloadAndExtract(fs afero.Fs, extractDir, uri, sha256sum, overrideFile string) error {
	var fetcher download.Fetcher = download.HTTPFetcher{}
	if overrideFile != "" {
		fetcher = download.NewFileFetcher(fs, overrideFile)
	}

	verifier := download.NewSha256Verifier(sha256sum)
	err := download.NewDownloader(verifier, fetcher).Get(uri, extractDir)
	return errors.Wrap(err, "failed to unpack the plugin archive")
}

// Uninstall will uninstall a plugin.
func Uninstall(p env.Factory, name string) error {
	if name == constants.DevctlPluginName {
		log.Errorf("Removing krew through krew is not supported.")
		if !IsWindows() { // assume POSIX-like
			log.Errorf("If you’d like to uninstall krew altogether, run:\n\trm -rf -- %q", p.Paths().Base())
		}
		return errors.New("self-uninstall not allowed")
	}
	log.Infof("Finding installed version to delete")

	if _, err := Load(p.Fs(), p.Paths().PluginInstallReceiptPath(name)); err != nil {
		if os.IsNotExist(err) {
			return ErrIsNotInstalled
		}
		return errors.Wrapf(err, "failed to look up install receipt for plugin %q", name)
	}

	log.Infof("Deleting plugin %s", name)

	symlinkPath := filepath.Join(p.Paths().BinPath(), pluginNameToBin(name, IsWindows()))
	log.Infof("Unlink %q", symlinkPath)
	if err := removeLink(symlinkPath); err != nil {
		return errors.Wrap(err, "could not uninstall symlink of plugin")
	}

	pluginInstallPath := p.Paths().PluginInstallPath(name)
	log.Infof("Deleting path %q", pluginInstallPath)
	if err := os.RemoveAll(pluginInstallPath); err != nil {
		return errors.Wrapf(err, "could not remove plugin directory %q", pluginInstallPath)
	}
	pluginReceiptPath := p.Paths().PluginInstallReceiptPath(name)
	log.Infof("Deleting plugin receipt %q", pluginReceiptPath)
	err := os.Remove(pluginReceiptPath)
	return errors.Wrapf(err, "could not remove plugin receipt %q", pluginReceiptPath)
}

func createOrUpdateLink(fs afero.Fs, binDir, binary, plugin string) error {
	dst := filepath.Join(binDir, pluginNameToBin(plugin, IsWindows()))

	// TODO: afero.Fs
	if err := removeLink(dst); err != nil {
		return errors.Wrap(err, "failed to remove old symlink")
	}
	if _, err := fs.Stat(binary); os.IsNotExist(err) {
		return errors.Wrapf(err, "can't create symbolic link, source binary (%q) cannot be found in extracted archive", binary)
	}

	// Create new
	log.Infof("Creating symlink to %q at %q", binary, dst)

	// TODO: afero.Fs
	if err := os.Symlink(binary, dst); err != nil {
		return errors.Wrapf(err, "failed to create a symlink from %q to %q", binary, dst)
	}
	log.Infof("Created symlink at %q", dst)

	return nil
}

// removeLink removes a symlink reference if exists.
func removeLink(path string) error {
	// TODO: afero.Fs
	fi, err := os.Lstat(path)
	if os.IsNotExist(err) {
		log.Infof("No file found at %q", path)
		return nil
	} else if err != nil {
		return errors.Wrapf(err, "failed to read the symlink in %q", path)
	}

	if fi.Mode()&os.ModeSymlink == 0 {
		return errors.Errorf("file %q is not a symlink (mode=%s)", path, fi.Mode())
	}
	// TODO: afero.Fs
	if err := os.Remove(path); err != nil {
		return errors.Wrapf(err, "failed to remove the symlink in %q", path)
	}
	log.Infof("Removed symlink from %q", path)
	return nil
}

// IsWindows sees if KREW_OS or runtime.GOOS to find out if current execution mode is win32.
func IsWindows() bool {
	goos := runtime.GOOS
	if env := os.Getenv("KREW_OS"); env != "" {
		goos = env
	}
	return goos == "windows"
}

// pluginNameToBin creates the name of the symlink file for the plugin name.
// It converts dashes to underscores.
func pluginNameToBin(name string, isWindows bool) string {
	name = strings.ReplaceAll(name, "-", "_")
	name = "devctl-" + name
	if isWindows {
		name += ".exe"
	}
	return name
}

// CleanupStaleKrewInstallations removes the versions that aren't the current version.
func CleanupStaleKrewInstallations(fs afero.Fs, dir, currentVersion string) error {
	ls, err := afero.ReadDir(fs, dir)
	if err != nil {
		return errors.Wrap(err, "failed to read krew store directory")
	}
	log.Infof("Found %d entries in krew store directory", len(ls))
	for _, d := range ls {
		log.Infof("Found a krew installation: %s (%s)", d.Name(), d.Mode())
		if d.IsDir() && d.Name() != currentVersion {
			log.Debugf("Deleting stale krew install directory: %s", d.Name())
			p := filepath.Join(dir, d.Name())
			if err := fs.RemoveAll(p); err != nil {
				return errors.Wrapf(err, "failed to remove stale krew version at path '%s'", p)
			}
			log.Debugf("Stale installation directory removed")
		}
	}
	return nil
}
