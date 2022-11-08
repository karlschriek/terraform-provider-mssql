package mssql

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/betr-io/terraform-provider-mssql/mssql/model"
	"github.com/betr-io/terraform-provider-mssql/sql"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"strings"
)

type mssqlProvider struct {
	factory model.ConnectorFactory
	logger  *zerolog.Logger
	data    *schema.ResourceData
}

const (
	providerLogFile = "terraform-provider-mssql.log"
)

var (
	defaultTimeout = schema.DefaultTimeout(30 * time.Second)
)

func New(version, commit string) func() *schema.Provider {
	return func() *schema.Provider {
		return Provider(sql.GetFactory())
	}
}

func Provider(factory model.ConnectorFactory) *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
				DefaultFunc: schema.EnvDefaultFunc("MSSQL_HOST", nil),
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  DefaultPort,
			},
			"login": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"login", "azure_login"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("MSSQL_USERNAME", nil),
						},
						"password": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc("MSSQL_PASSWORD", nil),
						},
					},
				},
			},
			"azure_login": {
				Type:         schema.TypeList,
				MaxItems:     1,
				Optional:     true,
				ExactlyOneOf: []string{"login", "azure_login"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_id": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("MSSQL_TENANT_ID", nil),
						},
						"client_id": {
							Type:        schema.TypeString,
							Required:    true,
							DefaultFunc: schema.EnvDefaultFunc("MSSQL_CLIENT_ID", nil),
						},
						"client_secret": {
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
							DefaultFunc: schema.EnvDefaultFunc("MSSQL_CLIENT_SECRET", nil),
						},
					},
				},
			},
			"debug": {
				Type:        schema.TypeBool,
				Description: fmt.Sprintf("Enable provider debug logging (logs to file %s)", providerLogFile),
				Optional:    true,
				Default:     false,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"mssql_login": resourceLogin(),
			"mssql_user":  resourceUser(),
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ConfigureContextFunc: func(ctx context.Context, data *schema.ResourceData) (interface{}, diag.Diagnostics) {
			return providerConfigure(ctx, data, factory)
		},
	}
}

func providerConfigure(ctx context.Context, data *schema.ResourceData, factory model.ConnectorFactory) (model.Provider, diag.Diagnostics) {
	isDebug := data.Get("debug").(bool)
	logger := newLogger(isDebug)

	logger.Info().Msg("Created provider")

	return mssqlProvider{factory: factory, logger: logger, data: data}, nil
}

func (p mssqlProvider) GetConnector(prefix string) (interface{}, error) {
	return p.factory.GetConnector(prefix, p.data)
}

func (p mssqlProvider) ResourceLogger(resource, function string) zerolog.Logger {
	return p.logger.With().Str("resource", resource).Str("func", function).Logger()
}

func (p mssqlProvider) DataSourceLogger(datasource, function string) zerolog.Logger {
	return p.logger.With().Str("datasource", datasource).Str("func", function).Logger()
}

func (p mssqlProvider) GetProviderData() *schema.ResourceData {
	return p.data
}

func newLogger(isDebug bool) *zerolog.Logger {
	var writer io.Writer = nil
	logLevel := zerolog.Disabled
	if isDebug {
		f, err := os.OpenFile(providerLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Err(err).Msg("error opening file")
		}
		writer = f
		logLevel = zerolog.DebugLevel
	}
	logger := zerolog.New(writer).Level(logLevel).With().Timestamp().Logger()
	return &logger
}
