package jumpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceSystemUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSystemUserCreate,
		ReadContext:   resourceSystemUserRead,
		UpdateContext: resourceSystemUserUpdate,
		DeleteContext: resourceSystemUserDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "common",
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ssh",
			},
			"login_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "auto",
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  81,
			},
			"sudo": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/bin/whoami",
			},
			"shell": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/bin/bash",
			},
			"sftp_root": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "tmp",
			},
			"home": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "/home/student",
			},
			"username_same_with_user": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"auto_push": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"su_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceSystemUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	// Prepare payload
	payload := map[string]interface{}{
		"name":                    d.Get("name").(string),
		"username":                d.Get("username").(string),
		"password":                d.Get("password").(string),
		"type":                    d.Get("type").(string),
		"protocol":                d.Get("protocol").(string),
		"login_mode":              d.Get("login_mode").(string),
		"priority":                d.Get("priority").(int),
		"sudo":                    d.Get("sudo").(string),
		"shell":                   d.Get("shell").(string),
		"sftp_root":               d.Get("sftp_root").(string),
		"home":                    d.Get("home").(string),
		"username_same_with_user": d.Get("username_same_with_user").(bool),
		"auto_push":               d.Get("auto_push").(bool),
		"su_enabled":              d.Get("su_enabled").(bool),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	// Send POST request to create system user
	url := fmt.Sprintf("%s/api/v1/assets/system-users/", c.BaseURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return diag.FromErr(err)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	} else {
		if err := signReq(req, c.AccessKey, c.SecretKey); err != nil {
			return diag.FromErr(err)
		}
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Failed to create system user. Status code: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	// Set resource ID
	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Failed to retrieve ID for created system user")
	}

	return diags
}

func resourceSystemUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	id := d.Id()
	url := fmt.Sprintf("%s/api/v1/assets/system-users/%s/", c.BaseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	} else {
		if err := signReq(req, c.AccessKey, c.SecretKey); err != nil {
			return diag.FromErr(err)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error fetching system user: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", result["name"].(string))
	d.Set("username", result["username"].(string))
	d.Set("type", result["type"].(string))
	d.Set("protocol", result["protocol"].(string))
	d.Set("login_mode", result["login_mode"].(string))

	// Handle numeric values with proper type conversion
	if priority, ok := result["priority"].(float64); ok {
		d.Set("priority", int(priority))
	} else {
		return diag.Errorf("Failed to parse 'priority' field from API response")
	}

	// Additional fields to set if available
	if home, ok := result["home"].(string); ok {
		d.Set("home", home)
	}
	if shell, ok := result["shell"].(string); ok {
		d.Set("shell", shell)
	}

	d.SetId(id)

	return diags
}

func resourceSystemUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	// Prepare payload for update
	payload := map[string]interface{}{
		"name":                    d.Get("name").(string),
		"username":                d.Get("username").(string),
		"password":                d.Get("password").(string),
		"type":                    d.Get("type").(string),
		"protocol":                d.Get("protocol").(string),
		"login_mode":              d.Get("login_mode").(string),
		"priority":                d.Get("priority").(int),
		"sudo":                    d.Get("sudo").(string),
		"shell":                   d.Get("shell").(string),
		"sftp_root":               d.Get("sftp_root").(string),
		"home":                    d.Get("home").(string),
		"username_same_with_user": d.Get("username_same_with_user").(bool),
		"auto_push":               d.Get("auto_push").(bool),
		"su_enabled":              d.Get("su_enabled").(bool),
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	id := d.Id()
	url := fmt.Sprintf("%s/api/v1/assets/system-users/%s/", c.BaseURL, id)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return diag.FromErr(err)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	} else {
		if err := signReq(req, c.AccessKey, c.SecretKey); err != nil {
			return diag.FromErr(err)
		}
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to update system user. Status code: %d", resp.StatusCode)
	}

	// Update Terraform state after successful update
	resourceSystemUserRead(ctx, d, m)

	return diags
}

func resourceSystemUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	id := d.Id()
	url := fmt.Sprintf("%s/api/v1/assets/system-users/%s/", c.BaseURL, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	} else {
		if err := signReq(req, c.AccessKey, c.SecretKey); err != nil {
			return diag.FromErr(err)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Failed to delete system user. Status code: %d", resp.StatusCode)
	}

	d.SetId("") // Mark resource as deleted

	return diags
}
