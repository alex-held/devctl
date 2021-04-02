package logging

import (
	"fmt"
	"strings"

	"github.com/muesli/termenv"
	"github.com/sirupsen/logrus"
)

type DevCtlFormatter struct{}

type lLevel logrus.Level

func (l lLevel) LevelStyle() termenv.Style {
	return l.ColorPrintf(strings.ToUpper(logrus.Level(l).String()))
}
func (l lLevel) ColorPrintf(format string, args ...interface{}) termenv.Style {
	clr := l.Color()
	text := fmt.Sprintf(format, args...)
	return termenv.
		String(text).
		Foreground(clr).
		Bold().
		Underline()
}

func (l lLevel) Color() termenv.Color {
	var clr termenv.Color
	switch logrus.Level(l) {
	case logrus.DebugLevel, logrus.TraceLevel:
		clr = termenv.ANSIBrightBlack
	case logrus.InfoLevel:
		clr = termenv.ANSIBrightCyan
	case logrus.WarnLevel:
		clr = termenv.ANSIBrightYellow
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		clr = termenv.ANSIBrightMagenta
	default:
		panic(fmt.Errorf("logrus.Level know supported"))
	}
	return clr
}

func (c DevCtlFormatter) Format(entry *logrus.Entry) (bytes []byte, err error) {
	sb := &strings.Builder{}
	lvlStyle := lLevel(entry.Level).LevelStyle()

	msg := entry.Message
	lLevel(entry.Level).LevelStyle()
	lvl := lvlStyle.String()

	normalizeStr := func(str string) []byte {
		styled := termenv.
			String(str).
			Foreground(termenv.ForegroundColor()).
			Faint()
		normalizedString := styled.String()
		return []byte(normalizedString)
	}

	prefix_start := normalizeStr("[ ")
	prefix_level := lvl
	prefix_end := normalizeStr(" ]")
	prefix_spacer := normalizeStr("\t")

	sb.Write(prefix_start)
	sb.Write([]byte(prefix_level))
	sb.Write(prefix_end)
	sb.Write(prefix_spacer)
	sb.Write([]byte(msg))
	sb.Write([]byte("\n"))

	text := sb.String()
	bytes = []byte(text)

	return bytes, err
}
