// @Version : 1.0
// @Author  : steven.wong
// @Email   : 'wangxk1991@gamil.com'
// @Time    : 2024/01/19 11:27:38
// Desc     :
package config

type Agent struct {
	KarmadaControlPlaneAddr string `yaml:"karmada-control-plane-addr,omitempty"`
}

type Business struct {
	Agent *Agent `yaml:"agent,omitempty"`
}
