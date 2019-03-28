package config

import (
	"fmt"

	toml "github.com/pelletier/go-toml"
)

var (
	Conf = New()
)

/**
 * 返回单例实例
 * @method New
 */
func New() (config *toml.Tree) {
	config, err := toml.LoadFile("./config/config.toml")

	if err != nil {
		fmt.Println("TomlError ", err.Error())
	}
	return
}
