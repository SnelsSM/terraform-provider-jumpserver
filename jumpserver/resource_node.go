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

func resourceNode() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNodeCreate,
		ReadContext:   resourceNodeRead,
		UpdateContext: resourceNodeUpdate,
		DeleteContext: resourceNodeDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the node.",
			},
			"parent": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the parent node.",
			},
		},
	}
}

func resourceNodeCreate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*Config)
	var diags diag.Diagnostics

	// Get parent node name from config
	parentName := d.Get("parent").(string)

	// Find parent node ID
	parentID, err := findNodeIDByName(c, parentName)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error finding parent node '%s': %v", parentName, err))
	}

	// Prepare payload for node creation
	nodeName := d.Get("name").(string)
	payload := map[string]string{
		"value": nodeName,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return diag.FromErr(err)
	}

	// Send POST request to create child node
	url := fmt.Sprintf("%s/api/v1/assets/nodes/%s/children/", c.BaseURL, parentID)
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

	client := c.NewHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return diag.Errorf("Failed to create node. Status code: %d", resp.StatusCode)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	// Set resource ID
	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Failed to retrieve ID for created node, response: %v", result)
	}

	return diags
}

func resourceNodeRead(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	id := d.Id()
	url := fmt.Sprintf("%s/api/v1/assets/nodes/%s/", c.BaseURL, id)

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

	client := c.NewHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Node was deleted outside of Terraform
		d.SetId("")
		return diags
	}

	if resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error fetching node: %s", resp.Status)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	// Update resource data with fetched values
	if name, ok := result["value"].(string); ok {
		d.Set("name", name)
	}

	return diags
}

func resourceNodeUpdate(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	// For now, we'll treat node updates as a no-op since JumpServer API doesn't seem
	// to support changing node parent or name via a simple PUT request
	// In a real implementation, you might want to handle name changes differently

	// If parent changes, that would require special handling which is not
	// straightforward in JumpServer API

	return resourceNodeRead(ctx, d, m)
}

func resourceNodeDelete(ctx context.Context, d *schema.ResourceData, m any) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	id := d.Id()
	url := fmt.Sprintf("%s/api/v1/assets/nodes/%s/", c.BaseURL, id)

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

	client := c.NewHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	// Check for 204 No Content status code
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return diag.Errorf("Error deleting node: %s (status code: %d)", resp.Status, resp.StatusCode)
	}

	d.SetId("") // Mark resource as destroyed
	return diags
}
