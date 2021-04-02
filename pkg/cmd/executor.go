package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/alex-held/devctl/internal/config"
	"github.com/alex-held/devctl/pkg/logging"
)

type Executor struct {
	rootCmd    *cobra.Command
	runCmd     *cobra.Command
	lintersCmd *cobra.Command

	exitCode              int
	version, commit, date string

	cfg              *config.DevEnvConfig
	log               *logging.Log
	contextLoader     *lint.ContextLoader
	goenv             *goutil.Env
	//fileCache         *fsutils.FileCache
	//lineCache         *fsutils.LineCache
//	pkgCache          *pkgcache.Cache
	debugf logging.DebugFunc
	sw     *timeutils.Stopwatch

	loadGuard *load.Guard
	flock     *flock.Flock
}

func (e *Executor) executeZshCompletion(cmd *cobra.Command, args []string) error {
	err := cmd.Root().GenZshCompletion(os.Stdout)
	if err != nil {
		return errors.Wrap(err, "unable to generate zsh completions: %v")
	}
	// Add extra compdef directive to support sourcing command directly.
	// https://github.com/spf13/cobra/issues/881
	// https://github.com/spf13/cobra/pull/887
	fmt.Println("compdef _devctl devctl")
	return nil
}



func NewExecutor(version, commit, date string) *Executor {
	startedAt := time.Now()
	e := &Executor{
		cfg:       config.Default(),
		version:   version,
		commit:    commit,
		date:      date,
		DBManager: lintersdb.NewManager(nil, nil),
		debugf:    logutils.Debug("exec"),
	}

	e.debugf("Starting execution...")
	e.log = report.NewLogWrapper(logutils.NewStderrLog(""), &e.reportData)

	// to setup log level early we need to parse config from command line extra time to
	// find `-v` option
	commandLineCfg, err := e.getConfigForCommandLine()
	if err != nil && err != pflag.ErrHelp {
		e.log.Fatalf("Can't get config for command line: %s", err)
	}
	if commandLineCfg != nil {
		logutils.SetupVerboseLog(e.log, commandLineCfg.Run.IsVerbose)

		switch commandLineCfg.Output.Color {
		case "always":
			color.NoColor = false
		case "never":
			color.NoColor = true
		case "auto":
			// nothing
		default:
			e.log.Fatalf("invalid value %q for --color; must be 'always', 'auto', or 'never'", commandLineCfg.Output.Color)
		}
	}

	// init of commands must be done before config file reading because
	// init sets config with the default values of flags
	e.initRoot()
	e.initRun()
	e.initHelp()
	e.initLinters()
	e.initConfig()
	e.initCompletion()
	e.initVersion()
	e.initCache()

	// init e.cfg by values from config: flags parse will see these values
	// like the default ones. It will overwrite them only if the same option
	// is found in command-line: it's ok, command-line has higher priority.

	r := config.NewFileReader(e.cfg, commandLineCfg, e.log.Child("config_reader"))
	if err = r.Read(); err != nil {
		e.log.Fatalf("Can't read config: %s", err)
	}

	// recreate after getting config
	e.DBManager = lintersdb.NewManager(e.cfg, e.log).WithCustomLinters()

	e.cfg.LintersSettings.Gocritic.InferEnabledChecks(e.log)
	if err = e.cfg.LintersSettings.Gocritic.Validate(e.log); err != nil {
		e.log.Fatalf("Invalid gocritic settings: %s", err)
	}

	// Slice options must be explicitly set for proper merging of config and command-line options.
	fixSlicesFlags(e.runCmd.Flags())
	fixSlicesFlags(e.lintersCmd.Flags())

	e.EnabledLintersSet = lintersdb.NewEnabledSet(e.DBManager,
		lintersdb.NewValidator(e.DBManager), e.log.Child("lintersdb"), e.cfg)
	e.goenv = goutil.NewEnv(e.log.Child("goenv"))
	e.fileCache = fsutils.NewFileCache()
	e.lineCache = fsutils.NewLineCache(e.fileCache)

	e.sw = timeutils.NewStopwatch("pkgcache", e.log.Child("stopwatch"))
	e.pkgCache, err = pkgcache.NewCache(e.sw, e.log.Child("pkgcache"))
	if err != nil {
		e.log.Fatalf("Failed to build packages cache: %s", err)
	}
	e.loadGuard = load.NewGuard()
	e.contextLoader = lint.NewContextLoader(e.cfg, e.log.Child("loader"), e.goenv,
		e.lineCache, e.fileCache, e.pkgCache, e.loadGuard)
	if err = e.initHashSalt(version); err != nil {
		e.log.Fatalf("Failed to init hash salt: %s", err)
	}
	e.debugf("Initialized executor in %s", time.Since(startedAt))
	return e
}
