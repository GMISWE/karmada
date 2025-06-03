package models

import (
	"database/sql"
	"time"
)

func Models() []interface{} {
	return []interface{}{
		&User{},
		&Image{},
		&Product{},
		&Inference{},
	}
}

type InferenceStatus string
type InferTaskStatus string
type Role string

const (
	InferenceStatusReject InferenceStatus = "reject" // 驳回
	InferenceStatusReady  InferenceStatus = "ready"  // 待审核
	InferenceStatusDeploy InferenceStatus = "deploy" // 部署中
)
const (
	InferTaskStatusOnline    InferTaskStatus = "online"
	InferTaskStatusOffline   InferTaskStatus = "offline"
	InferTaskStatusUnhealthy InferTaskStatus = "unhealthy"
	InferTaskStatusBusy      InferTaskStatus = "busy"
	InferTaskStatusWait      InferTaskStatus = "wait"
)

const (
	RoleAdmin  Role = "admin"
	RoleDevops Role = "devops"
	RoleAlgo   Role = "algo"
	RoleGuest  Role = "guest"
)

type (
	User struct {
		Uid       string         `gorm:"primaryKey;size:16"  json:"uid"`
		Username  string         `gorm:"unique;index:idx_user_name;not null;size:256"  json:"username"`
		Password  string         `gorm:"not null;size:256"  json:"password"`
		Email     sql.NullString `gorm:"not null;size:256"  json:"email"`
		Phone     sql.NullString `gorm:"not null;size:16"  json:"phone"`
		Role      Role           `gorm:"not null;size:256"  json:"role"`
		Namespace string         `gorm:"null;size:256"  json:"namespace"` // Kubernetes namespace
		CreateAt  time.Time      `gorm:"not null;autoUpdateTime"  json:"create_at"`
		UpdateAt  time.Time      `gorm:"not null;autoUpdateTime"  json:"update_at"`
		IsDelete  int64          `gorm:"not null;size:1"  json:"is_delete"` // 删除标志
	}
	Task struct {
		TaskID        uint   `gorm:"primaryKey"  json:"task_id"`
		TaskName      string `gorm:"unique;index:idx_task_name;not null;size:256"  json:"task_name"`
		TaskImage     string `gorm:"not null;type:text"  json:"task_image"`
		TaskCapacity  uint32 `gorm:"not null" json:"task_capacity"`
		TaskGPU       uint32 `gorm:"not null" json:"task_gpu"`
		TaskParameter string `gorm:"not null" json:"task_parameter"`
	}
	Image struct {
		Iid         string         `gorm:"primaryKey;size:16" json:"iid"`
		Uid         string         `gorm:"not null;size:16" json:"uid"`
		Share       int64          `gorm:"not null;size:1" json:"share"` // 是否分享
		Name        string         `gorm:"unique;index:idx_image_name;not null;size:256" json:"name"`
		Url         string         `gorm:"not null;size:256" json:"url"`
		Versions    string         `gorm:"not null;type:text" json:"versions"`
		Description sql.NullString `gorm:"not null" json:"description"`
		Params      string         `gorm:"params" json:"params"`
		CreateAt    time.Time      `gorm:"not null;autoUpdateTime" json:"create_at"`
		UpdateAt    time.Time      `gorm:"not null;autoUpdateTime" json:"update_at"`
		IsDelete    uint8          `gorm:"not null;size:1" json:"is_delete"` // 删除标志
	}
	Product struct {
		Pid         string         `gorm:"primaryKey;size:16" json:"pid"`
		Name        string         `gorm:"unique;index:idx_product_name;not null;size:256" json:"name"`
		Description sql.NullString `gorm:"not null" json:"description"`
		CreateAt    time.Time      `gorm:"not null;autoUpdateTime" json:"create_at"`
		UpdateAt    time.Time      `gorm:"not null;autoUpdateTime" json:"update_at"`
		IsDelete    uint8          `gorm:"not null;size:1" json:"is_delete"` // 删除标志
	}
	Inference struct {
		Iid string `gorm:"primaryKey;size:16" json:"iid"`
		Uid string `gorm:"not null;size:16" json:"uid"`
		Pid string `gorm:"not null;size:16" json:"pid"`

		TaskID         uint32 `gorm:"size:16" json:"task_id"`
		TaskName       string `gorm:"unique;index:idx_task_name;not null;size:256" json:"task_name"`
		TaskImages     string `gorm:"not null;type:text" json:"task_images"`
		TaskResources  string `gorm:"not null;type:text" json:"task_resources"`
		TaskParameters string `gorm:"not null;type:text" json:"task_parameters"`

		Name        string          `gorm:"unique;index:idx_image_name;not null;size:256" json:"name"`
		Description sql.NullString  `gorm:"not null" json:"description"`
		Version     string          `gorm:"not null;size:256" json:"version"`
		InferStatus InferenceStatus `gorm:"type:enum('ready', 'reject', 'deploy')" json:"infer_status"`
		TaskStatus  InferTaskStatus `gorm:"type:enum('unhealthy', 'online', 'offline', 'busy', 'wait')" json:"task_status"`
		CreateAt    time.Time       `gorm:"not null;autoUpdateTime" json:"create_at"`
		UpdateAt    time.Time       `gorm:"not null;autoUpdateTime" json:"update_at"`
		IsDelete    uint8           `gorm:"not null;size:1" json:"is_delete"` // 删除标志
	}
	// // 集群定义
	// Icluster struct {
	// 	Rid        uint      `gorm:"primaryKey" json:"rid"`
	// 	Name       string    `gorm:"unique;not null" json:"name"`
	// 	LinkConfig string    `gorm:"not null;type:text" json:"link_config"`
	// 	Status     int       `gorm:"not null;type:enum('init','active','deactive')" json:"status"`
	// 	CreateAt   time.Time `gorm:"not null;autoUpdateTime" json:"create_at"`
	// 	UpdateAt   time.Time `gorm:"not null;autoUpdateTime" json:"update_at"`
	// 	IsDelete   uint8     `gorm:"not null;size:1" json:"is_delete"` // 删除标志
	// }
)
