package ws

// AgentStatus 监控进程信息
type AgentStatus struct {
	ID uint // agent表ID
}

// StatusServer 监控服务器状态
type StatusServer struct {
	BootTime    uint64
	Uptime      uint64
	OS          string
	CPUPercent  float64
	DiskPercent float64
	MemPercent  float64
	SwapPercent float64
}

// StatusProcess 节点进程状态
type StatusProcess struct {
	NodeID     uint // chaindeploynode表ID
	PID        int32
	BootTime   int64
	CPUPercent float64
	MemPercent float32
}

// StatusRPC 节点进程RPC状态
type StatusRPC struct {
	NodeID uint // chaindeploynode表ID
	RPCID  uint // chainstatus表ID
	Value  string
}

// StatusCmd 节点进程命令状态
type StatusCmd struct {
	Cmd   string
	Value string `json:"value,omitempty"`
	Error string `json:"error,omitempty"`
}

// Node 节点进程信息
type Node struct {
	ID        uint
	Name      string
	ChainName string
	PID       int32 `json:"pid,omitempty"`
}
