package privilege

type Config struct {
	Npu     PrivilegeConfig `json:"npu"`
	Disable PrivilegeConfig `json:"disable"`
}

type PrivilegeConfig struct {
	Orgs []OrgIndex `json:"orgs"`
}

type OrgIndex struct {
	OrgId   string `json:"org_id"`
	OrgName string `json:"org_name"`
}
