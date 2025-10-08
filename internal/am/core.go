package am

import (
	"context"
	"sync"
	"time"
)

type Core interface {
	SetOpts(opts ...Option)
	Naming
	Logging
	Configuring
	Lifecycle
}

type Naming interface {
	Name() string
	SetName(name string)
}

type Logging interface {
	Log() Logger
	SetLog(log Logger)
}

type Configuring interface {
	Cfg() *Config
	SetCfg(cfg *Config)
}

type Lifecycle interface {
	Setup(ctx context.Context) error
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type BaseCore struct {
	name string
	log  Logger
	cfg  *Config
}

var (
	defaultLogger     Logger
	defaultLoggerOnce sync.Once
)

func getDefaultLogger() Logger {
	defaultLoggerOnce.Do(func() {
		defaultLogger = NewLogger("info")
	})
	return defaultLogger
}

// NewCore creates a new BaseCore instance with the provided opts.
func NewCore(name string, opts ...Option) *BaseCore {
	core := &BaseCore{name: name}
	for _, opt := range opts {
		opt(core)
	}
	return core
}

// SetOpts sets the options in BaseCore.
func (c *BaseCore) SetOpts(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

// Name returns the name in BaseCore.
func (c *BaseCore) Name() string {
	return c.name
}

// SetName sets the name in BaseCore.
func (c *BaseCore) SetName(name string) {
	c.name = name
}

func (c *BaseCore) Log() Logger {
	if c.log == nil {
		return getDefaultLogger()
	}
	return c.log
}

// SetLog sets the Logger in BaseCore.
func (c *BaseCore) SetLog(log Logger) {
	c.log = log
}

// Cfg returns the Config in BaseCore.
func (c *BaseCore) Cfg() *Config {
	return c.cfg
}

// SetCfg sets the Config in BaseCore.
func (c *BaseCore) SetCfg(cfg *Config) {
	c.cfg = cfg
}

// Setup is the default implementation for the Setup method in BaseCore.
func (c *BaseCore) Setup(ctx context.Context) error {
	return nil
}

// Start is the default implementation for the Start method in BaseCore.
func (c *BaseCore) Start(ctx context.Context) error {
	return nil
}

// Stop is the default implementation for the Stop method in BaseCore.
func (c *BaseCore) Stop(ctx context.Context) error {
	return nil
}

// Option defines a type for setting optional parameters in BaseCore.
type Option func(Core)


// WithConfigValue sets a specific key-value pair in the Config.
func WithConfigValue(key string, value interface{}) Option {
	return func(c Core) {
		if c.Cfg() == nil {
			c.SetCfg(NewConfig())
		}
		c.Cfg().Set(key, value)
	}
}

// XParams contains the essential dependencies for most components.
type XParams struct {
	Cfg *Config
	Log Logger
}

// NewCoreWithParams creates a new BaseCore instance with XParams.
func NewCoreWithParams(name string, params XParams) *BaseCore {
	return &BaseCore{
		name: name,
		log:  params.Log,
		cfg:  params.Cfg,
	}
}

// NewCoreWithParamsAndOpts creates a new BaseCore with XParams and additional options.
func NewCoreWithParamsAndOpts(name string, params XParams, opts ...Option) *BaseCore {
	core := NewCoreWithParams(name, params)
	for _, opt := range opts {
		opt(core)
	}
	return core
}

// Now returns the current time.
func Now() time.Time {
	return time.Now()
}
