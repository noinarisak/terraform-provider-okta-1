package okta

import (
	"fmt"
	"github.com/articulate/oktasdk-go/okta"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceIdentityProviders() *schema.Resource {
	return &schema.Resource{
		Create:        resourceIdentityProviderCreate,
		Read:          resourceIdentityProviderRead,
		Update:        resourceIdentityProviderUpdate,
		Delete:        resourceIdentityProviderDelete,
		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error { return nil },

		Schema: map[string]*schema.Schema{
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"GOOGLE"}, false),
				Description:  "Identity Provider Type: GOOGLE",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the Identity Provider Resource",
			},
			"protocol_type": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "OAUTH2",
				Description: "IDP Protocol type to use - ie. OAUTH2",
			},
			"protocol_scopes": &schema.Schema{
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Scopes provided to the Idp, e.g. 'openid', 'email', 'profile'",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OAUTH2 client ID",
			},
			"client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "OAUTH2 client secret",
			},
		},
	}
}

func resourceIdentityProviderCreate(d *schema.ResourceData, m interface{}) error {

	client := m.(*Config).oktaClient
	idp := client.IdentityProviders.IdentityProvider()

	idpClient := &okta.IdpClient{}

	credentials := &okta.Credentials{Client: idpClient}
	protocol := &okta.Protocol{Credentials: credentials}

	idpGroups := &okta.IdpGroups{Action:"NONE"}
	deprovisioned := &okta.Deprovisioned{Action:"NONE"}
	suspended := &okta.Suspended{Action:"NONE"}	
	
	conditions := &okta.Conditions{
		Deprovisioned: deprovisioned,
		Suspended: suspended,
	}

	provisioning := &okta.Provisioning{
		Action: "AUTO",
		ProfileMaster: true,
		Groups: idpGroups,
		Conditions: conditions,
	}

	accountLink := &okta.AccountLink{
		Action: "AUTO",
	}

	userNameTemplate := &okta.UserNameTemplate{
		Template: "idpuser.firstName",
	}
	
	subject := &okta.Subject{
		UserNameTemplate: userNameTemplate,
		MatchType: "USERNAME",
	}
	
	policy := &okta.IdpPolicy{
		Provisioning: provisioning,
		AccountLink: accountLink,
		Subject: subject,
		MaxClockSkew: 0,
	}

	idp.Type = d.Get("type").(string)
	idp.Name = d.Get("name").(string)

	protocol.Type = d.Get("protocol_type").(string)

	if len(d.Get("protocol_scopes").([]interface{})) > 0 {
		scopes := make([]string, 0)
		for _, vals := range d.Get("protocol_scopes").([]interface{}) {
			scopes = append(scopes, vals.(string))
		}
		protocol.Scopes = scopes
	}

	protocol.Credentials.Client.ClientID     = d.Get("client_id").(string)
	protocol.Credentials.Client.ClientSecret = d.Get("client_secret").(string)

	idp.Protocol = protocol
	idp.Policy   = policy

	_, _, err := client.IdentityProviders.CreateIdentityProvider(idp)
	if err != nil {
		fmt.Println("ERRORE OMG PROTECC ME!!!")
		fmt.Println(err)
		return err
	}
	return nil
}

func resourceIdentityProviderRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIdentityProviderUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceIdentityProviderDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
