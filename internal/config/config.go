package config

import (
	"os"
	"runtime"
	"time"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Server      ServerConfig      `yaml:"server"`
	Zeo         ZeoConfig         `yaml:"zeo"`
	Concurrency ConcurrencyConfig `yaml:"concurrency"`
	Cache       CacheConfig       `yaml:"cache"`
	Logging     LoggingConfig     `yaml:"logging"`
}

type ServerConfig struct {
	Port               string        `yaml:"port"`
	Host               string        `yaml:"host"`
	ReadTimeout        time.Duration `yaml:"read_timeout"`
	WriteTimeout       time.Duration `yaml:"write_timeout"`
	MaxMultipartMemory int64         `yaml:"max_multipart_memory"`
}

type ZeoConfig struct {
	ExecutablePath string        `yaml:"executable_path"`
	Workdir        string        `yaml:"workdir"`
	Timeout        time.Duration `yaml:"timeout"`
}

type ConcurrencyConfig struct {
	MaxWorkers           int   `yaml:"max_workers"`
	MaxQueueSize         int   `yaml:"max_queue_size"`
	RateLimitPerIP       int   `yaml:"rate_limit_per_ip"`
	MaxFileSize          int64 `yaml:"max_file_size"`
	MaxConcurrentUploads int   `yaml:"max_concurrent_uploads"`
}

type CacheConfig struct {
	Enabled   bool          `yaml:"enabled"`
	TTL       time.Duration `yaml:"ttl"`
	MaxSizeMB int64         `yaml:"max_size_mb"`
	Shards    int           `yaml:"shards"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Format string `yaml:"format"`
	Output string `yaml:"output"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Set defaults
	if cfg.Concurrency.MaxWorkers <= 0 {
		cfg.Concurrency.MaxWorkers = runtime.NumCPU()
	}
	if cfg.Cache.Shards <= 0 {
		cfg.Cache.Shards = 32
	}
	if cfg.Cache.TTL < 0 {
		cfg.Cache.TTL = 3600 * time.Second
	}

	return &cfg, nil
}

func LoadDefaultConfig() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port:               "8080",
			Host:               "0.0.0.0",
			ReadTimeout:        30 * time.Second,
			WriteTimeout:       30 * time.Second,
			MaxMultipartMemory: 32 << 20, // 32MB
		},
		Zeo: ZeoConfig{
			ExecutablePath: "network",
			Workdir:        "./workspace",
			Timeout:        5 * time.Minute,
		},
		Concurrency: ConcurrencyConfig{
			MaxWorkers:           runtime.NumCPU(),
			MaxQueueSize:         1000,
			RateLimitPerIP:       10,
			MaxFileSize:          int64(100) << 20, // 100MB
			MaxConcurrentUploads: 50,
		},
		Cache: CacheConfig{
			Enabled:   true,
			TTL:       3600 * time.Second,
			MaxSizeMB: 1024,
			Shards:    32,
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
			Output: "stdout",
		},
	}, nil
}
