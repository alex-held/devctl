package zsh

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func GenerateBanner(banner string) string {
	return generateBanner(http.DefaultClient, banner)
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
