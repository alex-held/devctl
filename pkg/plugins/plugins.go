package plugins

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/events"
	"github.com/karrick/godirwalk"
	"github.com/markbates/oncer"
	"github.com/markbates/safe"
	"github.com/sirupsen/logrus"

	"github.com/alex-held/devctl/pkg/constants"
)

type List map[string]Commands

var _list List

func Available() (available List, err error) {
	oncer.Do("plugins.Available", func() {
		defer func() {
			if err := saveCache(); err != nil {
				logrus.Error(err)
			}
		}()

		paths := []string{"plugins"}

		from, err := envy.MustGet(constants.DevctlEnvRootKey)
		if err != nil {
			logrus.Error(err)
		}
		paths = append(paths, strings.Split(from, string(os.PathListSeparator))...)

		list := List{}
		for _, p := range paths {
			if _, err := os.Stat(p); err != nil {
				continue
			}

			err := godirwalk.Walk(p, &godirwalk.Options{
				FollowSymbolicLinks: true,
				Callback: func(path string, info *godirwalk.Dirent) error {
					if err != nil {
						return nil
					}
					if info.IsDir() {
						return nil
					}
					base := filepath.Base(path)

					if hasPluginPrefix(base) {
						ctx, cancel := context.WithTimeout(context.Background(), timeout())
						commands := askBin(ctx, path)
						cancel()
						for _, c := range commands {
							devCtlCommand := c.DevCtlCommand
							if _, ok := list[devCtlCommand]; !ok {
								list[devCtlCommand] = Commands{}
							}
							c.Binary = path
							list[devCtlCommand] = append(list[devCtlCommand], c)
						}
					}
					return nil
				},
			})

			if err != nil {
				return
			}
		}

		_list = list
	})
	return _list, err
}

func askBin(ctx context.Context, path string) Commands {
	start := time.Now()
	defer func() {
		logrus.Debugf("askBin %s=%.4f s", path, time.Since(start).Seconds())
	}()

	commands := Commands{}
	if cp, ok := findInCache(path); ok {
		s := sum(path)
		if s == cp.CheckSum {
			logrus.Debugf("cache hit: %s", path)
			commands = cp.Commands
			return commands
		}
	}
	logrus.Debugf("cache miss: %s", path)
	if strings.HasPrefix(filepath.Base(path), "buffalo-no-sqlite") {
		return commands
	}

	cmd := exec.CommandContext(ctx, path, "available")
	bb := &bytes.Buffer{}
	cmd.Stdout = bb
	err := cmd.Run()
	if err != nil {
		logrus.Errorf("[PLUGIN] error loading plugin %s: %s\n", path, err)
		return commands
	}

	msg := bb.String()
	for len(msg) > 0 {
		err = json.NewDecoder(strings.NewReader(msg)).Decode(&commands)
		if err == nil {
			addToCache(path, cachedPlugin{
				Commands: commands,
			})
			return commands
		}
		msg = msg[1:]
	}
	logrus.Errorf("[PLUGIN] error decoding plugin %s: %s\n%s\n", path, err, msg)
	return commands
}

const timeoutEnv = "DEVCTL_PLUGIN_TIMEOUT"

var t = time.Second * 2

func timeout() time.Duration {
	oncer.Do("plugins.timeout", func() {
		rawTimeout, err := envy.MustGet(timeoutEnv)
		if err == nil {
			if parsed, err := time.ParseDuration(rawTimeout); err == nil {
				t = parsed
			} else {
				logrus.Errorf("%q value is malformed assuming default %q: %v", timeoutEnv, t, err)
			}
		} else {
			logrus.Debugf("%q not set, assuming default of %v", timeoutEnv, t)
		}
	})
	return t
}

func hasPluginPrefix(base string) bool {
	prefixes := []string{
		"devctl-",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(base, prefix) {
			return true
		}
	}
	return false
}

func LoadPlugins() (err error) {
	oncer.Do(constants.LoadPluginsEvent, func() {
		// don't send plugins events during testing
		if envy.Get(constants.DevctlEnvKey, "development") == "test" {
			return
		}

		plugs, err := Available()
		if err != nil {
			return
		}

		for _, commands := range plugs {
			for _, c := range commands {
				if c.DevCtlCommand != "events" {
					continue
				}

				err := func(c Command) error {
					return safe.RunE(func() error {
						n := fmt.Sprintf("[PLUGIN] %s %s", c.Binary, c.Name)
						fn := func(e events.Event) {
							b, err := json.Marshal(e)
							if err != nil {
								fmt.Println("error trying to marshal event", e, err)
								return
							}
							cmd := exec.Command(c.Binary, c.UseCommand, string(b))
							cmd.Stderr = os.Stderr
							cmd.Stdout = os.Stdout
							cmd.Stdin = os.Stdin
							if err := cmd.Run(); err != nil {
								fmt.Println("error trying to send event", strings.Join(cmd.Args, " "), err)
							}
						}

						_, err := events.NamedListen(n, events.Filter(c.ListenFor, fn))
						if err != nil {
							return err
						}
						return nil
					})
				}(c)
				if err != nil {
					err = err
					return
				}
			}
		}
	})
	return err
}
