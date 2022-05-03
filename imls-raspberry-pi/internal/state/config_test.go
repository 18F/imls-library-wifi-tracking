package state

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigSuite struct {
	suite.Suite
}

var configDBPath = "/tmp/config-test.sql"

func (suite *ConfigSuite) SetupTest() {
	os.Create(configDBPath)
	os.Chmod(configDBPath, 0777)
	SetConfigAtPath(configDBPath)
}

func (suite *ConfigSuite) AfterTest(suiteName, testName string) {
	dc := GetConfig()
	dc.Close()
	// ensure a clean run.
	os.Remove(configDBPath)
}

func (suite *ConfigSuite) TestConfigDefaults() {
	dc := GetConfig()
	_, filename, _, _ := runtime.Caller(0)
	path := filepath.Dir(filename)
	durationsPath := filepath.Join(path, "test", "durations.sqlite")
	dc.SetDurationsPath(durationsPath)
	var expected = []string{"local:stderr", "local:tmp", "api:directus"}
	result := dc.GetLoggers()
	for i := 0; i < 3; i += 1 {
		if result[i] != expected[i] {
			suite.Fail("loggers were not equal")
		}
	}
	if dc.GetDurationsDatabase().GetPath() != durationsPath {
		suite.Fail("duration path was not equal")
	}
	os.Remove(durationsPath)
}

func (suite *ConfigSuite) TestConfigWrite() {
	dc := GetConfig()
	dc.SetDeviceTag("a random string")
	result := dc.GetDeviceTag()
	if result != "a random string" {
		suite.Fail("write was not reflected")
	}
}

func TestConfigSuite(t *testing.T) {
	suite.Run(t, new(ConfigSuite))
}
