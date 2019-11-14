package awx

import (
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	awxgo "gitlab.com/dhendel/awx-go"
)

func resourceInventoryGroupObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceInventoryGroupCreate,
		Read:   resourceInventoryGroupRead,
		Update: resourceInventoryGroupUpdate,
		Delete: resourceInventoryGroupDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"inventory_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"variables": &schema.Schema{
				Type:      schema.TypeString,
				Optional:  true,
				Default:   "",
				StateFunc: normalizeJSONYaml,
			},
			"child_group_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

func resourceInventoryGroupCreate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.GroupService

	_, res, _ := awxService.ListGroups(map[string]string{
		"name":      d.Get("name").(string),
		"inventory": d.Get("inventory_id").(string),
	})
	if len(res.Results) >= 1 {
		return fmt.Errorf("InventoryGroup %s with id %d already exists", res.Results[0].Name, res.Results[0].ID)
	}

	result, err := awxService.CreateGroup(map[string]interface{}{
		"name":        d.Get("name").(string),
		"description": d.Get("description").(string),
		"inventory":   d.Get("inventory_id").(string),
		"variables":   d.Get("variables").(string),
	}, map[string]string{})
	if err != nil {
		return err
	}

	if childGroups, ok := d.GetOkExists("child_group_ids"); ok {
		for _, i := range childGroups.([]interface{}) {
			_, err := awxService.AddChildGroup(result.ID, i.(int))

			if err != nil {
				return fmt.Errorf("Failed to add child group %d to group id %d ", i, result.ID)
			}
		}
	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceInventoryGroupRead(d, m)

}

func resourceInventoryGroupUpdate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.GroupService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	_, res, _ := awxService.ListGroups(map[string]string{"id": d.Id()})
	if len(res.Results) >= 1 {

		_, err = awxService.UpdateGroup(id, map[string]interface{}{
			"name":        d.Get("name").(string),
			"description": d.Get("description").(string),
			"inventory":   d.Get("inventory_id").(string),
			"variables":   d.Get("variables").(string),
		}, nil)
		if err != nil {
			return err
		}

		return resourceInventoryGroupRead(d, m)
	}
	return fmt.Errorf("Group %s with id %d doesn't exist", d.Get("name").(string), id)

}

func resourceInventoryGroupDelete(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.GroupService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("InventoryGroup %d not found", id)
	}
	if _, err := awxService.DeleteGroup(id); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceInventoryGroupRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.GroupService
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return fmt.Errorf("InventoryGroup %d not found", id)
	}
	_, res, err := awxService.ListGroups(map[string]string{
		"name":      d.Get("name").(string),
		"inventory": d.Get("inventory_id").(string),
	})
	if err != nil {
		return err
	}
	d = setInventoryGroupResourceData(d, res.Results[0])
	return nil
}

func setInventoryGroupResourceData(d *schema.ResourceData, r *awxgo.Group) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("description", r.Description)
	d.Set("inventory_id", r.Inventory)
	d.Set("variables", normalizeJSONYaml(r.Variables))
	return d
}
