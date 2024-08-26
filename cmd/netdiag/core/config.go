package core

type NodeConfig struct {
	Enode        string
	Ip           string
	Key          string
	Zone         string
	InstanceName string
}
type Config struct {
	Nodes        []NodeConfig
	GcpProjectId string
	GcpUsername  string
}
