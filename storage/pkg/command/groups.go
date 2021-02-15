package command

import (
	"context"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/cs3org/reva/cmd/revad/runtime"
	"github.com/gofrs/uuid"
	"github.com/micro/cli/v2"
	"github.com/oklog/run"
	"github.com/owncloud/ocis/storage/pkg/config"
	"github.com/owncloud/ocis/storage/pkg/flagset"
	"github.com/owncloud/ocis/storage/pkg/server/debug"
)

// Groups is the entrypoint for the sharing command.
func Groups(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:  "groups",
		Usage: "Start groups service",
		Flags: flagset.GroupsWithConfig(cfg),
		Before: func(c *cli.Context) error {
			cfg.Reva.Groups.Services = c.StringSlice("service")

			return nil
		},
		Action: func(c *cli.Context) error {
			logger := NewLogger(cfg)

			if cfg.Tracing.Enabled {
				switch t := cfg.Tracing.Type; t {
				case "agent":
					logger.Error().
						Str("type", t).
						Msg("Reva only supports the jaeger tracing backend")

				case "jaeger":
					logger.Info().
						Str("type", t).
						Msg("configuring storage to use the jaeger tracing backend")

				case "zipkin":
					logger.Error().
						Str("type", t).
						Msg("Reva only supports the jaeger tracing backend")

				default:
					logger.Warn().
						Str("type", t).
						Msg("Unknown tracing backend")
				}

			} else {
				logger.Debug().
					Msg("Tracing is not enabled")
			}

			var (
				gr          = run.Group{}
				ctx, cancel = context.WithCancel(context.Background())
				//metrics     = metrics.New()
			)

			defer cancel()

			{
				uuid := uuid.Must(uuid.NewV4())
				pidFile := path.Join(os.TempDir(), "revad-"+c.Command.Name+"-"+uuid.String()+".pid")

				rcfg := map[string]interface{}{
					"core": map[string]interface{}{
						"max_cpus":             cfg.Reva.Groups.MaxCPUs,
						"tracing_enabled":      cfg.Tracing.Enabled,
						"tracing_endpoint":     cfg.Tracing.Endpoint,
						"tracing_collector":    cfg.Tracing.Collector,
						"tracing_service_name": c.Command.Name,
					},
					"shared": map[string]interface{}{
						"jwt_secret": cfg.Reva.JWTSecret,
					},
					"grpc": map[string]interface{}{
						"network": cfg.Reva.Groups.GRPCNetwork,
						"address": cfg.Reva.Groups.GRPCAddr,
						// TODO build services dynamically
						"services": map[string]interface{}{
							"groupprovider": map[string]interface{}{
								"driver": cfg.Reva.Groups.Driver,
								"drivers": map[string]interface{}{
									"json": map[string]interface{}{
										"groups": cfg.Reva.Groups.JSON,
									},
									"ldap": map[string]interface{}{
										"hostname":        cfg.Reva.LDAP.Hostname,
										"port":            cfg.Reva.LDAP.Port,
										"base_dn":         cfg.Reva.LDAP.BaseDN,
										"groupfilter":     cfg.Reva.LDAP.GroupFilter,
										"attributefilter": cfg.Reva.LDAP.GroupAttributeFilter,
										"findfilter":      cfg.Reva.LDAP.GroupFindFilter,
										"memberfilter":    cfg.Reva.LDAP.GroupMemberFilter,
										"bind_username":   cfg.Reva.LDAP.BindDN,
										"bind_password":   cfg.Reva.LDAP.BindPassword,
										"idp":             cfg.Reva.LDAP.IDP,
										"schema": map[string]interface{}{
											"dn":          "dn",
											"gid":         cfg.Reva.LDAP.Schema.GID,
											"mail":        cfg.Reva.LDAP.Schema.Mail,
											"displayName": cfg.Reva.LDAP.Schema.DisplayName,
											"cn":          cfg.Reva.LDAP.Schema.CN,
											"gidNumber":   cfg.Reva.LDAP.Schema.GIDNumber,
										},
									},
									"rest": map[string]interface{}{
										"client_id":                      cfg.Reva.UserGroupRest.ClientID,
										"client_secret":                  cfg.Reva.UserGroupRest.ClientSecret,
										"redis_address":                  cfg.Reva.UserGroupRest.RedisAddress,
										"redis_username":                 cfg.Reva.UserGroupRest.RedisUsername,
										"redis_password":                 cfg.Reva.UserGroupRest.RedisPassword,
										"group_members_cache_expiration": cfg.Reva.Groups.GroupMembersCacheExpiration,
										"id_provider":                    cfg.Reva.UserGroupRest.IDProvider,
										"api_base_url":                   cfg.Reva.UserGroupRest.APIBaseURL,
										"oidc_token_endpoint":            cfg.Reva.UserGroupRest.OIDCTokenEndpoint,
										"target_api":                     cfg.Reva.UserGroupRest.TargetAPI,
									},
								},
							},
						},
					},
				}

				gr.Add(func() error {
					runtime.RunWithOptions(
						rcfg,
						pidFile,
						runtime.WithLogger(&logger.Logger),
					)
					return nil
				}, func(_ error) {
					logger.Info().
						Str("server", c.Command.Name).
						Msg("Shutting down server")

					cancel()
				})
			}

			{
				server, err := debug.Server(
					debug.Name(c.Command.Name+"-debug"),
					debug.Addr(cfg.Reva.Users.DebugAddr),
					debug.Logger(logger),
					debug.Context(ctx),
					debug.Config(cfg),
				)

				if err != nil {
					logger.Info().
						Err(err).
						Str("server", c.Command.Name+"-debug").
						Msg("Failed to initialize server")

					return err
				}

				gr.Add(func() error {
					return server.ListenAndServe()
				}, func(_ error) {
					ctx, timeout := context.WithTimeout(ctx, 5*time.Second)

					defer timeout()
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						logger.Info().
							Err(err).
							Str("server", c.Command.Name+"-debug").
							Msg("Failed to shutdown server")
					} else {
						logger.Info().
							Str("server", c.Command.Name+"-debug").
							Msg("Shutting down server")
					}
				})
			}

			{
				stop := make(chan os.Signal, 1)

				gr.Add(func() error {
					signal.Notify(stop, os.Interrupt)

					<-stop

					return nil
				}, func(err error) {
					close(stop)
					cancel()
				})
			}

			return gr.Run()
		},
	}
}
