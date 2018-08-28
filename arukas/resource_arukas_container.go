package arukas

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/yamamoto-febc/go-arukas"
)

func resourceArukasContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArukasContainerCreate,
		Read:   resourceArukasContainerRead,
		Update: resourceArukasContainerUpdate,
		Delete: resourceArukasContainerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		SchemaVersion: 1,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"image": {
				Type:     schema.TypeString,
				Required: true,
			},
			"instances": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				Removed:  "Use `plan` instead. This attribute will be removed in a future version",
			},
			"plan": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      arukas.PlanFree,
				ValidateFunc: validation.StringInSlice(arukas.ValidPlans, false),
			},
			"endpoint": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"ports": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "tcp",
							ValidateFunc: validation.StringInSlice([]string{"tcp", "udp"}, false),
						},
						"number": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      "80",
							ValidateFunc: validation.IntBetween(1, 65535),
						},
					},
				},
			},
			"environments": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 20,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
			"cmd": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(1, 4096),
			},
			"port_mappings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"host": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ipaddress": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"service_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"endpoint_full_hostname": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_full_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArukasContainerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*arukasClient)

	name := d.Get("name").(string)
	image := d.Get("image").(string)
	instances := d.Get("instances").(int)
	cmd := d.Get("cmd").(string)
	subdomain := d.Get("endpoint").(string)

	plan := d.Get("plan").(string)
	if plan == "" {
		plan = arukas.PlanFree
	}

	var parsedEnvs []*arukas.Env
	if rawEnvs, ok := d.GetOk("environments"); ok {
		parsedEnvs = expandEnvs(rawEnvs)
	}

	var parsedPorts []*arukas.Port
	if rawPorts, ok := d.GetOk("ports"); ok {
		parsedPorts = expandPorts(rawPorts)
	}

	// create
	app, err := client.CreateApp(&arukas.RequestParam{
		Name:        name,
		Plan:        plan,
		Image:       image,
		Instances:   int32(instances),
		Command:     cmd,
		SubDomain:   subdomain,
		Environment: parsedEnvs,
		Ports:       parsedPorts,
	})
	if err != nil {
		return err
	}

	d.SetId(app.AppID())

	// start container
	if err := client.PowerOn(app.ServiceID()); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Target:  []string{arukas.StatusRunning},
		Pending: []string{arukas.StatusStopping, arukas.StatusStopped, arukas.StatusBooting},
		Timeout: client.timeout,
		Refresh: func() (interface{}, string, error) {
			service, err := client.ReadService(app.ServiceID())
			if err != nil {
				return nil, "", err
			}

			return service, service.Status(), nil
		},
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceArukasContainerRead(d, meta)
}

func resourceArukasContainerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*arukasClient)

	app, err := client.ReadApp(d.Id())
	if err != nil {
		if _, ok := err.(arukas.ErrorNotFound); ok {
			d.SetId("")
			return nil
		}
		return err
	}
	service, err := client.ReadService(app.ServiceID())
	if err != nil {
		return err
	}

	d.Set("name", app.Name())
	d.Set("image", service.Image())
	d.Set("instances", service.Instances())
	d.Set("cmd", service.Command())
	d.Set("ports", flattenPorts(service.Ports()))
	d.Set("port_mappings", flattenPortMappings(service.PortMapping()))
	d.Set("environments", flattenEnvs(service.Environment()))

	var region string
	var plan string
	plans := strings.Split(service.PlanID(), "/")
	if len(plans) == 2 {
		region = plans[0]
		plan = plans[1]
	}
	d.Set("region", region)
	d.Set("plan", plan)

	endpoint := service.EndPoint()
	if strings.HasSuffix(endpoint, ".arukascloud.io") {
		endpoint = strings.Replace(endpoint, ".arukascloud.io", "", -1)
	}
	d.Set("endpoint", endpoint)
	d.Set("endpoint_full_hostname", endpoint)
	d.Set("endpoint_full_url", fmt.Sprintf("https://%s", endpoint))

	d.Set("service_id", app.ServiceID())
	return nil
}

func resourceArukasContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*arukasClient)

	app, err := client.ReadApp(d.Id())
	if err != nil {
		if _, ok := err.(arukas.ErrorNotFound); ok {
			d.SetId("")
			return nil
		}
		return err
	}
	serviceID := app.ServiceID()

	image := d.Get("image").(string)
	instances := d.Get("instances").(int)
	cmd := d.Get("cmd").(string)
	subdomain := d.Get("endpoint").(string)

	plan := d.Get("plan").(string)
	if plan == "" {
		plan = arukas.PlanFree
	}

	var parsedEnvs []*arukas.Env
	if rawEnvs, ok := d.GetOk("environments"); ok {
		parsedEnvs = expandEnvs(rawEnvs)
	}

	var parsedPorts []*arukas.Port
	if rawPorts, ok := d.GetOk("ports"); ok {
		parsedPorts = expandPorts(rawPorts)
	}

	_, err = client.UpdateService(serviceID, &arukas.RequestParam{
		Plan:        plan,
		Image:       image,
		Instances:   int32(instances),
		Command:     cmd,
		SubDomain:   subdomain,
		Environment: parsedEnvs,
		Ports:       parsedPorts,
	})
	if err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Target:  []string{arukas.StatusRunning},
		Pending: []string{arukas.StatusRebooting, arukas.StatusBooting},
		Timeout: client.timeout,
		Refresh: func() (interface{}, string, error) {
			service, err := client.ReadService(serviceID)
			if err != nil {
				return nil, "", err
			}

			return service, service.Status(), nil
		},
	}
	_, err = stateConf.WaitForState()
	if err != nil {
		return err
	}

	return resourceArukasContainerRead(d, meta)

}

func resourceArukasContainerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*arukasClient)

	_, err := client.ReadApp(d.Id())
	if err != nil {
		if _, ok := err.(arukas.ErrorNotFound); ok {
			d.SetId("")
			return nil
		}
		return err
	}

	if err := client.DeleteApp(d.Id()); err != nil {
		return nil
	}

	d.SetId("")
	return nil
}
