package jumpserver

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Config struct {
	Token         string
	BaseURL       string
	Username      string
	Password      string
	SkipTLSVerify bool
}

func (c *Config) NewHTTPClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.SkipTLSVerify,
		},
	}
	return &http.Client{Transport: transport}
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"skip_tls_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, skip SSL certificate validation (insecure).",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"jumpserver_host":             resourceHost(),
			"jumpserver_user":             resourceUser(),
			"jumpserver_asset":            resourceAsset(),
			"jumpserver_system_user":      resourceSystemUser(),
			"jumpserver_asset_permission": resourceAssetPermission(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	baseURL := d.Get("base_url").(string)
	username := d.Get("username").(string)
	password := d.Get("password").(string)
	skipTLS := d.Get("skip_tls_verify").(bool)

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipTLS,
		},
	}
	client := &http.Client{Transport: transport}

	token, err := getToken(client, baseURL, username, password)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return &Config{
		Token:         token,
		BaseURL:       baseURL,
		Username:      username,
		Password:      password,
		SkipTLSVerify: skipTLS,
	}, diags
}

func getToken(client *http.Client, baseURL, username, password string) (string, error) {
	url := baseURL + "/api/v1/authentication/auth/"
	credentials := map[string]string{
		"username": username,
		"password": password,
	}
	jsonValue, _ := json.Marshal(credentials)

	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if token, ok := result["token"].(string); ok {
		return token, nil
	}
	return "", fmt.Errorf("unable to fetch token from %s", url)
}
