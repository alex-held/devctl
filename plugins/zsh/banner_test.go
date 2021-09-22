package zsh

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBanner_WithInternet(t *testing.T) {
	expected := "  ______   ______   .___  ___. .______    __       _______ .___________. __    ______   .__   __.      _______.\n /      | /  __  \\  |   \\/   | |   _  \\  |  |     |   ____||           ||  |  /  __  \\  |  \\ |  |     /       |\n|  ,----'|  |  |  | |  \\  /  | |  |_)  | |  |     |  |__   `---|  |----`|  | |  |  |  | |   \\|  |    |   (----`\n|  |     |  |  |  | |  |\\/|  | |   ___/  |  |     |   __|      |  |     |  | |  |  |  | |  . `  |     \\   \\    \n|  `----.|  `--'  | |  |  |  | |  |      |  `----.|  |____     |  |     |  | |  `--'  | |  |\\   | .----)   |   \n \\______| \\______/  |__|  |__| | _|      |_______||_______|    |__|     |__|  \\______/  |__| \\__| |_______/    \n"

	actual := GenerateBanner("Completions")

	if *genGoldenMaster {
		_ = os.WriteFile("testdata/banner_test.zsh.golden", []byte(actual), os.ModePerm)
		return
	}

	assert.Equal(t, expected, actual)
}
