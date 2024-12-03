package Modles

import "time"

type HardSource struct {
	HostName          string  `json:"hostName"`
	CPUPercent        float64 `json:"cpu_percent"`
	DiskTotal         float64 `json:"disk_total"`
	DiskUsed          float64 `json:"disk_used"`
	DiskFree          float64 `json:"disk_free"`
	DiskUsedPercent   float64 `json:"disk_used_percent"`
	MemoryTotal       float64 `json:"memory_total"`
	MemoryUsed        float64 `json:"memory_used"`
	MemoryFree        float64 `json:"memory_free"`
	MemoryBuffered    float64 `json:"memory_buffered"`
	MemoryCached      float64 `json:"memory_cached"`
	MemoryShared      float64 `json:"memory_shared"`
	MemoryAvailable   float64 `json:"memory_available"`
	MemoryUsedPercent float64 `json:"memory_used_percent"`
	CPULoad1          float64 `json:"cpu_load_1"`
	CPULoad5          float64 `json:"cpu_load_5"`
	CPULoad15         float64 `json:"cpu_load_15"`
}
type NginxSource struct {
	HostName       string `json:"hostName"`
	IsRun          int    `json:"isRun"`
	ReTotal        int    `json:"reTotal"`
	LoginUserCount int    `json:"loginUserCount"`
	RawTotal       int    `json:"rawTotal"`
	Udptotal       int    `json:"udptotal"`
	TcpTotal       int    `json:"tcpTotal"`
	TotalTcp       int    `json:"totaltcp"`
	InetTotal      int    `json:"inetTotal"`
	FragTotal      int    `json:"fragTotal"`
	TcpEstab       int    `json:"tcpEstab"`
	TcpClosed      int    `json:"tcpClosed"`
	TcpOrphaned    int    `json:"tcpOrphaned"`
	TcpTimewait    int    `json:"tcpTimewait"`
}

type SslSource struct {
	Domain     string    `json:"domain"`
	Comment    string    `json:"comment"`
	Expiration time.Time `json:"expiration"`
	DaysLeft   int       `json:"days_left"`
	Status     string    `json:"status"`
	Resolve    bool      `json:"resolve"`
}

type ContainerResource struct {
	Namespace      string  `json:"namespace"`
	PodName        string  `json:"podName"`
	ControllerName string  `json:"controllerName"`
	Container      string  `json:"container"`
	LimitCpu       float64 `json:"limitCpu"`      // CPU 核心数
	LimitMemory    int64   `json:"limitMemory"`   // 内存字节数
	RequestCpu     float64 `json:"requestCpu"`    // CPU 核心数
	RequestMemory  int64   `json:"requestMemory"` // 内存字节数
	UseCpu         float64 `json:"useCpu"`        // 使用的 CPU 核心数
	UseMemory      int64   `json:"useMemory"`     // 使用的内存字节数
	RestartCount   int     `json:"restartCount"`  // 重启次数
}
type ControllerResource struct {
	Namespace      string `json:"namespace"`
	ControllerName string `json:"controller_name"`
	Container      string `json:"container"`
	ReplicaCount   int32  `json:"replica"`
}

// MetricWithTimestamp 用于存储指标值和最后更新的时间戳
type MetricWithTimestamp struct {
	Value      float64
	MetricName string
	Timestamp  time.Time
}
type HeartSource struct {
	IsActive int    `json:"isActive"`
	Project  string `json:"project"`
	Hostname string `json:"hostname"`
}