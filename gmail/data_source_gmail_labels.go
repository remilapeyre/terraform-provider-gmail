package gmail

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceGmailLabels() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGmailLabelsRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"labels": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"messages": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"threads": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"visibility": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"message_list": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"label_list": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"color": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceGmailLabelsRead(d *schema.ResourceData, meta interface{}) error {
	srv, userID := getService(d, meta)

	call := srv.Users.Labels.List(userID)
	labels, err := call.Do()
	if err != nil {
		return fmt.Errorf("failed to get the list of labels: %v", err)
	}

	res := make([]interface{}, len(labels.Labels))

	for i, l := range labels.Labels {
		m := parseLabel(l)
		res[i] = m
	}

	d.SetId("labels")

	if err = d.Set("labels", res); err != nil {
		return fmt.Errorf("failed to set labels: %v", err)
	}

	return nil
}
