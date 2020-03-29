package gmail

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/gmail/v1"
)

func dataSourceGmailLabel() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGmailLabelRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
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
	}
}

func dataSourceGmailLabelRead(d *schema.ResourceData, meta interface{}) error {
	srv, userID := getService(d, meta)
	name := d.Get("name").(string)

	call := srv.Users.Labels.List(userID)
	labels, err := call.Do()
	if err != nil {
		return fmt.Errorf("failed to get the list of labels: %v", err)
	}

	var label *gmail.Label

	for _, l := range labels.Labels {
		if l.Name != name {
			continue
		}
		if label != nil {
			return fmt.Errorf("multiple labels have the name '%s'", name)
		}
		label = l
	}

	m := parseLabel(label)
	for k, v := range m {
		if k == "id" {
			continue
		}
		if err = d.Set(k, v); err != nil {
			return fmt.Errorf("failed to set '%s': %v", k, err)
		}
	}

	d.SetId(m["id"].(string))

	return nil
}

func parseLabel(l *gmail.Label) map[string]interface{} {
	color := map[string]interface{}{
		"text":       "",
		"background": "",
	}
	if l.Color != nil {
		color["text"] = l.Color.TextColor
		color["background"] = l.Color.BackgroundColor
	}
	visibility := map[string]interface{}{
		"message_list": l.MessageListVisibility,
		"label_list":   l.LabelListVisibility,
	}
	threads := map[string]interface{}{
		"total":  int(l.ThreadsTotal),
		"unread": int(l.ThreadsUnread),
	}
	messages := map[string]interface{}{
		"total":  int(l.MessagesTotal),
		"unread": int(l.MessagesUnread),
	}
	return map[string]interface{}{
		"id":         l.Id,
		"name":       l.Name,
		"messages":   messages,
		"threads":    threads,
		"visibility": visibility,
		"color":      color,
	}
}
