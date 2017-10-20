package disk

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/ipfs/ipfs-cluster/config"
)

const configKey = "disk"

// Default values for disk Config
const (
	DefaultMetricTTL  = 30 * time.Second
	DefaultMetricType = MetricFreeSpace
)

// String returns a string representation for MetricType.
func (t MetricType) String() string {
	switch t {
	case MetricFreeSpace:
		return "freespace"
	case MetricRepoSize:
		return "reposize"
	}
	return ""
}

// Config is used to initialize an Informer and customize
// the type and parameters of the metric it produces.
type Config struct {
	config.Saver

	MetricTTL time.Duration
	Type      MetricType
}

type jsonConfig struct {
	MetricTTL string `json:"metric_ttl"`
	Type      string `json:"metric_type"`
}

// ConfigKey returns a human-friendly identifier for this type of Metric.
func (cfg *Config) ConfigKey() string {
	return configKey
}

// Default initializes this Config with sensible values.
func (cfg *Config) Default() error {
	cfg.MetricTTL = DefaultMetricTTL
	cfg.Type = DefaultMetricType
	return nil
}

// Validate checks that the fields of this Config have working values,
// at least in appearance.
func (cfg *Config) Validate() error {
	if cfg.MetricTTL <= 0 {
		return errors.New("disk.metric_ttl is invalid")
	}

	if _, ok := metricToRPC[cfg.Type]; !ok {
		return errors.New("disk.metric_type is invalid")
	}
	return nil
}

// LoadJSON reads the fields of this Config from a JSON byteslice as
// generated by ToJSON.
func (cfg *Config) LoadJSON(raw []byte) error {
	jcfg := &jsonConfig{}
	err := json.Unmarshal(raw, jcfg)
	if err != nil {
		logger.Error("Error unmarshaling disk informer config")
		return err
	}

	t, _ := time.ParseDuration(jcfg.MetricTTL)
	cfg.MetricTTL = t

	switch jcfg.Type {
	case "reposize":
		cfg.Type = MetricRepoSize
	case "freespace":
		cfg.Type = MetricFreeSpace
	default:
		return errors.New("disk.metric_type is invalid")
	}

	return cfg.Validate()
}

// ToJSON generates a JSON-formatted human-friendly representation of this
// Config.
func (cfg *Config) ToJSON() (raw []byte, err error) {
	jcfg := &jsonConfig{}

	jcfg.MetricTTL = cfg.MetricTTL.String()
	jcfg.Type = cfg.Type.String()

	raw, err = config.DefaultJSONMarshal(jcfg)
	return
}
