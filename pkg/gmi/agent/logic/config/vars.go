/*
 @Version : 1.0
 @Author  : steven.wong
 @Email   : 'wangxk1991@gamil.com'
 @Time    : 2024/01/23 09:37:29
 Desc     :
*/

package config

import (
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

type CtxKey string

const (
	SUCCESS           = 0
	ERR_API_AIRGAP    = 100
	ERR_API_USER      = 101
	ERR_API_INFERENCE = 102
	ERR_API_INFER     = 103

	ERR_REQ_KUBE      = 110
	ERR_REQ_MYSQL     = 111
	ERR_HTTP_REQ      = 112

	ERR_UNAUTHORIZED  = 401
	ERR_NO_PERMISSION = 402
	ERR_RELOGIN       = 403
)

var (
	BASE_DIR = getProjectAbPath()
)

const (
	COMMIT_SHA    = ""
	DB_NOT_DELETE = 0
	UUID_BASE_KEY = "UUID_BASE_KEY"
	CIPHER_KEY    = "uOvKLmVfztaXGpNYd4Z0I1SiT7MweJhl"
)

const (
	CtxUidKey CtxKey = "uid"
)

const (
	REDIS_KEY_BASE       = "ifs:"
	REDIS_KEY_SESSION    = REDIS_KEY_BASE + "session:"
	REDIS_KEY_MODELSTOPO = REDIS_KEY_BASE + "modelstopo:"
	REDIS_KEY_PERMISSION = REDIS_KEY_BASE + "permission"
)

func GetSha() string {
	return COMMIT_SHA
}

func getProjectAbPath() string {
	currentPath := getCurrentAbPath()
	if strings.Contains(currentPath, "/pkg/apis") || strings.Contains(currentPath, "/config") {
		return filepath.Dir(filepath.Dir(currentPath))
	}
	return currentPath
}

// 最终方案-全兼容
func getCurrentAbPath() string {
	dir := getCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return getCurrentAbPathByCaller()
	}
	return dir
}

// 获取当前执行文件绝对路径
func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

// 获取当前执行文件绝对路径（go run）
func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}
