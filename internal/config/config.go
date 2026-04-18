// Package config loads the Fireflies CLI configuration from env and TOML file.
//
// Precedence:
//  1. FIREFLIES_API_KEY env var (wins for the active profile's API key)
//  2. ~/.config/fireflies/config.toml
//  3. interactive login (written back to the config file)
package config

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sync"

	"github.com/pelletier/go-toml/v2"

	ferr "github.com/fvdm-otinga/fireflies-cli/internal/errors"
)

// defaultHost is the expected GraphQL host; any other host triggers a
// one-shot warning because the API key will be transmitted to it.
const defaultHost = "api.fireflies.ai"

// endpointWarnOnce ensures we only emit the non-default-endpoint warning
// a single time per process, no matter how many profiles are inspected.
var endpointWarnOnce sync.Once

// validateEndpoint enforces https:// on any configured GraphQL endpoint
// and warns once if the host is not the default. It returns a config
// CLIError if the endpoint is malformed or non-HTTPS.
func validateEndpoint(endpoint string) error {
	if endpoint == "" {
		return nil
	}
	u, err := url.Parse(endpoint)
	if err != nil {
		return ferr.Usage(fmt.Sprintf("invalid endpoint %q: %v", endpoint, err))
	}
	if u.Scheme != "https" {
		return ferr.Usage(fmt.Sprintf(
			"endpoint must use https:// (got %q); API key would be sent in clear",
			endpoint))
	}
	if u.Host != defaultHost {
		endpointWarnOnce.Do(func() {
			fmt.Fprintf(os.Stderr,
				"warning: using non-default GraphQL endpoint %s; API key will be sent to this host\n",
				endpoint)
		})
	}
	return nil
}

const (
	EnvAPIKey     = "FIREFLIES_API_KEY"
	EnvConfigPath = "FIREFLIES_CONFIG"
	DefaultFile   = "fireflies/config.toml"
)

type Profile struct {
	APIKey   string `toml:"api_key"`
	Endpoint string `toml:"endpoint,omitempty"`
}

type File struct {
	Active   string             `toml:"active"`
	Profiles map[string]Profile `toml:"profiles"`
}

// Loader reads and writes the configuration.
type Loader struct {
	path   string
	data   *File
	loaded bool
}

func New() *Loader { return &Loader{} }

func NewWithPath(p string) *Loader { return &Loader{path: p} }

// Path returns the resolved config file path.
func (l *Loader) Path() (string, error) {
	if l.path != "" {
		return l.path, nil
	}
	if env := os.Getenv(EnvConfigPath); env != "" {
		l.path = env
		return env, nil
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("resolve config dir: %w", err)
	}
	l.path = filepath.Join(dir, DefaultFile)
	return l.path, nil
}

// Load reads the config file (creating it with no profiles if missing).
func (l *Loader) Load() error {
	p, err := l.Path()
	if err != nil {
		return err
	}
	b, err := os.ReadFile(p)
	if err != nil {
		if os.IsNotExist(err) {
			l.data = &File{Profiles: map[string]Profile{}}
			l.loaded = true
			return nil
		}
		return fmt.Errorf("read config: %w", err)
	}
	f := &File{Profiles: map[string]Profile{}}
	if err := toml.Unmarshal(b, f); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}
	if f.Profiles == nil {
		f.Profiles = map[string]Profile{}
	}
	// Reject obviously unsafe endpoints at load time; the warning for a
	// non-default host is deferred to Profile() so it fires only for the
	// profile actually in use.
	for _, p := range f.Profiles {
		if p.Endpoint == "" {
			continue
		}
		u, err := url.Parse(p.Endpoint)
		if err != nil {
			return ferr.Usage(fmt.Sprintf("invalid endpoint %q: %v", p.Endpoint, err))
		}
		if u.Scheme != "https" {
			return ferr.Usage(fmt.Sprintf(
				"endpoint must use https:// (got %q); API key would be sent in clear",
				p.Endpoint))
		}
	}
	l.data = f
	l.loaded = true
	return nil
}

// Profile returns a named profile. If name is empty, returns the active profile.
// Env var FIREFLIES_API_KEY always overrides the profile's api_key.
func (l *Loader) Profile(name string) (Profile, error) {
	if !l.loaded {
		if err := l.Load(); err != nil {
			return Profile{}, err
		}
	}
	if name == "" {
		name = l.data.Active
		if name == "" {
			name = "default"
		}
	}
	p, ok := l.data.Profiles[name]
	if !ok {
		p = Profile{}
	}
	if env := os.Getenv(EnvAPIKey); env != "" {
		p.APIKey = env
	}
	if p.APIKey == "" {
		return Profile{}, ferr.Auth(fmt.Sprintf(
			"no API key found. Set %s env var or run `fireflies auth login`", EnvAPIKey))
	}
	if err := validateEndpoint(p.Endpoint); err != nil {
		return Profile{}, err
	}
	return p, nil
}

// Save writes the current config to disk with mode 0600.
func (l *Loader) Save() error {
	p, err := l.Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return fmt.Errorf("mkdir config dir: %w", err)
	}
	b, err := toml.Marshal(l.data)
	if err != nil {
		return fmt.Errorf("marshal config: %w", err)
	}
	if err := os.WriteFile(p, b, 0600); err != nil {
		return fmt.Errorf("write config: %w", err)
	}
	return nil
}

// SetProfile upserts a profile and optionally sets it active.
func (l *Loader) SetProfile(name string, p Profile, makeActive bool) error {
	if !l.loaded {
		if err := l.Load(); err != nil {
			return err
		}
	}
	if l.data.Profiles == nil {
		l.data.Profiles = map[string]Profile{}
	}
	l.data.Profiles[name] = p
	if makeActive || l.data.Active == "" {
		l.data.Active = name
	}
	return l.Save()
}

// DeleteProfile removes a named profile.
func (l *Loader) DeleteProfile(name string) error {
	if !l.loaded {
		if err := l.Load(); err != nil {
			return err
		}
	}
	delete(l.data.Profiles, name)
	if l.data.Active == name {
		l.data.Active = ""
	}
	return l.Save()
}

// All returns a snapshot of the loaded file.
func (l *Loader) All() (*File, error) {
	if !l.loaded {
		if err := l.Load(); err != nil {
			return nil, err
		}
	}
	cp := *l.data
	return &cp, nil
}
