// Package appcontext provides config options
package appcontext

/**
 ██████╗ ██████╗ ███╗   ██╗███████╗██╗ ██████╗
██╔════╝██╔═══██╗████╗  ██║██╔════╝██║██╔════╝
██║     ██║   ██║██╔██╗ ██║█████╗  ██║██║  ███╗
██║     ██║   ██║██║╚██╗██║██╔══╝  ██║██║   ██║
╚██████╗╚██████╔╝██║ ╚████║██║     ██║╚██████╔╝
 ╚═════╝ ╚═════╝ ╚═╝  ╚═══╝╚═╝     ╚═╝ ╚═════╝
*/

import (
	"encoding/json"
	"strings"
	"time"
)

// Webserver config
type ServerConfig struct {
	HTTP struct {
		// Listen the [host]:port
		Listen string `default:":8080" field:"listen" json:"listen" yaml:"listen" cli:"http-listen" env:"LISTEN"`

		// ReadTimeout for all income request in milliseconds
		ReadTimeout time.Duration `default:"120s" field:"read_timeout" json:"read_timeout" yaml:"read_timeout" env:"READ_TIMEOUT"`

		// ReadTimeout for all request in milliseconds
		WriteTimeout time.Duration `default:"120s" field:"write_timeout" json:"write_timeout" yaml:"write_timeout" env:"WRITE_TIMEOUT"`
	} `json:"HTTP" yaml:"HTTP" envPrefix:"HTTP_"`

	Profile struct {
		Mode   string `json:"mode" yaml:"mode" default:"" env:"MODE"`
		Listen string `json:"listen" yaml:"listen" default:"" env:"LISTEN"`
	} `json:"profile" yaml:"profile" envPrefix:"PROFILE_"`

	// Hostname of the server
	Hostname string `field:"hostname" json:"hostname" yaml:"hostname" env:"HOSTNAME"`

	// All IPs of this server
	IPs []string `field:"ips" json:"ips" yaml:"ips"`

	// Datacenter helps sinchronize balances between DC
	// We have to split spends cross DC because statistic delivery time is different.
	// Moreover perhapse we have to synchronise balances by hosts
	Datacenter struct {
		// Code of DC for 2 chars like: EU, US, AS, etc.
		Code string `field:"code" json:"code" yaml:"code" env:"ISO2_GEO_CODE" default:"XX"`

		// We can split the traffic between DC instancess according with countries closest to this DC
		// For example EU, US DCs
		// For EU DC all pure EU campaigns have to receive 70% of traffic and other 30% to other DC
		Countries []string `field:"countries" json:"countries" yaml:"countries" env:"AVAILABLE_COUNTRIES"`

		// Count of services in this DC
		ServiceCount int `field:"service_count" json:"service_count" yaml:"service_count" env:"SERVICE_COUNT"`
	} `field:"datacenter" yaml:"datacenter" json:"datacenter" envPrefix:"DATACENTER_"`

	// Service discovery
	Registry struct {
		// Connection to the registry
		// Examples:
		//   - consul://host:port?arg=value
		//   - etcd://host:port?arg=value
		//   - zookeeper://host:port?arg=value
		Connection string `field:"connection" json:"connection" yaml:"connection" env:"CONNECTION"`

		// Hostname of the service in the registry
		Hostname string `field:"hostname" json:"hostname" yaml:"hostname" env:"HOSTNAME"`

		// Port of the service in the registry
		Port int `field:"port" json:"port" yaml:"port" env:"PORT"`
	} `field:"registry" yaml:"registry" json:"registry" envPrefix:"REGISTRY_"`
}

type adeventer struct {
	EventQueue struct {
		Connection string `field:"connection" json:"connection" yaml:"connection" env:"CONNECTION"`
	} `yaml:"event_queue" json:"event_queue" envPrefix:"EVENT_QUEUE_"`

	WinQueue struct {
		Connection string `field:"connection" json:"connection" yaml:"connection" env:"CONNECTION"`
	} `yaml:"wins_queue" json:"wins_queue" envPrefix:"WINS_QUEUE_"`
}

type adstorage struct {
	// Connection to database of path to directory with data or DB connection
	// Examples:
	//   - fs://directory/path
	//   - postgresql://login:password@hostname:port/dbname
	Connection string `field:"connection" json:"connection" yaml:"connection" env:"CONNECTION"`

	// Formats list of format objects
	Formats string `field:"formats" json:"formats" yaml:"formats" env:"FORMATS"`

	// Zones list of zone objects
	Zones string `field:"zones" json:"zones" yaml:"zones" env:"ZONES"`

	// Sources list of ad sources
	Sources string `field:"sources" json:"sources" yaml:"sources" env:"SOURCES"`
}

type adLogic struct {
	Direct struct {
		DefaultURL string `field:"default_url" yaml:"default_url" json:"default_url" env:"ADSERVER_LOGIC_DIRECT_DEFAULT_URL"`
	} `field:"direct" yaml:"direct" json:"direct"`
}

type adInfoConfig struct {
	ComplaintAdURL string `json:"complaint_ad_url" yaml:"complaint_ad_url" env:"COMPLAINT_AD_URL"`
	AboutAdURL     string `json:"about_ad_url" yaml:"about_ad_url" env:"ABOUT_AD_URL"`
}

type AdServerConfig struct {
	// Storage of the advertisement
	Storage adstorage `field:"storage" yaml:"storage" json:"storage" envPrefix:"ADSTORAGE_"`

	// Event pipelines
	// Event adeventer `field:"event" yaml:"event" json:"event" envPrefix:"ADEVENT_"`

	// Default tracker domain name
	TrackerHost string `field:"pixel_host" yaml:"pixel_host" json:"pixel_host" env:"TRACKER_HOST"`

	// Default CDN domain name
	CDNDomain string `field:"cdn" yaml:"cdn" json:"cdn" env:"CDN_DOMAIN" default:"localhost:8090"`

	// Lib CDN domain name
	LibDomain string `field:"libcdn" yaml:"libcdn" json:"libcdn" env:"CDN_LIB_DOMAIN" default:"localhost:8090"`

	// Configuration of Source accessor
	AdSource adsourceConfig `field:"adsource" yaml:"adsource" json:"adsource" envPrefix:"ADSOURCE_"`

	// EventPipeline of the results
	EventPipeline adeventer `field:"event_pipeline" yaml:"event_pipeline" json:"event_pipeline" envPrefix:"EVENTPIPELINE_"`

	// Logic of adserver behavier
	Logic adLogic `field:"logic" yaml:"logic" json:"logic" envPrefix:"LOGIC_"`

	// Information about ads
	Info adInfoConfig `field:"info" yaml:"info" json:"info" envPrefix:"ADINFO_"`
}

type adsourceConfig struct {
	// Maximum amount of requests from source accessor
	MaxParallelRequests int `field:"max_parallel_requests" json:"max_parallel_requests" yaml:"max_parallel_requests" env:"REQUEST_MAX_PARALLEL_REQUESTS" default:"10"`

	// Maximal request timeout (one for all internal requests)
	RequestTimeout int64 `field:"request_timeout" json:"request_timeout" yaml:"request_timeout" env:"REQUEST_TIMEOUT" default:"500"`
}

type PersonConfig struct {
	// Connect to the source of information about GEO and Device
	//
	// Macros:
	//  {ip} - IP address of current request
	//  {ua} - User-Agent string value of current request
	//
	// Supported by default HTTP protocol
	// 	http://geodevice.domain.com/get-geo-and-device-info?ip={ip}&user-agent={ua}
	// GRPC connection
	//  grpc://hostname:1234
	// or UNIX socket + GRPC
	//  grpc+unix://hostname:1234
	Connect          string        `field:"connect" json:"connect" yaml:"connect" env:"CONNECT"`
	MaxConn          int           `field:"max_conn" json:"max_conn" yaml:"max_conn" env:"MAX_CONN" default:"10"`
	RequestTimeout   time.Duration `field:"request_timeout" json:"request_timeout" yaml:"request_timeout" env:"REQUEST_TIMEOUT" default:"1s"`
	KeepAliveTimeout time.Duration `field:"keepalive_timeout" json:"keepalive_timeout" yaml:"keepalive_timeout" env:"KEEPALIVE_TIMEOUT" default:"300s"`

	UUIDCookieName   string        `field:"uuid_cookie_name" json:"uuid_cookie_name" yaml:"uuid_cookie_name" env:"UUID_COOKIE_NAME"`
	SessiCookiedName string        `field:"session_cookie_name" json:"session_cookie_name" yaml:"session_cookie_name" env:"SESSION_COOKIE_NAME"`
	SessionLifetime  time.Duration `field:"session_lifetime" json:"session_lifetime" yaml:"session_lifetime" env:"SESSION_LIFETIME" default:"72h"`
}

// Config contains all configurations of the application
type Config struct {
	ServiceName    string `json:"service_name" yaml:"service_name" env:"SERVICE_NAME" default:"sspserver"`
	DatacenterName string `json:"datacenter_name" yaml:"datacenter_name" env:"DC_NAME" default:"??"`
	Hostname       string `json:"hostname" yaml:"hostname" env:"HOSTNAME" default:""`
	Hostcode       string `json:"hostcode" yaml:"hostcode" env:"HOSTCODE" default:""`

	LogAddr    string `json:"log_addr" default:"" env:"LOG_ADDR"`
	LogLevel   string `json:"log_level" default:"error" env:"LOG_LEVEL"`
	LogEncoder string `json:"log_encoder" env:"LOG_ENCODER"`

	// Server config
	Server ServerConfig `field:"server" json:"server" yaml:"server" envPrefix:"SERVER_"`

	// Configuration of Advertisement server
	AdServer AdServerConfig `field:"adserver" yaml:"adserver" json:"adserver" envPrefix:"ADSERVER_"`

	// Person data extraction service
	Person PersonConfig `field:"person" yaml:"person" json:"person" envPrefix:"PERSON_"`
}

// String implementation of Stringer interface
func (cfg *Config) String() (res string) {
	if data, err := json.MarshalIndent(cfg, "", "  "); nil != err {
		res = `{"error":"` + err.Error() + `"}`
	} else {
		res = string(data)
	}
	return
}

// IsDebug mode
func (cfg *Config) IsDebug() bool {
	return strings.EqualFold(cfg.LogLevel, "debug")
}
