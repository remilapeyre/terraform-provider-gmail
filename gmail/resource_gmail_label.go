package gmail

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/gmail/v1"
)

func resourceGmailLabel() *schema.Resource {
	return &schema.Resource{
		Create: resourceGmailLabelCreate,
		Read:   resourceGmailLabelRead,
		Update: resourceGmailLabelUpdate,
		Delete: resourceGmailLabelDelete,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"visibility": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"color": {
				Type:         schema.TypeMap,
				Required:     true,
				ValidateFunc: validateColor,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			// Out parameters
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
		},
	}
}

func validateColor(i interface{}, k string) (s []string, es []error) {
	colors := []string{"#000000", "#434343", "#666666", "#999999", "#cccccc",
		"#ffffff", "#fb4c2f", "#ffad47", "#fad165", "#16a766", "#43d692", "#4a86e8",
		"#a479e2", "#f691b3", "#f6c5be", "#ffe6c7", "#fef1d1", "#b9e4d0", "#c6f3de",
		"#c9daf8", "#e4d7f5", "#fcdee8", "#efa093", "#ffd6a2", "#fce8b3", "#89d3b2",
		"#a0eac9", "#a4c2f4", "#d0bcf1", "#fbc8d9", "#e66550", "#ffbc6b", "#fcda83",
		"#44b984", "#68dfa9", "#6d9eeb", "#b694e8", "#f7a7c0", "#cc3a21", "#eaa041",
		"#f2c960", "#149e60", "#3dc789", "#3c78d8", "#8e63ce", "#e07798", "#ac2b16",
		"#cf8933", "#d5ae49", "#0b804b", "#2a9c68", "#285bac", "#653e9b", "#b65775",
		"#822111", "#a46a21", "#aa8831", "#076239", "#1a764d", "#1c4587", "#41236d",
		"#83334c", "#464646", "#e7e7e7", "#0d3472", "#b6cff5", "#0d3b44", "#98d7e4",
		"#3d188e", "#e3d7ff", "#711a36", "#fbd3e0", "#8a1c0a", "#f2b2a8", "#7a2e0b",
		"#ffc8af", "#7a4706", "#ffdeb5", "#594c05", "#fbe983", "#684e07", "#fdedc1",
		"#0b4f30", "#b3efd3", "#04502e", "#a2dcc1", "#c2c2c2", "#4986e7", "#2da2bb",
		"#b99aff", "#994a64", "#f691b2", "#ff7537", "#ffad46", "#662e37", "#ebdbde",
		"#cca6ac", "#094228", "#42d692", "#16a765", "#efefef", "#f3f3f3", "",
	}
	for key, value := range i.(map[string]interface{}) {
		found := false
		for _, c := range colors {
			if value.(string) == c {
				found = true
				break
			}
		}
		if !found {
			es = append(es, fmt.Errorf("expected color.%s to be one of %v, got %s", key, colors, value))
		}
	}
	return
}

func resourceGmailLabelCreate(d *schema.ResourceData, meta interface{}) error {
	srv, userID := getService(d, meta)
	label := getLabel(d)
	call := srv.Users.Labels.Create(userID, label)
	label, err := call.Do()
	if err != nil {
		return fmt.Errorf("failed to create label: %v", err)
	}

	d.SetId(label.Id)

	return resourceGmailLabelRead(d, meta)
}

func resourceGmailLabelUpdate(d *schema.ResourceData, meta interface{}) error {
	srv, userID := getService(d, meta)
	label := getLabel(d)
	call := srv.Users.Labels.Update(userID, d.Id(), label)
	_, err := call.Do()
	if err != nil {
		return fmt.Errorf("failed to update label: %v", err)
	}
	return nil
}

func resourceGmailLabelRead(d *schema.ResourceData, meta interface{}) error {
	srv, userID := getService(d, meta)
	labelID := d.Id()

	call := srv.Users.Labels.Get(userID, labelID)
	label, err := call.Do()
	if err != nil {
		if strings.Contains(err.Error(), "Error 404: Not Found") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read label '%s': %v", labelID, err)
	}

	if err = d.Set("name", label.Name); err != nil {
		return fmt.Errorf("failed to set 'name': %v", err)
	}
	visibility := map[string]interface{}{
		"message_list": label.MessageListVisibility,
		"label_list":   label.LabelListVisibility,
	}
	if err = d.Set("visibility", visibility); err != nil {
		return fmt.Errorf("failed to set 'visibility': %v", err)
	}
	color := map[string]interface{}{
		"text":       "",
		"background": "",
	}
	if label.Color != nil {
		color["text"] = label.Color.TextColor
		color["background"] = label.Color.BackgroundColor
	}
	if err = d.Set("color", color); err != nil {
		return fmt.Errorf("failed to set 'color': %v", err)
	}
	messages := map[string]interface{}{
		"total":  label.MessagesTotal,
		"unread": label.MessagesUnread,
	}
	if err = d.Set("messages", messages); err != nil {
		return fmt.Errorf("failed to set 'messages': %v", err)
	}
	threads := map[string]interface{}{
		"total":  label.ThreadsTotal,
		"unread": label.ThreadsUnread,
	}
	if err = d.Set("threads", threads); err != nil {
		return fmt.Errorf("failed to set 'threads': %v", err)
	}

	return nil
}

func resourceGmailLabelDelete(d *schema.ResourceData, meta interface{}) error {
	srv, userID := getService(d, meta)
	labelID := d.Id()

	call := srv.Users.Labels.Delete(userID, labelID)
	err := call.Do()
	if err != nil {
		return fmt.Errorf("failed to delete label '%s': %v", labelID, err)
	}

	d.SetId("")

	return nil
}

func getLabel(d *schema.ResourceData) *gmail.Label {
	color := d.Get("color").(map[string]interface{})
	c := &gmail.LabelColor{}
	if backgroundColor := color["background"]; backgroundColor != nil {
		c.BackgroundColor = backgroundColor.(string)
	}
	if textColor := color["text"]; textColor != nil {
		c.TextColor = textColor.(string)
	}

	label := &gmail.Label{
		Color: c,
		Name:  d.Get("name").(string),
	}

	visibility := d.Get("visibility").(map[string]interface{})
	if messageList := visibility["message_list"]; messageList != nil {
		label.MessageListVisibility = messageList.(string)
	}
	if labelList := visibility["label_list"]; labelList != nil {
		label.LabelListVisibility = labelList.(string)
	}
	return label
}
