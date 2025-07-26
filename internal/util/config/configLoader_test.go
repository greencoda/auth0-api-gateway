package config_test

import (
	"os"
	"testing"

	config_util "github.com/greencoda/auth0-api-gateway/internal/util/config"
	"github.com/greencoda/confiq"
	yaml_loader "github.com/greencoda/confiq/loaders/yaml"
	. "github.com/smartystreets/goconvey/convey"

	_ "embed"
)

type testConfig struct {
	Name  string `cfg:"name"`
	Port  int    `cfg:"port,default=8080"`
	Debug bool   `cfg:"debug,default=false"`
}

type testConfig_NoDefaults struct {
	Name  string `cfg:"name"`
	Port  int    `cfg:"port"`
	Debug bool   `cfg:"debug"`
}

//go:embed testdata/config_valid.yaml
var config_valid string

//go:embed testdata/config_invalid.yaml
var config_invalid string

//go:embed testdata/config_valid_prefix.yaml
var config_valid_prefix string

//go:embed testdata/config_valid_noPrefix.yaml
var config_valid_noPrefix string

//go:embed testdata/config_valid_partialData.yaml
var config_valid_partial_data string

//go:embed testdata/config_valid_different.yaml
var config_valid_different string

func Test_LoadConfigYAML(t *testing.T) {
	Convey("When loading config from YAML file", t, func() {
		Convey("With valid YAML file", func() {
			tmpFile, err := os.CreateTemp("", "test-config-*.yaml")
			So(err, ShouldBeNil)
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(config_valid)
			So(err, ShouldBeNil)
			tmpFile.Close()

			configSet, err := config_util.LoadConfigYAML(config_util.ConfigFilename(tmpFile.Name()))
			So(err, ShouldBeNil)
			So(configSet, ShouldNotBeNil)

			var config testConfig
			err = configSet.Decode(&config, confiq.FromPrefix("test"))
			So(err, ShouldBeNil)

			So(config, ShouldResemble, testConfig{
				Name:  "test-service",
				Port:  9090,
				Debug: true,
			})
		})

		Convey("With non-existent file", func() {
			configSet, err := config_util.LoadConfigYAML(config_util.ConfigFilename("non-existent-file.yaml"))
			So(err, ShouldNotBeNil)
			So(configSet, ShouldBeNil)
		})

		Convey("With invalid YAML file", func() {
			// Create a temporary invalid YAML file
			tmpFile, err := os.CreateTemp("", "invalid-config-*.yaml")
			So(err, ShouldBeNil)
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(config_invalid)
			So(err, ShouldBeNil)
			tmpFile.Close()

			configSet, err := config_util.LoadConfigYAML(config_util.ConfigFilename(tmpFile.Name()))
			So(err, ShouldNotBeNil)
			So(configSet, ShouldBeNil)
		})
	})
}

func Test_LoadConfigFromSetWithPrefix(t *testing.T) {
	Convey("When loading config from ConfigSet with prefix", t, func() {
		Convey("With valid config set and prefix", func() {
			configSet := confiq.New()
			err := configSet.Load(
				yaml_loader.Load().FromString(config_valid_prefix),
			)
			So(err, ShouldBeNil)

			config, err := config_util.LoadConfigFromSetWithPrefix[testConfig](configSet, "test")
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)

			So(config, ShouldResemble, &testConfig{
				Name:  "prefixed-service",
				Port:  3000,
				Debug: true,
			})
		})

		Convey("With valid config set and empty prefix", func() {
			configSet := confiq.New()
			err := configSet.Load(
				yaml_loader.Load().FromString(config_valid_noPrefix),
			)
			So(err, ShouldBeNil)

			config, err := config_util.LoadConfigFromSetWithPrefix[testConfig](configSet, "")
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)

			So(config, ShouldResemble, &testConfig{
				Name:  "no-prefix-service",
				Port:  4000,
				Debug: false,
			})
		})

		Convey("With nil config set", func() {
			config, err := config_util.LoadConfigFromSetWithPrefix[testConfig](nil, "test")
			So(err, ShouldEqual, config_util.ErrNoConfigSet)
			So(config, ShouldBeNil)
		})

		Convey("With valid config set but non-existent prefix", func() {
			configSet := confiq.New()
			err := configSet.Load(
				yaml_loader.Load().FromString(config_valid_prefix),
			)
			So(err, ShouldBeNil)

			config, err := config_util.LoadConfigFromSetWithPrefix[testConfig](configSet, "nonexistent")
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)

			So(config, ShouldResemble, &testConfig{
				Name:  "",
				Port:  8080,  // default value
				Debug: false, // default value
			})
		})

		Convey("With empty config set", func() {
			configSet := confiq.New()
			config, err := config_util.LoadConfigFromSetWithPrefix[testConfig](configSet, "test")
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)

			So(config, ShouldResemble, &testConfig{
				Name:  "",
				Port:  8080,
				Debug: false,
			})
		})

		Convey("With config set containing partial data", func() {
			configSet := confiq.New()
			err := configSet.Load(
				yaml_loader.Load().FromString(config_valid_partial_data),
			)
			So(err, ShouldBeNil)

			config, err := config_util.LoadConfigFromSetWithPrefix[testConfig](configSet, "test")
			So(err, ShouldBeNil)
			So(config, ShouldNotBeNil)

			So(config, ShouldResemble, &testConfig{
				Name:  "partial-service",
				Port:  8080,  // default value
				Debug: false, // default value
			})
		})

		Convey("With config set containing different keys and no default values", func() {
			configSet := confiq.New()
			err := configSet.Load(
				yaml_loader.Load().FromString(config_valid_different),
			)
			So(err, ShouldBeNil)

			config, err := config_util.LoadConfigFromSetWithPrefix[testConfig_NoDefaults](configSet, "test")
			So(err, ShouldNotBeNil)
			So(config, ShouldBeNil)
		})
	})
}
