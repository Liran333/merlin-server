package permission

// config:
// permissions
type Config struct {
	Permissions []PermObject `json:"permissions"`
}

type PermObject struct {
	ObjectType string `json:"object_type"`
	Rules      []Rule `json:"rules"`
}

type Rule struct {
	Role      string   `json:"role"`
	Operation []string `json:"operation"`
}
