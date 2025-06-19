package jumpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceUserCreate,
		ReadContext:   resourceUserRead,
		UpdateContext: resourceUserUpdate,
		DeleteContext: resourceUserDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_active": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_roles": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
		},
	}
}

func resourceUserCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	user := map[string]interface{}{
		"name":         d.Get("name").(string),
		"username":     d.Get("username").(string),
		"email":        d.Get("email").(string),
		"is_active":    d.Get("is_active").(bool),
		"system_roles": d.Get("system_roles").([]interface{}),
	}

	url := c.BaseURL + "/api/v1/users/users/"
	jsonValue, _ := json.Marshal(user)

	// Log request body
	log.Printf("Request Body: %s\n", string(jsonValue))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
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

	// Check for 201 Created status code
	if resp.StatusCode != http.StatusCreated {
		return diag.Errorf("Error creating user: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return diag.FromErr(err)
	}

	// Log the entire response
	log.Printf("Response Body: %v\n", result)

	if id, ok := result["id"].(string); ok {
		d.SetId(id)
	} else {
		return diag.Errorf("Error retrieving user ID after creation, response: %v", result)
	}

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	url := fmt.Sprintf("%s/api/v1/users/users/%s/", c.BaseURL, d.Id())
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
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return diag.FromErr(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			d.SetId("")
			return nil
		}
		return diag.Errorf("Error reading user: %s", resp.Status)
	}

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", user["name"].(string))
	d.Set("username", user["username"].(string))
	d.Set("email", user["email"].(string))
	d.Set("is_active", user["is_active"].(bool))
	d.Set("system_roles", user["system_roles"].([]interface{}))

	return diags
}

func resourceUserUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	user := map[string]interface{}{
		"name":         d.Get("name").(string),
		"username":     d.Get("username").(string),
		"email":        d.Get("email").(string),
		"is_active":    d.Get("is_active").(bool),
		"system_roles": d.Get("system_roles").([]interface{}),
	}

	url := fmt.Sprintf("%s/api/v1/users/users/%s/", c.BaseURL, d.Id())
	jsonValue, _ := json.Marshal(user)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonValue))
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
		return diag.Errorf("Error updating user: %s", resp.Status)
	}

	resourceUserRead(ctx, d, m)

	return diags
}

func resourceUserDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*Config)

	var diags diag.Diagnostics

	id := d.Id()
	url := fmt.Sprintf("%s/api/v1/users/users/%s/", c.BaseURL, id)

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

	// Check for 204 No Content status code
	if resp.StatusCode != http.StatusNoContent {
		return diag.Errorf("Error deleting user: %s", resp.Status)
	}

	d.SetId("") // Mark resource as destroyed
	return diags
}
