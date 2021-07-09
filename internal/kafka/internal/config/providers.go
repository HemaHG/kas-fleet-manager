package config

import (
	"errors"
	"fmt"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/environments"
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/shared"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v2"
)

type Region struct {
	Name    string `json:"name"`
	Default bool   `json:"default"`
}

type RegionList []Region

func (rl RegionList) GetByName(regionName string) (Region, bool) {
	for _, r := range rl {
		if r.Name == regionName {
			return r, true
		}
	}
	return Region{}, false
}

func (rl RegionList) String() string {
	var names []string
	for _, r := range rl {
		names = append(names, r.Name)
	}
	return fmt.Sprint(names)
}

type Provider struct {
	Name    string     `json:"name"`
	Default bool       `json:"default"`
	Regions RegionList `json:"regions"`
}

type ProviderList []Provider

func (pl ProviderList) GetByName(providerName string) (Provider, bool) {
	for _, p := range pl {
		if p.Name == providerName {
			return p, true
		}
	}
	return Provider{}, false
}

func (pl ProviderList) String() string {
	var names []string
	for _, p := range pl {
		names = append(names, p.Name)
	}
	return fmt.Sprint(names)
}

type ProviderConfiguration struct {
	SupportedProviders ProviderList `yaml:"supported_providers"`
}

type ProviderConfig struct {
	ProvidersConfig     ProviderConfiguration `json:"providers"`
	ProvidersConfigFile string                `json:"providers_config_file"`
}

func NewSupportedProvidersConfig() *ProviderConfig {
	return &ProviderConfig{
		ProvidersConfigFile: "config/provider-configuration.yaml",
	}
}

var _ environments.ServiceValidator = &ProviderConfig{}

func (c *ProviderConfig) Validate() error {
	providerDefaultCount := 0
	for _, p := range c.ProvidersConfig.SupportedProviders {
		if err := p.Validate(); err != nil {
			return err
		}
		if p.Default {
			providerDefaultCount++
		}
	}
	if providerDefaultCount != 1 {
		return fmt.Errorf("expected 1 default provider in provider list, got %d", providerDefaultCount)
	}
	return nil
}

func (provider Provider) Validate() error {
	defaultCount := 0
	for _, p := range provider.Regions {
		if p.Default {
			defaultCount++
		}
	}
	if defaultCount != 1 {
		return fmt.Errorf("expected 1 default region in provider %s, got %d", provider.Name, defaultCount)
	}
	return nil
}

func (c *ProviderConfig) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.ProvidersConfigFile, "providers-config-file", c.ProvidersConfigFile, "SupportedProviders configuration file")
}

func (c *ProviderConfig) ReadFiles() error {
	return readFileProvidersConfig(c.ProvidersConfigFile, &c.ProvidersConfig)
}

// Read the contents of file into the providers config
func readFileProvidersConfig(file string, val *ProviderConfiguration) error {
	fileContents, err := shared.ReadFile(file)
	if err != nil {
		return err
	}
	return yaml.UnmarshalStrict([]byte(fileContents), val)
}

func (c ProviderList) GetDefault() (Provider, error) {
	for _, p := range c {
		if p.Default {
			return p, nil
		}
	}
	return Provider{}, errors.New("no default provider found in list of supported providers")
}

func (provider Provider) GetDefaultRegion() (Region, error) {
	for _, r := range provider.Regions {
		if r.Default {
			return r, nil
		}
	}
	return Region{}, fmt.Errorf("no default region found for provider %s", provider.Name)
}

func (provider Provider) IsRegionSupported(regionName string) bool {
	_, ok := provider.Regions.GetByName(regionName)
	return ok
}

func (h ProviderList) IsRegionSupportedForProvider(provider string, id string) bool {
	if p, ok := h.GetByName(provider); ok {
		return p.IsRegionSupported(id)
	}
	return false
}