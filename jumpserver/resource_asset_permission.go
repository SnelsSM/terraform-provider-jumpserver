package jumpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
)

func resourceAssetPermission() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetPermissionCreate,
		ReadContext:   resourceAssetPermissionRead,
		UpdateContext: resourceAssetPermissionUpdate,
		DeleteContext: resourceAssetPermissionDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"users_display": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"assets_display": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"system_users_display": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceAssetPermissionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	permission := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"is_active":            d.Get("is_active").(bool),
		"users_display":        d.Get("users_display").([]interface{}),
		"assets_display":       d.Get("assets_display").([]interface{}),
		"system_users_display": d.Get("system_users_display").([]interface{}),
	}

	url := c.BaseURL + "/api/v1/perms/asset-permissions/"
	jsonValue, _ := json.Marshal(permission)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Check for 201 Created status code
	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating asset permission: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving asset permission ID after creation, response: %v", result)
	}

	resourceAssetPermissionRead(ctx, d, m)

	return diags
}

func resourceAssetPermissionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	id := d.Id()
	url := fmt.Sprintf("%s/api/v1/perms/asset-permissions/%s/", c.BaseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Check for 200 OK status code
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error fetching asset permission: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	// Update resource data with fetched values
	d.Set("name", result["name"].(string))
	d.Set("is_active", result["is_active"].(bool))
	d.Set("users_display", result["users_display"].([]interface{}))
	d.Set("assets_display", result["assets_display"].([]interface{}))
	d.Set("system_users_display", result["system_users_display"].([]interface{}))

	return diags
}

func resourceAssetPermissionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	permission := map[string]interface{}{
		"name":                 d.Get("name").(string),
		"is_active":            d.Get("is_active").(bool),
		"users_display":        d.Get("users_display").([]interface{}),
		"assets_display":       d.Get("assets_display").([]interface{}),
		"system_users_display": d.Get("system_users_display").([]interface{}),
	}

	id := d.Id()
	url := fmt.Sprintf("%s/api/v1/perms/asset-permissions/%s/", c.BaseURL, id)
	jsonValue, _ := json.Marshal(permission)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonValue))
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Check for 200 OK status code
	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error updating asset permission: %s", resp.Status)
	}

	resourceAssetPermissionRead(ctx, d, m)

	return diags
}

func resourceAssetPermissionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	id := d.Id()
	url := fmt.Sprintf("%s/api/v1/perms/asset-permissions/%s/", c.BaseURL, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return diag.FromErr(err)
	}
	req.Header.Set("Authorization", "Bearer "+c.Token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Check for 204 No Content status code
	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting asset permission: %s", resp.Status)
	}

	d.SetId("") // Mark resource as destroyed
	return diags
}
