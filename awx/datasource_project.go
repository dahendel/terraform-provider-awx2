package awx

import (
	"github.com/hashicorp/terraform/helper/schema"
	"gitlab.com/dhendel/awx-go"
	"strconv"
)

func dataSourceProjectObject() *schema.Resource {
	return &schema.Resource{
		Read:   dataSourceProjectObjectRead,
		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of this project",
			},
			"id": &schema.Schema {
				Type:	schema.TypeInt,
				Computed: true,
				Description: "Id of the ansible project",
			},
		},
	}
}


func dataSourceProjectObjectRead(d *schema.ResourceData, meta interface{}) error {
	awx := meta.(*awx.AWX)
	awxService := awx.ProjectService
	_, res, err := awxService.ListProjects(map[string]string{
		"name": d.Get("name").(string)})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}
	d.SetId(strconv.Itoa(res.Results[0].ID))
	d = setProjectDataSourceData(d, res.Results[0])
	return nil
}

func setProjectDataSourceData(d *schema.ResourceData, r *awx.Project) *schema.ResourceData {
	d.Set("name", r.Name)
	d.Set("id", r.ID)
	return d
}