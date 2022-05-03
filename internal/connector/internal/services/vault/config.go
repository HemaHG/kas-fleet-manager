package vault

import (
	"github.com/bf2fc6cc711aee1a0c2a/kas-fleet-manager/pkg/shared"
	"github.com/spf13/pflag"
)

type Config struct {
	// Used for OSD Cluster creation with OCM
	Kind                string `json:"kind"`
	AccessKey           string `json:"access_key"`
	AccessKeyFile       string `json:"access_key_file"`
	SecretAccessKey     string `json:"secret_access_key"`
	SecretAccessKeyFile string `json:"secret_access_key_file"`
	SecretPrefix        string `json:"secret_prefix"`
	Region              string `json:"region"`
}

func NewConfig() *Config {
	return &Config{
		Kind:                KindTmp,
		AccessKeyFile:       "secrets/vault.accesskey",
		SecretAccessKeyFile: "secrets/vault.secretaccesskey",
		Region:              DefaultRegion,
		SecretPrefix:        "managed-connectors",
	}
}

func (c *Config) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.Kind, "vault-kind", c.Kind, "The kind of vault to use: aws|tmp")
	fs.StringVar(&c.AccessKeyFile, "vault-access-key-file", c.AccessKeyFile, "File containing vault access key")
	fs.StringVar(&c.SecretAccessKeyFile, "vault-secret-access-key-file", c.SecretAccessKeyFile, "File containing vault secret access key")
	fs.StringVar(&c.SecretPrefix, "vault-secret-prefix", c.SecretPrefix, "Prefix to use before all managed connectors secret names in aws vault")
	fs.StringVar(&c.Region, "vault-region", c.Region, "The region of the vault")
}

func (c *Config) ReadFiles() error {
	if c.Kind == KindAws {
		err := shared.ReadFileValueString(c.AccessKeyFile, &c.AccessKey)
		if err != nil {
			return err
		}
		err = shared.ReadFileValueString(c.SecretAccessKeyFile, &c.SecretAccessKey)
		if err != nil {
			return err
		}
	}
	return nil
}
