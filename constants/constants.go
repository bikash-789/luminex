package constants

const (
	Env = "ENV"
)

const (
	LocalEnv = "local"
	ProdEnv  = "prod"
)

var ConfigFileMap = map[string]string{
	LocalEnv: "config_local.yaml",
	ProdEnv:  "config_prod.yaml",
}
