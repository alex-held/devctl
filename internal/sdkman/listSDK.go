package sdkman

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	
	"github.com/jedib0t/go-pretty/table"
	"github.com/jedib0t/go-pretty/text"
)

type TableFormatter func(f fmt.State, verb rune)

type SDKList []string

func (s SDKList) Format(f fmt.State, verb rune) {

}

func (s SDKList) String() string {
	tw := table.NewWriter()
	tw.SetCaption("sdkman - list")
	tw.SetStyle(table.StyleColoredDark)
	tw.SetColumnConfigs([]table.ColumnConfig{
		{
			Number:       1,
			Align:        text.AlignCenter,
			AlignHeader:  text.AlignCenter,
			Colors:       text.Colors{text.FgCyan},
			VAlign:       text.VAlignMiddle,
			VAlignFooter: text.VAlignMiddle,
			VAlignHeader: text.VAlignMiddle,
			WidthMin:     50,
			WidthMax:     50,
		},
		{
			Number:       2,
			Align:        text.AlignCenter,
			AlignHeader:  text.AlignCenter,
			Colors:       text.Colors{text.FgMagenta},
			VAlign:       text.VAlignMiddle,
			VAlignFooter: text.VAlignMiddle,
			VAlignHeader: text.VAlignMiddle,
			WidthMin:     100,
			WidthMax:     200,
		},
	})
	tw.SetTitle("Available SDK's")
	tw.AppendHeader(table.Row{"#", "sdk"})
	for i, sdk := range s {
		tw.AppendRow(table.Row{i, sdk})
	}
	return tw.Render()
}

type ListAllSDKService service

// CreateListAllAvailableSDKURI gets all available SDK and returns them as an array of strings
// https://api.sdkman.io/2/candidates/all
func (s *ListAllSDKService) ListAllSDK(ctx context.Context) (sdks SDKList, resp *http.Response, err error) {
	// CreateListAllAvailableSDKURI creates the URI to list all available SDK
	// https://api.sdkman.io/2/candidates/all
	req, err := s.client.NewRequest(ctx, "GET", "candidates/all", nil)
	if err != nil {
		return nil, nil, err
	}
	resp, err = s.client.client.Do(req)
	if err != nil {
		return nil, resp, err
	}
	defer resp.Body.Close()
	
	responseBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, resp, err
	}
	
	sdkCSV := string(responseBodyBytes)
	sdkList := strings.Split(sdkCSV, ",")
	return sdkList, resp, nil
}
