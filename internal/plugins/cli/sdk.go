package cli

type SdkPlugin interface {
	Download(version string) (err error)
}

type sdkPlugin struct {
	baseUrl string
}

/*

response, err := http.Get("https://golang.org/dl/go1.16.3.darwin-amd64.tar.gz")
if err != nil {
t.Fatal(err)
}
defer response.Body.Close()
size := response.ContentLength

//var bar = progressbar.DefaultBytes(size, "downloading...")
var bar = progressbar.NewOptions64(size,
	progressbar.OptionSetWriter(os.Stdout),
	progressbar.OptionEnableColorCodes(true),
	progressbar.OptionFullWidth(),
	progressbar.OptionSetPredictTime(true),
	progressbar.OptionShowCount(),
	progressbar.OptionSetDescription("[cyan][1/3][reset] downloading go sdk..."),
	progressbar.OptionSetTheme(progressbar.Theme{
		Saucer:        "[green]=[reset]",
		SaucerHead:    "[green]>[reset]",
		SaucerPadding: " ",
		BarStart:      "[",
		BarEnd:        "]",
	}),
	progressbar.OptionShowBytes(true),
)

io.Copy(bar, response.Body)
}
*/
