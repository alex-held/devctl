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
	var clr = termenv.RGBColor("#17AFE1")
	switch logrus.Level(l) {
	case logrus.DebugLevel, logrus.TraceLevel:
		clr = termenv.RGBColor("#17AFE1")
	case logrus.InfoLevel:
		clr = termenv.RGBColor("#17AFE1")
	case logrus.WarnLevel:
		clr = termenv.RGBColor("#F9FE00")
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		clr = termenv.RGBColor("#FF3D6D")
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

	prefixStart := normalizeStr("[ ")
	prefixLevel := lvl
	prefixnd := normalizeStr(" ]")
	prefixSpacer := normalizeStr("\t")

	sb.Write(prefixStart)
	sb.Write([]byte(prefixLevel))
	sb.Write(prefixnd)
	sb.Write(prefixSpacer)
	sb.Write([]byte(msg))
	sb.Write([]byte("\n"))

	text := sb.String()
	bytes = []byte(text)

	return bytes, err
}
