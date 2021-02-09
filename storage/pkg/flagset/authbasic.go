package flagset

import (
	"github.com/micro/cli/v2"
	"github.com/owncloud/ocis/storage/pkg/config"
)

// AuthBasicWithConfig applies cfg to the root flagset
func AuthBasicWithConfig(cfg *config.Config) []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:        "debug-addr",
			Value:       "0.0.0.0:9147",
			Usage:       "Address to bind debug server",
			EnvVars:     []string{"STORAGE_AUTH_BASIC_DEBUG_ADDR"},
			Destination: &cfg.Reva.AuthBasic.DebugAddr,
		},
		&cli.StringFlag{
			Name:        "auth-driver",
			Value:       "ldap",
			Usage:       "auth driver: 'demo', 'json' or 'ldap'",
			EnvVars:     []string{"STORAGE_AUTH_DRIVER"},
			Destination: &cfg.Reva.AuthProvider.Driver,
		},
		&cli.StringFlag{
			Name:        "auth-json",
			Value:       "",
			Usage:       "Path to users.json file",
			EnvVars:     []string{"STORAGE_AUTH_JSON"},
			Destination: &cfg.Reva.AuthProvider.JSON,
		},
		&cli.StringFlag{
			Name:        "network",
			Value:       "tcp",
			Usage:       "Network to use for the storage auth-basic service, can be 'tcp', 'udp' or 'unix'",
			EnvVars:     []string{"STORAGE_AUTH_BASIC_GRPC_NETWORK"},
			Destination: &cfg.Reva.AuthBasic.GRPCNetwork,
		},
		&cli.StringFlag{
			Name:        "addr",
			Value:       "0.0.0.0:9146",
			Usage:       "Address to bind storage service",
			EnvVars:     []string{"STORAGE_AUTH_BASIC_GRPC_ADDR"},
			Destination: &cfg.Reva.AuthBasic.GRPCAddr,
		},
		&cli.StringSliceFlag{
			Name:    "service",
			Value:   cli.NewStringSlice("authprovider"),
			Usage:   "--service authprovider [--service otherservice]",
			EnvVars: []string{"STORAGE_AUTH_BASIC_SERVICES"},
		},
	}

	flags = append(flags, TracingWithConfig(cfg)...)
	flags = append(flags, DebugWithConfig(cfg)...)
	flags = append(flags, SecretWithConfig(cfg)...)
	flags = append(flags, LDAPWithConfig(cfg)...)

	return flags
}
