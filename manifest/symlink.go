package manifest

import (
    scriptish "github.com/ganbarodigital/go_scriptish"
    "os"
)

// ln -s <source> <target>
func Symlink(source string, target string) (func(p *scriptish.Pipe) (int, error)) {

	return func(p *scriptish.Pipe) (int, error) {

		// expand our input
		expandedSource := p.Env.Expand(source)
		expandedTarget := p.Env.Expand(target)
		/*
		   // debugging support
		   //        Tracef("Mkdir(%#v, 0%o)", source, mode)
		   //        Tracef("=> Mkdir(%#v, 0%o)", expFilepath, mode)
		*/
		err := os.Symlink(expandedSource, expandedTarget)

		if err != nil {
			return scriptish.StatusNotOkay, err
		}

		// all done
		return scriptish.StatusOkay, nil
	}
}
