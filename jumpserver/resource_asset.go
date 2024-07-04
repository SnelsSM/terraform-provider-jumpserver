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

func resourceAsset() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetCreate,
		ReadContext:   resourceAssetRead,
		UpdateContext: resourceAssetUpdate,
		DeleteContext: resourceAssetDelete,

		Schema: map[string]*schema.Schema{
			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocols": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"nodes_display": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceAssetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	asset := map[string]interface{}{
		"hostname":      d.Get("hostname").(string),
		"ip":            d.Get("ip").(string),
		"platform":      d.Get("platform").(string),
		"protocols":     d.Get("protocols").([]interface{}),
		"nodes_display": d.Get("nodes_display").([]interface{}),
	}

	url := c.BaseURL + "/api/v1/assets/assets/"
	jsonValue, _ := json.Marshal(asset)

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
		return diag.Errorf("Error creating asset: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving asset ID after creation, response: %v", result)
	}

	resourceAssetRead(ctx, d, m)

	return diags
}

func resourceAssetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	id := d.Id()
	url := fmt.Sprintf("%s/api/v1/assets/assets/%s/", c.BaseURL, id)

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
		return diag.Errorf("Error fetching asset: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	// Update resource data with fetched values
	d.Set("hostname", result["hostname"].(string))
	d.Set("ip", result["ip"].(string))
	d.Set("platform", result["platform"].(string))
	d.Set("protocols", result["protocols"].([]interface{}))
	d.Set("nodes_display", result["nodes_display"].([]interface{}))

	return diags
}

func resourceAssetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	asset := map[string]interface{}{
		"hostname":      d.Get("hostname").(string),
		"ip":            d.Get("ip").(string),
		"platform":      d.Get("platform").(string),
		"protocols":     d.Get("protocols").([]interface{}),
		"nodes_display": d.Get("nodes_display").([]interface{}),
	}

	id := d.Id()
	url := fmt.Sprintf("%s/api/v1/assets/assets/%s/", c.BaseURL, id)
	jsonValue, _ := json.Marshal(asset)

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
		return diag.Errorf("Error updating asset: %s", resp.Status)
	}

	resourceAssetRead(ctx, d, m)

	return diags
}

func resourceAssetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	id := d.Id()
	url := fmt.Sprintf("%s/api/v1/assets/assets/%s/", c.BaseURL, id)

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
		return diag.Errorf("Error deleting asset: %s", resp.Status)
	}

	d.SetId("") // Mark resource as destroyed
	return diags
}
