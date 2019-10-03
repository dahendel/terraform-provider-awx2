package awx

import (
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/dhendel/awx-go"
	"strconv"
)

func dataSourceJobTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceJobTemplateRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this job template",
			},
			"id": &schema.Schema {
				Type:	schema.TypeInt,
				Computed: true,
				Description: "Id of the ansible job template",
			},
			"prompt_inventory": &schema.Schema{
				Type: schema.TypeBool,
				Computed: true,
				Description: "Requires an inventory ID be passed",
			},

		},
	}
}


func dataSourceJobTemplateRead(d *schema.ResourceData, meta interface{}) error {
	awx := meta.(*awx.AWX)
	awxService := awx.JobTemplateService
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"name": d.Get("name").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}
	d.SetId(strconv.Itoa(res.Results[0].ID))
	d = setJobTemplateDataSourceData(d, res.Results[0])
	return nil
}

func setJobTemplateDataSourceData(d *schema.ResourceData, r *awx.JobTemplate) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("id", r.ID)
	d.Set("prompt_inventory", r.AskInventoryOnLaunch)
	return d
}