package config

import (
	"encoding/json"
	"fmt"
	"gitlab.alipay-inc.com/afe/mosn/pkg/api/v2"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

)

//global instance for load & dump
var ConfigPath string
var config MOSNConfig

type FilterConfig struct {
	Type   string                 `json:"type,omitempty"`
	Config map[string]interface{} `json:"config,omitempty"`
}

type AccessLogConfig struct {
	LogPath   string `json:"log_path,omitempty"`
	LogFormat string `json:"log_format,omitempty"`
}

type ListenerConfig struct {
	Name           string         `json:"name,omitempty"`
	Address        string         `json:"address,omitempty"`
	BindToPort     bool           `json:"bind_port"`
	NetworkFilters []FilterConfig `json:"network_filters,service_registry"`
	StreamFilters  []FilterConfig `json:"stream_filters,omitempty"`

	//logger
	LogPath  string `json:"log_path,omitempty"`
	LogLevel string `json:"log_level,omitempty"`

	//access log
	AccessLogs []AccessLogConfig `json:"access_logs,omitempty"`

	// only used in http2 case
	DisableConnIo bool `json:"disable_conn_io"`
}

type ServerConfig struct {
	//default logger
	DefaultLogPath  string `json:"default_log_path,omitempty"`
	DefaultLogLevel string `json:"default_log_level,omitempty"`

	//graceful shutdown config
	GracefulTimeout DurationConfig `json:"graceful_timeout"`

	//go processor number
	Processor int

	Listeners []ListenerConfig `json:"listeners,omitempty"`
}

type HostConfig struct {
	Address  string `json:"address,omitempty"`
	Hostname string `json:"hostname,omitempty"`
	Weight   uint32 `json:"weight,omitempty"`
}

type HealthCheckConfig struct {
	Timeout            DurationConfig
	HealthyThreshold   uint32 `json:"healthy_threshold"`
	UnhealthyThreshold uint32 `json:"unhealthy_threshold"`
	Interval           DurationConfig
	IntervalJitter     DurationConfig `json:"interval_jitter"`
	CheckPath          string         `json:"check_path,omitempty"`
	ServiceName        string         `json:"service_name,omitempty"`
}

type ClusterSpecConfig struct {
	Subscribes []SubscribeSpecConfig `json:"subscribe,omitempty"`
}

type SubscribeSpecConfig struct {
	ServiceName string `json:"service_name,omitempty"`
}

type ClusterConfig struct {
	Name              string
	Type              string
	SubType           string             `json:"sub_type"`
	LbType            string             `json:"lb_type"`
	MaxRequestPerConn uint32
	CircuitBreakers   v2.CircuitBreakers `json:"circuit_breakers"`
	HealthCheck       v2.HealthCheck     `json:"health_check,omitempty"` //v2.HealthCheck
	ClusterSpecConfig ClusterSpecConfig  `json:"spec,omitempty"`         //	ClusterSpecConfig
	Hosts             []v2.Host          `json:"hosts,omitempty"`        //v2.Host
}

type ClusterManagerConfig struct {
	AutoDiscovery bool            `json:"auto_discovery"`
	Clusters      []ClusterConfig `json:"clusters,omitempty"`
}

type ServiceRegistryConfig struct {
	ServiceAppInfo ServiceAppInfoConfig   `json:"application"`
	ServicePubInfo []ServicePubInfoConfig `json:"publish_info,omitempty"`
}

type ServiceAppInfoConfig struct {
	AntShareCloud bool   `json:"ant_share_cloud"`
	DataCenter    string `json:"data_center,omitempty"`
	AppName       string `json:"app_name,omitempty"`
}

type ServicePubInfoConfig struct {
	ServiceName string `json:"service_name,omitempty"`
	PubData     string `json:"pub_data,omitempty"`
}

type MOSNConfig struct {
	Servers         []ServerConfig        `json:"servers,omitempty"`         //server config
	ClusterManager  ClusterManagerConfig  `json:"cluster_manager,omitempty"` //cluster config
	ServiceRegistry ServiceRegistryConfig `json:"service_registry"`          //service registry config, used by service discovery module
	//tracing config
	RawDynamicResources json.RawMessage `json:"dynamic_resources,omitempty"`  //dynamic_resources raw message
	RawStaticResources  json.RawMessage `json:"static_resources,omitempty"`   //static_resources raw message
}

//wrapper for time.Duration, so time config can be written in '300ms' or '1h' format
type DurationConfig struct {
	time.Duration
}

func (d *DurationConfig) UnmarshalJSON(b []byte) (err error) {
	d.Duration, err = time.ParseDuration(strings.Trim(string(b), `"`))
	return
}

func (d DurationConfig) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

func Load(path string) *MOSNConfig {
	log.Println("load config from : ", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("load config failed, ", err)
		os.Exit(1)
	}
	ConfigPath, _ = filepath.Abs(path)
	// todo delete
	//ConfigPath = "../../resource/mosn_config_dump_result.json"

	err = json.Unmarshal(content, &config)
	if err != nil {
		log.Fatalln("json unmarshal config failed, ", err)
		os.Exit(1)
	}
	return &config
}
