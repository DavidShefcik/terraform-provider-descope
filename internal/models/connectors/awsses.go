package connectors

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var AWSSESAttributes = map[string]schema.Attribute{
	"id":          stringattr.IdentifierMatched(),
	"name":        stringattr.Required(stringattr.StandardLenValidator),
	"description": stringattr.Default(""),

	"access_key_id":  stringattr.SecretRequired(),
	"secret":         stringattr.SecretRequired(),
	"endpoint":       stringattr.Default(""),
	"region":         stringattr.Required(),
	"sender_address": stringattr.Required(),
	"sender_name":    stringattr.Default(""),
}

// Model

type AWSSESModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`

	AccessKeyID   types.String `tfsdk:"access_key_id"`
	Secret        types.String `tfsdk:"secret"`
	Endpoint      types.String `tfsdk:"endpoint"`
	Region        types.String `tfsdk:"region"`
	SenderAddress types.String `tfsdk:"sender_address"`
	SenderName    types.String `tfsdk:"sender_name"`
}

func (m *AWSSESModel) Values(h *helpers.Handler) map[string]any {
	data := connectorValues(m.ID, m.Name, m.Description, h)
	data["type"] = "ses"
	data["configuration"] = m.ConfigurationValues(h)
	return data
}

func (m *AWSSESModel) SetValues(h *helpers.Handler, data map[string]any) {
	// all connector values are specified in the schema
}

// Configuration

func (m *AWSSESModel) ConfigurationValues(h *helpers.Handler) map[string]any {
	c := map[string]any{}
	stringattr.Get(m.AccessKeyID, c, "accessKeyId")
	stringattr.Get(m.Secret, c, "secret")
	stringattr.Get(m.Endpoint, c, "endpoint")
	stringattr.Get(m.Region, c, "region")
	stringattr.Get(m.SenderAddress, c, "senderAddress")
	stringattr.Get(m.SenderName, c, "senderName")
	return c
}

// Matching

func (m *AWSSESModel) GetName() types.String {
	return m.Name
}

func (m *AWSSESModel) GetID() types.String {
	return m.ID
}

func (m *AWSSESModel) SetID(id types.String) {
	m.ID = id
}
