/*
 * @Version : 1.0
 * @Author  : wangxiaokang
 * @Email   : xiaokang.w@gmicloud.ai
 * @Date    : 2025/05/09
 * @Desc    : nfs 存储类型实现
 */

package storage

import "context"

type NFS struct {
	Storage
	ctx context.Context
}

func (n *NFS) Mount() error {
	panic("not implemented")
}

func (n *NFS) Unmount() error {
	panic("not implemented")
}
