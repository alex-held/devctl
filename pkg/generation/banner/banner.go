package banner

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

	"github.com/common-nighthawk/go-figure"
)

var (
	curlStrategy = 0
	goStrategy   = 1
)
var defaultBannerStrategy = goStrategy

func GenerateBanner(banner string, kind OutputKind) (out string) {
	switch defaultBannerStrategy {
	case curlStrategy:
		out = generateBanner(http.DefaultClient, banner)
	case goStrategy:
		fallthrough
	default:
		out = Banner(banner, kind)
	}
	return out
}

func generateBanner(c *http.Client, banner string) string {
	defaultBanner := fmt.Sprintf("# %s", banner)

	url := fmt.Sprintf("https://devops.datenkollektiv.de/renderBannerTxt?text=%s&font=starwars", banner)
	resp, err := c.Get(url)
	if err != nil {
		return defaultBanner
	}

	defer resp.Body.Close()

	b := &bytes.Buffer{}
	_, err = io.Copy(b, resp.Body)
	if err != nil {
		return defaultBanner
	}

	bannerStr := b.String()

	i := strings.Index(bannerStr, "\n")
	if (i+1)%2 != 0 {
		i++
	}

	line := strings.Repeat("- ", i/2+1)
	banner = fmt.Sprintf("%s\n%s%s\n", line, bannerStr, line)

	bannerLines := strings.Split(banner, "\n")
	banner = ""
	for _, bl := range bannerLines {
		banner += "# " + bl + "\n"
	}

	li := strings.LastIndex(banner, "-")
	return banner[:li+1]
}

//go:embed fonts/starwars.flf
var starwarsFont string

func getStarwarsFontReader() io.Reader {
	return strings.NewReader(starwarsFont)
}

func Banner(text string, kind OutputKind) string {
	fig := figure.NewFigureWithFont(text, getStarwarsFontReader(), true)
	banner := fig.String()

	b := newBuilder(banner, kind)

	prettified := b.Prettify()

	return prettified
}

type bannerBuilder struct {
	b      *strings.Builder
	lines  []string
	banner string
	kind   OutputKind
}

type OutputKind int

const (
	KIND_SHELL OutputKind = iota
	KIND_YAML
	KIND_GO
)

var (
	sanitizers = map[OutputKind]string{
		KIND_SHELL: "# %v\n",
		KIND_YAML:  "# %v\n",
		KIND_GO:    "// %v\n",
	}
)

type headerLines []string

func (h headerLines) Len() int {
	return len(h)
}
func (h headerLines) Less(i, j int) bool {
	return len(h[i]) > len(h[j])
}
func (h headerLines) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (b *bannerBuilder) decorationLine() (line string) {
	lines := strings.Split(b.sanitizeBanner(), "\n")

	l := headerLines(lines)
	sort.Sort(l)
	repeatCount := len(l[0]) - 1

	if len(l[0])%2 != 0 {
		repeatCount++
	}

	line = b.sanitizeLine(strings.Repeat(" -", repeatCount/2))
	return line
}

func newBuilder(text string, kind OutputKind) *bannerBuilder {
	return &bannerBuilder{
		b:      &strings.Builder{},
		banner: text,
		kind:   kind,
		lines:  strings.Split(text, "\n"),
	}
}

func (b *bannerBuilder) sanitizeLine(str string) string {
	s := sanitizers[b.kind]
	return fmt.Sprintf(s, str)
}

func (b *bannerBuilder) Prettify() (out string) {
	line := b.decorationLine()

	if _, err := b.b.WriteString(line); err != nil {
		return b.sanitizeLine(b.banner)
	}

	banner := b.sanitizeBanner()

	if _, err := b.b.WriteString(banner); err != nil {
		return b.sanitizeLine(b.banner)
	}

	if _, err := b.b.WriteString(line); err != nil {
		return b.sanitizeLine(b.banner)
	}

	return b.b.String()
}

func (b *bannerBuilder) sanitizeBanner() string {
	out := ""
	for _, line := range b.lines {
		out += b.sanitizeLine(line)
	}
	if out == "" {
		return b.sanitizeLine(b.banner)
	}
	return out
}
