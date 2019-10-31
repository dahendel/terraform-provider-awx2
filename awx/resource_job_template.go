package awx

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/validation"
	uuid "github.com/satori/go.uuid"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	awxgo "gitlab.com/dhendel/awx-go"
)

func resourceJobTemplateObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceJobTemplateCreate,
		Read:   resourceJobTemplateRead,
		Delete: resourceJobTemplateDelete,
		Update: resourceJobTemplateUpdate,
		Importer: &schema.ResourceImporter{
			State: importJobTemplateData,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Description: "Job Template name",
				Required: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Description: "Description of job template",
				Optional: true,
				Default:  "",
			},
			// Run, Check, Scan
			"job_type": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "One of: run, check, scan",
				ValidateFunc: validation.StringInSlice([]string{
					"run",
					"check",
					"scan",
				}, true),
			},
			"inventory_id": &schema.Schema{
				Type:     schema.TypeString,
				Description: "Inventory ID to run job template against. Can't be set when ask_inventory_on_launch is also set",
				Optional: true,
				ConflictsWith: []string{"ask_inventory_on_launch"},
			},
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Description: "AWX Project ID the job template will run",
				Required: true,
			},
			"playbook": &schema.Schema{
				Type:     schema.TypeString,
				Description: "Select the project containing the playbook you want this job to execute.",
				Optional: true,
				Default:  "",
			},
			"credential_id": &schema.Schema{
				Type:     schema.TypeString,
				Description: "Credential that allows AWX to access the nodes this job will be ran against",
				Optional: true,
			},
			"extra_credential_ids": &schema.Schema{
				Type:     schema.TypeList,
				Description: "Extra AWX credential ID's needed for the job template to run properly",
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"vault_credential_id": &schema.Schema{
				Type:     schema.TypeString,
				Description: "AWX Vault Credential ID to decrypt vault variables",
				Optional: true,
			},
			"forks": &schema.Schema{
				Type:     schema.TypeInt,
				Description: "The number of parallel or simultaneous processes to use while executing the playbook. " +
					"An empty value, or a value less than 1 will use the Ansible default which is usually 5. " +
					"The default number of forks can be overwritten with a change to ansible.cfg. " +
					"Refer to the Ansible documentation for details about the configuration file.",
				Optional: true,
				Default:  0,
			},
			"limit": &schema.Schema{
				Type:     schema.TypeString,
				Description: "Provide a host pattern to further constrain the list of hosts that will be managed or affected by the playbook. " +
					"Multiple patterns are allowed. Refer to Ansible documentation for more information and examples on patterns.",
				Optional: true,
				Default:  "",
			},
			//0,1,2,3,4,5
			"verbosity": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
				Description: "Log level: 0,1,2,3,4,5",
				ValidateFunc: validation.IntInSlice([]int{0,1,2,3,4,5},
				),
			},
			"extra_vars": &schema.Schema{
				Type:     schema.TypeString,
				Description: "Pass extra command line variables to the playbook. Provide key/value pairs using either YAML or JSON. " +
					"Refer to the Ansible Tower documentation for example syntax.",
				Optional: true,
				Default:  "",
			},
			"job_tags": &schema.Schema{
				Type:     schema.TypeString,
				Description: "Tags are useful when you have a large playbook, and you want to run a specific part of a play or task. Use commas to separate multiple tags. Refer to Ansible Tower documentation for details on the usage of tags.",
				Optional: true,
				Default:  "",
			},
			"force_handlers": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Force handlers to run",
				Optional: true,
				Default:  false,
			},
			"skip_tags": &schema.Schema{
				Type:     schema.TypeString,
				Description: "Skip tags are useful when you have a large playbook, and you want to skip specific parts of a play or task. Use commas to separate multiple tags.",
				Optional: true,
				Default:  "",
			},
			"start_at_task": &schema.Schema{
				Type:     schema.TypeString,
				Description: "Task to start at in plabook",
				Optional: true,
				Default:  "",
			},
			"timeout": &schema.Schema{
				Type:     schema.TypeInt,
				Description: "The amount of time (in seconds) to run before the task is canceled. Defaults to 0 for no job timeout.",
				Optional: true,
				Default:  0,
			},
			"use_fact_cache": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"host_config_key": &schema.Schema{
				Type:     schema.TypeString,
				Description:"If enabled, use cached facts if available and store discovered facts in the cache.",
				Computed: true,
				Default:  "",
			},
			"ask_diff_mode_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Prompt for diff_mode on launch. This conflicts with `diff_mode`." +
					"If true it `diff_mode` must be passed to the job when the job is launched",
				Optional: true,
				Default:  false,
				ConflictsWith: []string{"diff_mode"},
			},
			"ask_limit_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Prompt for limit at job launch. This conflicts with `limit`. " +
					"If true it `limit` must be passed to the job when the job is launched",
				Optional: true,
				Default:  false,
			},
			"ask_tags_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Prompt for tags at job launch. This conflicts with `tags`. " +
					"If true it `limit` must be passed to the job when the job is launched",
				Optional: true,
				Default:  false,
				ConflictsWith: []string{"tags"},
			},
			"ask_verbosity_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Prompt for verbosity at job launch. This conflicts with `verbosity`. " +
					"If true it `verbosity` must be passed to the job when the job is launched",
				Optional: true,
				Default:  false,
				ConflictsWith: []string{"verbosity"},
			},
			"ask_inventory_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Prompt for inventory_id at job launch. This conflicts with `inventory_id`. " +
					"If true it `inventory_id` must be passed to the job when the job is launched",
				Optional: true,
				Default:  false,
				ConflictsWith: []string{"inventory_id"},
			},
			"ask_variables_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Prompt for extra_vars at job launch. This conflicts with `extra_vars`. " +
					"If true it `extra_vars` must be passed to the job when the job is launched",
				Optional: true,
				Default:  false,
				ConflictsWith: []string{"extra_vars"},
			},
			"ask_credential_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Prompt for credential_id at job launch. This conflicts with `credential_id`. " +
					"If true it `credential_id` must be passed to the job when the job is launched",
				Optional: true,
				Default:  false,
			},
			"survey_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Survey is the equivalent of vars_prompt. When enabled, `extra_vars` must be passed to the job" +
					"when it is launched.",
				Optional: true,
				Default:  false,
			},
			"become_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "If enabled, run this playbook as an administrator",
				Optional: true,
				Default:  false,
			},
			"diff_mode": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "If enabled, textual changes made to any templated files on the host are shown in the standard output",
				Optional: true,
				Default:  false,
			},
			"ask_skip_tags_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Prompt for skip_tags at job launch. This conflicts with `skip_tags`. " +
					"If true it `skip_tags` must be passed to the job when the job is launched",
				Optional: true,
				Default:  false,
				ConflictsWith: []string{"skip_tags"},
			},
			"allow_simultaneous": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "If enabled, simultaneous runs of this job template will be allowed.",
				Optional: true,
				Default:  false,
			},
			"custom_virtualenv": &schema.Schema{
				Type:     schema.TypeString,
				Description: "Select the custom Python virtual environment for this job template to run on.",
				Optional: true,
				Default:  "",
			},
			"ask_job_type_on_launch": &schema.Schema{
				Type:     schema.TypeBool,
				Description: "Prompt for job_type at job launch. This conflicts with `job_type`. " +
					"If true it `job_type` must be passed to the job when the job is launched",
				Optional: true,
				Default:  false,
				ConflictsWith: []string{"job_type"},
			},
			"job_id": {
				Type: schema.TypeInt,
				Description: "This is the ID of the job template, generally this is computed." +
					"If passed the job template would be updated.",
				Optional: true,
				Default: 0,
			},
			"allow_callbacks": {
				Type: schema.TypeBool,
				Description: "Enables call backs for the job template. This option is not available when ask_inventory_on_launch is true",
				Optional: true,
				Default: false,
				ConflictsWith: []string{"ask_inventory_on_launch"},
			},
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},
	}
}

func resourceJobTemplateCreate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService
	var jobID int
	var finished time.Time
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"name":    d.Get("name").(string),
		"project": d.Get("project_id").(string)},
	)

	if err != nil {
		return err
	}

	if len(res.Results) >= 1 {
		return fmt.Errorf("JobTemplate with name %s already exists",
			d.Get("name").(string))
	}
	_, prj, err := awx.ProjectService.ListProjects(map[string]string{
		"id": d.Get("project_id").(string)},
	)
	if err != nil {
		return err
	}
	if prj.Results[0].SummaryFields.CurrentJob["id"] != nil {
		jobID = int(prj.Results[0].SummaryFields.CurrentJob["id"].(float64))
	} else if prj.Results[0].SummaryFields.LastJob["id"] != nil {
		jobID = int(prj.Results[0].SummaryFields.LastJob["id"].(float64))
	}

	if jobID != 0 {
		// check if finished is 0
		for finished.IsZero() {
			prj, _ := awx.ProjectUpdatesService.ProjectUpdateGet(jobID)
			if prj != nil {
				finished = prj.Finished
				time.Sleep(1 * time.Second)
			}
		}
	}

	payload := map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"job_type":                 d.Get("job_type").(string),
		"inventory":                AtoipOr(d.Get("inventory_id").(string), nil),
		"project":                  AtoipOr(d.Get("project_id").(string), nil),
		"playbook":                 d.Get("playbook").(string),
		"forks":                    d.Get("forks").(int),
		"limit":                    d.Get("limit").(string),
		"verbosity":                d.Get("verbosity").(int),
		"extra_vars":               d.Get("extra_vars").(string),
		"job_tags":                 d.Get("job_tags").(string),
		"force_handlers":           d.Get("force_handlers").(bool),
		"skip_tags":                d.Get("skip_tags").(string),
		"start_at_task":            d.Get("start_at_task").(string),
		"timeout":                  d.Get("timeout").(int),
		"use_fact_cache":           d.Get("use_fact_cache").(bool),
		"host_config_key":          d.Get("host_config_key").(string),
		"ask_diff_mode_on_launch":  d.Get("ask_diff_mode_on_launch").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"ask_tags_on_launch":       d.Get("ask_tags_on_launch").(bool),
		"ask_skip_tags_on_launch":  d.Get("ask_skip_tags_on_launch").(bool),
		"ask_job_type_on_launch":   d.Get("ask_job_type_on_launch").(bool),
		"ask_verbosity_on_launch":  d.Get("ask_verbosity_on_launch").(bool),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_credential_on_launch": d.Get("ask_credential_on_launch").(bool),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"become_enabled":           d.Get("become_enabled").(bool),
		"diff_mode":                d.Get("diff_mode").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"custom_virtualenv":        AtoipOr(d.Get("custom_virtualenv").(string), nil),
		"credential":               AtoipOr(d.Get("credential_id").(string), nil),
		"vault_credential":         AtoipOr(d.Get("vault_credential_id").(string), nil),
	}

	result, err := awxService.CreateJobTemplate(payload, map[string]string{})
	if err != nil {
		return err
	}

	if creds, ok := d.GetOk("extra_credential_ids"); ok {
		for _, c := range creds.([]interface{}) {
			_, err := awxService.AddJobTemplateCredential(result.ID, c.(int))
			if err != nil {
				return err
			}
		}

	}

	d.SetId(strconv.Itoa(result.ID))
	return resourceJobTemplateRead(d, m)
}

func resourceJobTemplateUpdate(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"id":      d.Id(),
		"project": d.Get("project_id").(string)},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return fmt.Errorf("JobTemplate with name %s doesn't exists",
			d.Get("name").(string))
	}
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	result, err := awxService.UpdateJobTemplate(id, map[string]interface{}{
		"name":                     d.Get("name").(string),
		"description":              d.Get("description").(string),
		"job_type":                 d.Get("job_type").(string),
		"inventory":                AtoipOr(d.Get("inventory_id").(string), nil),
		"project":                  AtoipOr(d.Get("project_id").(string), nil),
		"playbook":                 d.Get("playbook").(string),
		"forks":                    d.Get("forks").(int),
		"limit":                    d.Get("limit").(string),
		"verbosity":                d.Get("verbosity").(int),
		"extra_vars":               d.Get("extra_vars").(string),
		"job_tags":                 d.Get("job_tags").(string),
		"force_handlers":           d.Get("force_handlers").(bool),
		"skip_tags":                d.Get("skip_tags").(string),
		"start_at_task":            d.Get("start_at_task").(string),
		"timeout":                  d.Get("timeout").(int),
		"use_fact_cache":           d.Get("use_fact_cache").(bool),
		"host_config_key":          d.Get("host_config_key").(string),
		"ask_diff_mode_on_launch":  d.Get("ask_diff_mode_on_launch").(bool),
		"ask_variables_on_launch":  d.Get("ask_variables_on_launch").(bool),
		"ask_limit_on_launch":      d.Get("ask_limit_on_launch").(bool),
		"ask_tags_on_launch":       d.Get("ask_tags_on_launch").(bool),
		"ask_skip_tags_on_launch":  d.Get("ask_skip_tags_on_launch").(bool),
		"ask_job_type_on_launch":   d.Get("ask_job_type_on_launch").(bool),
		"ask_verbosity_on_launch":  d.Get("ask_verbosity_on_launch").(bool),
		"ask_inventory_on_launch":  d.Get("ask_inventory_on_launch").(bool),
		"ask_credential_on_launch": d.Get("ask_credential_on_launch").(bool),
		"survey_enabled":           d.Get("survey_enabled").(bool),
		"become_enabled":           d.Get("become_enabled").(bool),
		"diff_mode":                d.Get("diff_mode").(bool),
		"allow_simultaneous":       d.Get("allow_simultaneous").(bool),
		"custom_virtualenv":        AtoipOr(d.Get("custom_virtualenv").(string), nil),
		"credential":               AtoipOr(d.Get("credential_id").(string), nil),
		"vault_credential":         AtoipOr(d.Get("vault_credential_id").(string), nil),
	}, map[string]string{})
	if err != nil {
		return err
	}

	if creds, ok := d.GetOk("extra_credential_ids"); ok {
		for _, c := range creds.([]interface{}) {
			_, err := awxService.AddJobTemplateCredential(result.ID, c.(int))
			if err != nil {
				return err
			}
		}

	}

	return resourceJobTemplateRead(d, m)
}

func resourceJobTemplateRead(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"id": strconv.Itoa(d.Get("job_id").(int)),
		//"name":    d.Get("name").(string),
		//"project": d.Get("project_id").(string),
	})
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		return nil
	}
	d = setJobTemplateResourceData(d, res.Results[0])
	return nil
}

func resourceJobTemplateDelete(d *schema.ResourceData, m interface{}) error {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService
	_, res, err := awxService.ListJobTemplates(map[string]string{
		"id":      d.Id(),
		"project": d.Get("project_id").(string)},
	)
	if err != nil {
		return err
	}
	if len(res.Results) == 0 {
		d.SetId("")
		return nil
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	_, err = awxService.DeleteJobTemplate(id)
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func setJobTemplateResourceData(d *schema.ResourceData, r *awxgo.JobTemplate) *schema.ResourceData {
	d.Set("job_id", r.ID)
	d.Set("allow_simultaneous", r.AllowSimultaneous)
	d.Set("ask_credential_on_launch", r.AskCredentialOnLaunch)
	d.Set("ask_job_type_on_launch", r.AskJobTypeOnLaunch)
	d.Set("ask_limit_on_launch", r.AskLimitOnLaunch)
	d.Set("ask_skip_tags_on_launch", r.AskSkipTagsOnLaunch)
	d.Set("ask_tags_on_launch", r.AskTagsOnLaunch)
	d.Set("ask_variables_on_launch", r.AskVariablesOnLaunch)
	d.Set("credential_id", r.Credential)
	d.Set("description", r.Description)
	d.Set("extra_vars", r.ExtraVars)
	d.Set("force_handlers", r.ForceHandlers)
	d.Set("forks", r.Forks)
	d.Set("host_config_key", r.HostConfigKey)
	d.Set("inventory_id", r.Inventory)
	d.Set("job_tags", r.JobTags)
	d.Set("job_type", r.JobType)
	d.Set("diff_mode", r.DiffMode)
	d.Set("custom_virtualenv", r.CustomVirtualenv)
	d.Set("vault_credential_id", r.VaultCredential)
	d.Set("limit", r.Limit)
	d.Set("name", r.Name)
	d.Set("become_enabled", r.BecomeEnabled)
	d.Set("use_fact_cache", r.UseFactCache)
	d.Set("playbook", r.Playbook)
	d.Set("project_id", r.Project)
	d.Set("skip_tags", r.SkipTags)
	d.Set("start_at_task", r.StartAtTask)
	d.Set("survey_enabled", r.SurveyEnabled)
	d.Set("verbosity", r.Verbosity)
	extraIDs := getExtraIDs(r)
	d.Set("extra_credential_ids", extraIDs)
	return d
}

func importJobTemplateData(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	awx := m.(*awxgo.AWX)
	awxService := awx.JobTemplateService

	id, err :=strconv.Atoi(d.Id())

	job, err := awxService.GetJobTemplate(id)

	if err != nil {
		return nil, err
	}

	if job == nil {
		return nil, err
	}

	resources := []*schema.ResourceData{setJobTemplateResourceData(d, job)}

	return resources, nil
}

func getExtraIDs(template *awxgo.JobTemplate) []int {
	creds := template.SummaryFields.ExtraCredentials
	var ids []int
	for _, c := range creds {
		if cred, ok := c.(*awxgo.Credential); ok {
			ids = append(ids, cred.ID)
		}
	}

	return ids
}

func genHostConfigKey() string {
	return uuid.NewV4().String()
}