package jumpserver

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gopkg.in/twindagger/httpsig.v1"
)

type Config struct {
	BaseURL       string
	Username      string
	Password      string
	Token         string
	AccessKey     string
	SecretKey     string
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

func signReq(r *http.Request, accessKey string, secretKey string) error {
	headers := []string{"(request-target)", "date"}
	gmtFmt := "Mon, 02 Jan 2006 15:04:05 GMT"
	r.Header.Add("Date", time.Now().Format(gmtFmt))
	r.Header.Add("Accept", "application/json")

	signer, err := httpsig.NewRequestSigner(accessKey, secretKey, "hmac-sha256")
	if err != nil {
		return err
	}
	return signer.SignRequest(r, headers, nil)
}

func getStringFromEnv(d *schema.ResourceData, key string, envKey string) string {
	if v, ok := d.GetOk(key); ok {
		return v.(string)
	}
	return os.Getenv(envKey)
}

func getBoolFromEnv(d *schema.ResourceData, key string, envKey string) bool {
	if v, ok := d.GetOk(key); ok {
		return v.(bool)
	}
	envVal := os.Getenv(envKey)
	return envVal == "true" || envVal == "1" || envVal == "yes"
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"base_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUMPSERVER_BASE_URL", nil),
				Description: "Jumpserver Base URL. Can also be set via environment variable JUMPSERVER_BASE_URL.",
			},
			"username": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("JUMPSERVER_USERNAME", nil),
				Description: "Jumpserver Username. Can also be set via environment variable JUMPSERVER_USERNAME.",
			},
			"password": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("JUMPSERVER_PASSWORD", nil),
				Description: "Jumpserver Password. Can also be set via environment variable JUMPSERVER_PASSWORD.",
			},
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("JUMPSERVER_ACCESS_KEY", nil),
				Description: "Jumpserver API Access Key. Can also be set via environment variable JUMPSERVER_ACCESS_KEY.",
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("JUMPSERVER_SECRET_KEY", nil),
				Description: "Jumpserver API Secret Key. Can also be set via environment variable JUMPSERVER_SECRET_KEY.",
			},
			"skip_tls_verify": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If true, skip SSL certificate validation (insecure). Can also be set via environment variable JUMPSERVER_SKIP_TLS_VERIFY.",
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
	var token string

	baseURL := getStringFromEnv(d, "base_url", "JUMPSERVER_BASE_URL")
	username := getStringFromEnv(d, "username", "JUMPSERVER_USERNAME")
	password := getStringFromEnv(d, "password", "JUMPSERVER_PASSWORD")
	accessKey := getStringFromEnv(d, "access_key", "JUMPSERVER_ACCESS_KEY")
	secretKey := getStringFromEnv(d, "secret_key", "JUMPSERVER_SECRET_KEY")
	skipTLS := getBoolFromEnv(d, "skip_tls_verify", "JUMPSERVER_SKIP_TLS_VERIFY")

	if baseURL == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Missing base URL",
			Detail:   "The Jumpserver base URL must be set either in the configuration or via the JUMPSERVER_BASE_URL environment variable.",
		})
		return nil, diags
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: skipTLS,
		},
	}
	client := &http.Client{Transport: transport}

	if accessKey == "" || secretKey == "" {
		var err error
		token, err = getToken(client, baseURL, username, password)
		if err != nil {
			return nil, diag.FromErr(err)
		}
	}

	return &Config{
		Token:         token,
		BaseURL:       baseURL,
		Username:      username,
		Password:      password,
		AccessKey:     accessKey,
		SecretKey:     secretKey,
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
