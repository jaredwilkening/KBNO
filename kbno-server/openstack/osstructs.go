package openstack

type bootRes struct {
	Server bootResS `json:"server"`
}

type bootResS struct {
	ID           string  `json:"id"`
	OSdiskConfig string  `json:"OS-DCF:diskConfig"`
	Links        []linkS `json:"links"`
	AdminPass    string  `json:"adminPass"`
}

type linkS struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
}

type tokenRes struct {
	Access tokenResS `json:"access"`
}

type tokenResS struct {
	ServiceCatalog []serviceS `json:"serviceCatalog"`
	Token          tokenS     `json:"token"`
	User           userS      `json:"user"`
}

type serviceS struct {
	Endpoints      []endpointS `json:"endpoints"`
	EndpointsLinks []linkS     `json:"endpoints_links"`
	Name           string      `json:"name"`
	Type           string      `json:"type"`
}

type endpointS struct {
	AdminURL    string `json:"adminURL"`
	InternalURL string `json:"internalURL"`
	PublicURL   string `json:"publicURL"`
	Region      string `json:"region"`
}

type tokenS struct {
	Expires string  `json:"expires"`
	ID      string  `json:"id"`
	Tenant  tenantS `json:"tenant"`
}

type tenantS struct {
	Description string `json:"description"`
	Enable      bool   `json:"enable"`
	ID          string `json:"id"`
	Name        string `json:"name"`
}

type userS struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	Roles      []roleS `json:"roles"`
	RolesLinks []linkS `json:"roles_links"`
	Username   string  `json:"username"`
}

type roleS struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type statusRes struct {
	Server Server `json:"server"`
}

type statusAllRes struct {
	Servers []Server `json:"servers"`
}

type Server struct {
	ID        string      `json:"id"`
	Name      string      `json:"name"`
	HostID    string      `json:"hostId"`
	TentantID string      `json:"tentant_id"`
	UserID    string      `json:"user_id"`
	Status    string      `json:"status"`
	Updated   string      `json:"updated"`
	Metadata  interface{} `json:"metadata"`
	KeyName   string      `json:"key_name"`
	Links     []linkS     `json:"links"`
	Image     idlinksS    `json:"image"`
	Flavor    idlinksS    `json:"flavor"`
	Addresses addressesS  `json:"addresses"`
}

type idlinksS struct {
	ID    string  `json:"id"`
	Links []linkS `json:"links"`
}

type addressesS struct {
	Services []addressS `json:"services"`
}

type addressS struct {
	Addr    string `json:"addr"`
	Version int    `json:"version"`
}
