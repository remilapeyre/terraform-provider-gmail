package gmail

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"google.golang.org/api/gmail/v1"
)

func resourceGmailFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceGmailFilterCreate,
		Read:   resourceGmailFilterRead,
		Delete: resourceGmailFilterDelete,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"from": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"to": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"subject": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"query": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"negated_query": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"has_attachment": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			"exclude_chats": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
				Default:  false,
			},
			"size": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Optional: true,
				Default:  0,
			},
			"size_comparison": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"add_labels": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"remove_labels": {
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"forward": {
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
		},
	}
}

func resourceGmailFilterCreate(d *schema.ResourceData, meta interface{}) error {
	srv, userID := getService(d, meta)

	addLabels := make([]string, d.Get("add_labels").(*schema.Set).Len())
	for i, l := range d.Get("add_labels").(*schema.Set).List() {
		addLabels[i] = l.(string)
	}

	removeLabels := make([]string, d.Get("remove_labels").(*schema.Set).Len())
	for i, l := range d.Get("remove_labels").(*schema.Set).List() {
		removeLabels[i] = l.(string)
	}

	filter := &gmail.Filter{
		Action: &gmail.FilterAction{
			Forward:        d.Get("forward").(string),
			AddLabelIds:    addLabels,
			RemoveLabelIds: removeLabels,
		},
		Criteria: &gmail.FilterCriteria{
			ExcludeChats:   d.Get("exclude_chats").(bool),
			From:           d.Get("from").(string),
			HasAttachment:  d.Get("has_attachment").(bool),
			NegatedQuery:   d.Get("negated_query").(string),
			Query:          d.Get("query").(string),
			Size:           int64(d.Get("size").(int)),
			SizeComparison: d.Get("size_comparison").(string),
			Subject:        d.Get("subject").(string),
			To:             d.Get("to").(string),
		},
	}
	call := srv.Users.Settings.Filters.Create(userID, filter)

	filter, err := call.Do()
	if err != nil {
		return fmt.Errorf("failed to create filter: %v", err)
	}

	d.SetId(filter.Id)

	return resourceGmailFilterRead(d, meta)
}

func resourceGmailFilterRead(d *schema.ResourceData, meta interface{}) error {
	srv, userID := getService(d, meta)
	filterID := d.Id()

	call := srv.Users.Settings.Filters.Get(userID, filterID)
	filter, err := call.Do()
	if err != nil {
		if strings.Contains(err.Error(), "Error 404: Not Found") {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("failed to read filter '%s': %v", filterID, err)
	}

	if err = d.Set("from", filter.Criteria.From); err != nil {
		return fmt.Errorf("failed to set 'from': %v", err)
	}
	if err = d.Set("to", filter.Criteria.To); err != nil {
		return fmt.Errorf("failed to set 'to': %v", err)
	}
	if err = d.Set("subject", filter.Criteria.Subject); err != nil {
		return fmt.Errorf("failed to set 'subject': %v", err)
	}
	if err = d.Set("query", filter.Criteria.Query); err != nil {
		return fmt.Errorf("failed to set 'query': %v", err)
	}
	if err = d.Set("negated_query", filter.Criteria.NegatedQuery); err != nil {
		return fmt.Errorf("failed to set 'negated_query': %v", err)
	}
	if err = d.Set("has_attachment", filter.Criteria.HasAttachment); err != nil {
		return fmt.Errorf("failed to set 'has_attachment': %v", err)
	}
	if err = d.Set("exclude_chats", filter.Criteria.ExcludeChats); err != nil {
		return fmt.Errorf("failed to set 'exclude_chats': %v", err)
	}
	if err = d.Set("size", filter.Criteria.Size); err != nil {
		return fmt.Errorf("failed to set 'size': %v", err)
	}
	if err = d.Set("size_comparison", filter.Criteria.SizeComparison); err != nil {
		return fmt.Errorf("failed to set 'size_comparison': %v", err)
	}
	if err = d.Set("add_labels", filter.Action.AddLabelIds); err != nil {
		return fmt.Errorf("failed to set 'add_labels': %v", err)
	}
	if err = d.Set("remove_labels", filter.Action.RemoveLabelIds); err != nil {
		return fmt.Errorf("failed to set 'remove_labels': %v", err)
	}
	if err = d.Set("forward", filter.Action.Forward); err != nil {
		return fmt.Errorf("failed to set 'forward': %v", err)
	}

	return nil
}

func resourceGmailFilterDelete(d *schema.ResourceData, meta interface{}) error {
	srv, userID := getService(d, meta)
	filterID := d.Id()

	call := srv.Users.Settings.Filters.Delete(userID, filterID)
	err := call.Do()
	if err != nil {
		return fmt.Errorf("failed to delete filter '%s': %v", filterID, err)
	}

	d.SetId("")

	return nil
}
