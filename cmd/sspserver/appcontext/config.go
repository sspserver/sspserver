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
		Listen string `default:":8080" field:"listen" json:"listen" yaml:"listen" cli:"http-listen" env:"SERVER_HTTP_LISTEN"`

		// ReadTimeout for all income request in milliseconds
		ReadTimeout time.Duration `default:"120s" field:"read_timeout" json:"read_timeout" yaml:"read_timeout" env:"SERVER_HTTP_READ_TIMEOUT"`

		// ReadTimeout for all request in milliseconds
		WriteTimeout time.Duration `default:"120s" field:"write_timeout" json:"write_timeout" yaml:"write_timeout" env:"SERVER_HTTP_WRITE_TIMEOUT"`
	} `json:"HTTP" yaml:"HTTP"`

	Profile struct {
		Mode   string `json:"mode" yaml:"mode" default:"" env:"SERVER_PROFILE_MODE"`
		Listen string `json:"listen" yaml:"listen" default:"" env:"SERVER_PROFILE_LISTEN"`
	} `json:"profile" yaml:"profile"`

	// Hostname of the server
	Hostname string `field:"hostname" json:"hostname" yaml:"hostname" env:"SERVER_HOSTNAME"`

	// All IPs of this server
	IPs []string `field:"ips" json:"ips" yaml:"ips"`

	// Datacenter helps sinchronize balances between DC
	// We have to split spends cross DC because statistic delivery time is different.
	// Moreover perhapse we have to synchronise balances by hosts
	Datacenter struct {
		// Code of DC for 2 chars like: EU, US, AS, etc.
		Code string `field:"code" json:"code" yaml:"code" env:"DATACENTER_ISO2_GEO_CODE" default:"XX"`

		// We can split the traffic between DC instancess according with countries closest to this DC
		// For example EU, US DCs
		// For EU DC all pure EU campaigns have to receive 70% of traffic and other 30% to other DC
		Countries []string `field:"countries" json:"countries" yaml:"countries" env:"DATACENTER_AVAILABLE_COUNTRIES"`

		// Count of services in this DC
		ServiceCount int `field:"service_count" json:"service_count" yaml:"service_count" env:"DATACENTER_SERVICE_COUNT"`
	} `field:"datacenter" yaml:"datacenter" json:"datacenter"`

	// Service discovery
	Registry struct {
		// Connection to the registry
		// Examples:
		//   - consul://host:port?arg=value
		//   - etcd://host:port?arg=value
		//   - zookeeper://host:port?arg=value
		Connection string `field:"connection" json:"connection" yaml:"connection" env:"REGISTRY_CONNECTION"`

		// Hostname of the service in the registry
		Hostname string `field:"hostname" json:"hostname" yaml:"hostname" env:"REGISTRY_HOSTNAME"`

		// Port of the service in the registry
		Port int `field:"port" json:"port" yaml:"port" env:"REGISTRY_PORT"`
	} `field:"registry" yaml:"registry" json:"registry"`
}

type adeventer struct {
	EventQueue struct {
		Connection string `field:"connection" json:"connection" yaml:"connection" env:"EVENTSTREAM_EVENT_QUEUE_CONNECTION"`
	} `yaml:"event_queue" json:"event_queue"`

	WinQueue struct {
		Connection string `field:"connection" json:"connection" yaml:"connection" env:"EVENTSTREAM_WINS_QUEUE_CONNECTION"`
	} `yaml:"wins_queue" json:"wins_queue"`
}

type adstorage struct {
	// Connection to database of path to directory with data or DB connection
	// Examples:
	//   - fs://directory/path
	//   - postgresql://login:password@hostname:port/dbname
	Connection string `field:"connection" json:"connection" yaml:"connection" env:"ADSTORAGE_CONNECTION"`

	// Formats list of format objects
	Formats string `field:"formats" json:"formats" yaml:"formats" env:"ADSTORAGE_FORMATS"`

	// Zones list of zone objects
	Zones string `field:"zones" json:"zones" yaml:"zones" env:"ADSTORAGE_ZONES"`

	// Sources list of ad sources
	Sources string `field:"sources" json:"sources" yaml:"sources" env:"ADSTORAGE_SOURCES"`
}

type adLogic struct {
	Direct struct {
		DefaultURL string `field:"default_url" yaml:"default_url" json:"default_url" env:"ADSERVER_LOGIC_DIRECT_DEFAULT_URL"`
	} `field:"direct" yaml:"direct" json:"direct"`
}

type AdServerConfig struct {
	// Storage of the advertisement
	Storage adstorage `field:"storage" yaml:"storage" json:"storage"`

	// Event pipelines
	Event adeventer `field:"event" yaml:"event" json:"event"`

	// Default tracker domain name
	TrackerHost string `field:"pixel_host" yaml:"pixel_host" json:"pixel_host" env:"ADSERVER_TRACKER_HOST"`

	// Default CDN domain name
	CDNDomain string `field:"cdn" yaml:"cdn" json:"cdn" env:"ADSERVER_CDN_DOMAIN" default:"localhost:8090"`

	// Lib CDN domain name
	LibDomain string `field:"libcdn" yaml:"libcdn" json:"libcdn" env:"ADSERVER_CDN_LIB_DOMAIN" default:"localhost:8090"`

	// Configuration of Source accessor
	AdSource adsourceConfig `field:"adsource" yaml:"adsource" json:"adsource"`

	// EventPipeline of the results
	EventPipeline adeventer `field:"event_pipeline" yaml:"event_pipeline" json:"event_pipeline"`

	// Logic of adserver behavier
	Logic adLogic `field:"logic" yaml:"logic" json:"logic"`
}

type adsourceConfig struct {
	// Maximum amount of requests from source accessor
	MaxParallelRequests int `field:"max_parallel_requests" json:"max_parallel_requests" yaml:"max_parallel_requests" env:"ADSOURCE_REQUEST_MAX_PARALLEL_REQUESTS" default:"10"`

	// Maximal request timeout (one for all internal requests)
	RequestTimeout int64 `field:"request_timeout" json:"request_timeout" yaml:"request_timeout" env:"ADSOURCE_REQUEST_TIMEOUT" default:"500"`
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
	Connect          string        `field:"connect" json:"connect" yaml:"connect" env:"PERSON_CONNECT"`
	MaxConn          int           `field:"max_conn" json:"max_conn" yaml:"max_conn" env:"PERSON_MAX_CONN" default:"10"`
	RequestTimeout   time.Duration `field:"request_timeout" json:"request_timeout" yaml:"request_timeout" env:"PERSON_REQUEST_TIMEOUT" default:"1s"`
	KeepAliveTimeout time.Duration `field:"keepalive_timeout" json:"keepalive_timeout" yaml:"keepalive_timeout" env:"PERSON_KEEPALIVE_TIMEOUT" default:"300s"`

	UUIDCookieName   string        `field:"uuid_cookie_name" json:"uuid_cookie_name" yaml:"uuid_cookie_name" env:"PERSON_UUID_COOKIE_NAME"`
	SessiCookiedName string        `field:"session_cookie_name" json:"session_cookie_name" yaml:"session_cookie_name" env:"PERSON_SESSION_COOKIE_NAME"`
	SessionLifetime  time.Duration `field:"session_lifetime" json:"session_lifetime" yaml:"session_lifetime" env:"PERSON_SESSION_LIFETIME" default:"72h"`
}

// Config contains all configurations of the application
type Config struct {
	ServiceName    string `json:"service_name" yaml:"service_name" env:"SERVICE_NAME" default:"sspserver"`
	DatacenterName string `json:"datacenter_name" yaml:"datacenter_name" env:"DC_NAME" default:"??"`
	Hostname       string `json:"hostname" yaml:"hostname" env:"HOSTNAME" default:""`
	Hostcode       string `json:"hostcode" yaml:"hostcode" env:"HOSTCODE" default:""`

	LogAddr    string `default:"" env:"LOG_ADDR"`
	LogLevel   string `default:"error" env:"LOG_LEVEL"`
	LogEncoder string `json:"log_encoder" env:"LOG_ENCODER"`

	// Server config
	Server ServerConfig `field:"server" json:"server" yaml:"server"`
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
