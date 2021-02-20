package cli

import (
	logger "github.com/sirupsen/logrus"
)

// VDebugLog is a "Verbose" debug logger; enable it if you really
// want spam and/or minutiae
type VDebugLog struct {
	log *logger.Logger
	// dumpSiteLoadUser bool
	// dumpPayload      bool

	// lock sync.RWMutex
	// lev  VDebugLevel
}

type VDebugLevel int

func NewVDebugLog(l *logger.Logger) *VDebugLog {
	return &VDebugLog{log: l}
}

type GlobalContext struct {
	Log *logger.Logger // Handles all logging
	VDL *VDebugLog     // verbose debug log
	Env *Env           // Env variables, cmdline args & config

	// API                              API
	/* How to make a REST call to the server */
	// XAPI                             ExternalAPI

	/* for contacting Twitter, Github, etc. */

	// DesktopAppState                  *DesktopAppState
	/* The state of focus for the currently running instance of the app */

}

// Contextified objects have explicit references to the GlobalContext,
// so that G can be swapped out for something else.  We're going to incrementally
// start moving objects over to this system.
type Contextified struct {
	g *GlobalContext
}

func (c Contextified) G() *GlobalContext {
	return c.g
}

func (c Contextified) GStrict() *GlobalContext {
	return c.g
}

func (c *Contextified) SetGlobalContext(g *GlobalContext) { c.g = g }

func NewContextified(gc *GlobalContext) Contextified {
	return Contextified{g: gc}
}

type Contextifier interface {
	G() *GlobalContext
}
