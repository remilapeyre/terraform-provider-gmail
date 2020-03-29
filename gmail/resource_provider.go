package gmail

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"google.golang.org/api/gmail/v1"
)

// Provider returns a terraform.ResourceProvider
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"token_file": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "token.json",
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "me",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"gmail_labels": dataSourceGmailLabels(),
			"gmail_label":  dataSourceGmailLabel(),
		},
		ResourcesMap: map[string]*schema.Resource{
			"gmail_filter": resourceGmailFilter(),
			"gmail_label":  resourceGmailLabel(),
		},
		ConfigureFunc: providerConfigure,
	}
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	b, err := ioutil.ReadFile("credentials.json")
	if err != nil {
		return nil, err
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}
	tokFile := d.Get("token_file").(string)
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		return nil, err
	}
	client := config.Client(context.Background(), tok)
	srv, err := gmail.New(client)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"srv":     srv,
		"user_id": d.Get("user_id").(string),
	}, nil
}

func getService(d *schema.ResourceData, meta interface{}) (*gmail.Service, string) {
	m := meta.(map[string]interface{})

	userID := d.Get("user_id").(string)
	if userID == "" {
		userID = m["user_id"].(string)
	}
	return m["srv"].(*gmail.Service), userID
}
