package internal

import (
	"fmt"

	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/mitchellh/mapstructure"
)

type Plugin struct {
	environment  string
	provider     string
	globalConfig *Auth0Config
	siteConfigs  map[string]*Auth0Config
	enabled      bool
}

func NewAuth0Plugin() schema.MachComposerPlugin {
	state := &Plugin{
		provider:    "1.26.0",
		siteConfigs: map[string]*Auth0Config{},
	}
	return plugin.NewPlugin(&schema.PluginSchema{
		Identifier: "Auth0",

		Configure: state.Configure,
		IsEnabled: state.IsEnabled,

		GetValidationSchema: state.GetValidationSchema,

		// Config
		SetGlobalConfig: state.SetGlobalConfig,
		SetSiteConfig:   state.SetSiteConfig,

		// Renders
		RenderTerraformProviders: state.TerraformRenderProviders,
		RenderTerraformResources: state.TerraformRenderResources,
		RenderTerraformComponent: state.RenderTerraformComponent,
	})
}

func (p *Plugin) Configure(environment string, provider string) error {
	p.environment = environment
	if provider != "" {
		p.provider = provider
	}
	return nil
}

func (p *Plugin) IsEnabled() bool {
	return p.enabled
}

func (p *Plugin) GetValidationSchema() (*schema.ValidationSchema, error) {
	result := getSchema()
	return result, nil
}

func (p *Plugin) SetGlobalConfig(data map[string]any) error {
	cfg := Auth0Config{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.globalConfig = &cfg
	p.enabled = true
	return nil
}

func (p *Plugin) SetSiteConfig(site string, data map[string]any) error {
	if data == nil {
		return nil
	}

	cfg := Auth0Config{}
	if err := mapstructure.Decode(data, &cfg); err != nil {
		return err
	}
	p.siteConfigs[site] = &cfg
	p.enabled = true
	return nil
}

func (p *Plugin) TerraformRenderStateBackend(site string) (string, error) {
	return "", nil
}

func (p *Plugin) TerraformRenderProviders(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	result := fmt.Sprintf(`
		auth0 = {
			source = "auth0/auth0"
			version = "%s"
		}`, helpers.VersionConstraint(p.provider))
	return result, nil
}

func (p *Plugin) TerraformRenderResources(site string) (string, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return "", nil
	}

	template := `
		  provider "auth0" {
			domain        = {{ .Domain|printf "%q" }}
			client_id     = {{ .ClientID|printf "%q" }}
			client_secret = {{ .ClientSecret|printf "%q" }}
		  }
	`
	return helpers.RenderGoTemplate(template, cfg)
}

func (p *Plugin) RenderTerraformComponent(site string, component string) (*schema.ComponentSchema, error) {
	cfg := p.getSiteConfig(site)
	if cfg == nil {
		return nil, nil
	}

	result := &schema.ComponentSchema{}
	return result, nil
}

func (p *Plugin) getSiteConfig(site string) *Auth0Config {
	if p.globalConfig == nil {
		return nil
	}
	cfg, ok := p.siteConfigs[site]
	if !ok {
		cfg = &Auth0Config{}
	}
	return cfg.extendConfig(p.globalConfig)
}
