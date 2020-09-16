package config

import (
	"fmt"
	"regexp"
	"strings"

	rhiconfig "github.com/redhatinsights/app-common-go/pkg/api/v1"
	"github.com/spf13/viper"
)

// IngressConfig represents the runtime configuration
type IngressConfig struct {
	Hostname             string
	MaxSize              int64
	StageBucket          string
	Auth                 bool
	KafkaBrokers         []string
	KafkaGroupID         string
	KafkaTrackerTopic    string
	ValidTopics          []string
	Port                 int
	Profile              bool
	OpenshiftBuildCommit string
	Version              string
	MinioDev             bool
	MinioEndpoint        string
	MinioAccessKey       string
	MinioSecretKey       string
	Debug                bool
	DebugUserAgent       *regexp.Regexp
}

// Get returns an initialized IngressConfig
func Get() *IngressConfig {

	options := viper.New()
	options.SetDefault("MaxSize", 10*1024*1024)
	options.SetDefault("StageBucket", "available")
	options.SetDefault("Auth", true)
	options.SetDefault("KafkaGroupID", "ingress")
	options.SetDefault("KafkaTrackerTopic", "platform.payload-status")
	options.SetDefault("ValidTopics", "unit")
	options.SetDefault("OpenshiftBuildCommit", "notrunninginopenshift")
	options.SetDefault("Profile", false)
	options.SetDefault("Debug", false)
	options.SetDefault("DebugUserAgent", `unspecified`)
	options.SetEnvPrefix("INGRESS")
	options.AutomaticEnv()
	kubenv := viper.New()
	kubenv.SetDefault("Openshift_Build_Commit", "notrunninginopenshift")
	kubenv.SetDefault("Hostname", "Hostname_Unavailable")
	kubenv.AutomaticEnv()
	return &IngressConfig{
		Hostname:             kubenv.GetString("Hostname"),
		MaxSize:              options.GetInt64("MaxSize"),
		StageBucket:          options.GetString("StageBucket"),
		Auth:                 options.GetBool("Auth"),
		KafkaBrokers:         []string{fmt.Sprintf("%v:%v", rhiconfig.LoadedConfig.Kafka.Brokers[0].Hostname, *(rhiconfig.LoadedConfig.Kafka.Brokers[0].Port))},
		KafkaGroupID:         options.GetString("KafkaGroupID"),
		KafkaTrackerTopic:    options.GetString("KafkaTrackerTopic"),
		ValidTopics:          strings.Split(options.GetString("ValidTopics"), ","),
		Port:                 int(rhiconfig.LoadedConfig.WebPort),
		Profile:              options.GetBool("Profile"),
		Debug:                options.GetBool("Debug"),
		DebugUserAgent:       regexp.MustCompile(options.GetString("DebugUserAgent")),
		OpenshiftBuildCommit: kubenv.GetString("Openshift_Build_Commit"),
		Version:              "1.0.8",
		MinioDev:             options.GetBool("MinioDev"),
		MinioEndpoint:        rhiconfig.LoadedConfig.ObjectStore.Endpoint,
		MinioAccessKey:       rhiconfig.LoadedConfig.ObjectStore.AccessKey,
		MinioSecretKey:       rhiconfig.LoadedConfig.ObjectStore.SecretKey,
	}
}
