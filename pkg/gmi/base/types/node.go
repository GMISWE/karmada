/*
 @Version : 1.0
 @Author  : wangxiaokang
 @Email   : xiaokang.w@gmicloud.ai
 @Time    : 2025/05/30 10:21:50
 Desc     : node types
*/

package types

type Cpu struct {
	Cores     int32   `json:"cores,omitempty"`
	Used      float64 `json:"used,omitempty"`
	Family    string  `json:"family,omitempty"`
	ModelName string  `json:"model_name,omitempty"`
	Mhz       float64 `json:"mhz,omitempty"`
}

type Mem struct {
	Size   string  `json:"size,omitempty"`
	Active string  `json:"active,omitempty"`
	Free   string  `json:"free,omitempty"`
	Used   float64 `json:"used,omitempty"`
}

type Disk struct {
	Device string  `json:"device,omitempty"`
	Fstype string  `json:"fstype,omitempty"`
	Path   string  `json:"path,omitempty"`
	Total  string  `json:"total,omitempty"`
	Free   string  `json:"free,omitempty"`
	Used   float64 `json:"used,omitempty"`
}

type Addr struct {
	Ip   string `json:"ip,omitempty"`
	Port uint32 `json:"port,omitempty"`
}

type Conection struct {
	Family  uint32 `json:"family,omitempty"`
	Laddr   *Addr  `json:"laddr,omitempty"`
	Raddr   *Addr  `json:"raddr,omitempty"`
	Status  string `json:"status,omitempty"`
	Pid     int32  `json:"pid,omitempty"`
	Process string `json:"process,omitempty"`
}

type Interface struct {
	Name         string   `json:"name,omitempty"`
	HardwareAddr string   `json:"hardware_addr,omitempty"`
	Flags        []string `json:"flags,omitempty"`
	Addrs        []string `json:"addrs,omitempty"`
}

type Net struct {
	Connections []*Conection `json:"connections,omitempty"`
	Interfaces  []*Interface `json:"interfaces,omitempty"`
}

type Host struct {
	Hostname        string `json:"hostname,omitempty"`
	Os              string `json:"os,omitempty"`
	Platform        string `json:"platform,omitempty"`
	PlatformFamily  string `json:"platform_family,omitempty"`
	PlatformVersion string `json:"platform_version,omitempty"`
	KernalVersion   string `json:"kernal_version,omitempty"`
	KernelArch      string `json:"kernel_arch,omitempty"`

	Cpu  *Cpu   `json:"cpu,omitempty"`
	Mem  *Mem   `json:"mem,omitempty"`
	Net  *Net   `json:"net,omitempty"`
	Disk []Disk `json:"disk,omitempty"`
}

// GPUProcess GPU进程信息
type GPUProcess struct {
	PID           uint32 `json:"pid"`
	ProcessName   string `json:"process_name"`
	GPUMemoryUsed uint64 `json:"gpu_memory_used"` // MB
	ContainerID   string `json:"container_id"`
	PodName       string `json:"pod_name"`
	PodNamespace  string `json:"pod_namespace"`
	PodUID        string `json:"pod_uid"`
}

type Gpu struct {
	ID          string       `json:"id"`
	UUID        string       `json:"uuid"`
	Model       string       `json:"model"`
	Core        int64        `json:"core"`
	Temp        int64        `json:"temp"`
	Power       int64        `json:"power"`
	Fan         int64        `json:"fan"`
	Load        int64        `json:"load"`
	UseMem      string       `json:"use_mem"`
	FreeMem     string       `json:"free_mem"`
	Utilization string       `json:"utilization"`
	Cluster     string       `json:"cluster"`
	Processes   []GPUProcess `json:"processes,omitempty"`
}

type Node struct {
	UUID string `json:"uuid"` // node uuid
	Host Host   `json:"host"` // node host info
	GPUs []Gpu  `json:"gpus"` // node gpus info, if no gpus, this field is empty
}
