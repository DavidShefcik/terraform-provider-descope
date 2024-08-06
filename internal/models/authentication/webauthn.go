package authentication

import (
	"github.com/descope/terraform-provider-descope/internal/models/helpers"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/boolattr"
	"github.com/descope/terraform-provider-descope/internal/models/helpers/stringattr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var WebAuthnAttributes = map[string]schema.Attribute{
	"enabled":          boolattr.Optional(),
	"top_level_domain": stringattr.Optional(),
}

type WebAuthnModel struct {
	Enabled        types.Bool   `tfsdk:"enabled"`
	TopLevelDomain types.String `tfsdk:"top_level_domain"`
}

func (m *WebAuthnModel) Values(h *helpers.Handler) map[string]any {
	data := map[string]any{}
	boolattr.Get(m.Enabled, data, "enabled")
	stringattr.Get(m.TopLevelDomain, data, "relyingPartyId")
	return data
}

func (m *WebAuthnModel) SetValues(h *helpers.Handler, data map[string]any) {
	boolattr.Set(&m.Enabled, data, "enabled")
	stringattr.Set(&m.TopLevelDomain, data, "relyingPartyId")
}
