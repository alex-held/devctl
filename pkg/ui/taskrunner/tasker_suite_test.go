package taskrunner

import (
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	"github.com/onsi/gomega"
)

//TestTaskerSuite tests Tasker, ConditionalTask, Task and TaskRunner
func TestTaskerSuite(t *testing.T) {
	config.GinkgoConfig.DryRun = false
	config.GinkgoConfig.EmitSpecProgress = true
	config.GinkgoConfig.DebugParallel = false
	config.DefaultReporterConfig.Verbose = false
	config.DefaultReporterConfig.Succinct = true
	config.DefaultReporterConfig.NoisyPendings = true

	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Tasker")
}
