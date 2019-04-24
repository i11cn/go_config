// Package config 包封装了配置项的数据结构，并且支持从Json、Yaml、INI文件，以及命令行、环境变量中读取更新配置项
//
// 根据加载配置项的先后顺序，可以视为配置项的优先级，例如：先加载Yaml，再加载环境变量，最后加载命令行，则最终生效的是命令行参数，
// 未设置的命令行参数则环境变量会生效，未设置的环境变量，则Yaml配置文件会生效。
//
package config

// TODO: 后续增加配置项变更后的触发功能

type (
	Config interface {
		// Add 给指定的Key上增加一个值，如果Key原来已有对应的值，则扩展成数组存放。返回值为本Config对象，可以级联使用，例如:
		//
		// config.Add("value", "key").Add("value2", "Key")
		Add(value interface{}, path string, mpath ...string) Config

		// Set 给指定的Key上设置一个值，如果Key原来已有对应的值，则原有数据会被丢弃，设置为新值。返回值为本Config对象，可以级联使用
		Set(value interface{}, path string, mpath ...string) Config

		// Delete 删除指定的Key，无论其下还有什么数据，均会被删除，包括子配置项。返回值为本Config对象，可以级联使用
		Delete(path string, mpath ...string) Config

		// LoadYaml 将输入数据作为Yaml格式解析，并且加载到本配置对象中，之前的数据会被清除
		LoadYaml(in []byte) (Config, error)

		// LoadYamlFile 读取Yaml配置文件到本配置对象中，之前的数据会被清除
		LoadYamlFile(file string) (Config, error)

		// LoadJson 将输入数据作为Json格式解析，并且加载到本配置对象中，之前的数据会被清除
		LoadJson(in []byte) (Config, error)

		// LoadJsonFile 读取Json配置文件到本配置对象中，之前的数据会被清除
		LoadJsonFile(file string) (Config, error)

		// LoadIni 将输入数据作为INI格式解析，并且加载到本配置对象中，之前的数据会被清除，由于INI文件中的Key不能重复，通常对于数组，
		// 都是在Key之后增加序号来实现的，因此此处提供一个函数参数kp来处理这一类配置项，默认截取所有Key之后的数字，
		// 如果需要自己处理，则需要自己传入key_preprocess
		LoadIni(in []byte, kp ...func(string) string) (Config, error)

		// LoadIniFile 读取INI配置文件到本配置对象中，之前的数据会被清除
		LoadIniFile(file string, kp ...func(string) string) (Config, error)

		// ToYaml 将本配置对象的内容导出为Yaml格式
		ToYaml() string

		// ToJson 将本配置对象的内容导出为Json格式
		ToJson() string

		// ToIni 将本配置对象的内容导出为INI格式
		ToIni() string

		// Get 以指定类型获取数据，要求必须为对应类型，类型不匹配则会返回错误
		Get(v interface{}, path string, mpath ...string) error

		// Convert 以指定类型获取数据，尽可能的做类型转换的尝试，包括数值类型之间的转换，以及各种类型和字符串类型之间的转换
		Convert(v interface{}, path string, mpath ...string) error

		// Keys 返回本配置对象中的所有配置项的名称
		Keys() []string

		// AddCommandFlag 从命令行中加载指定名称的参数，以Add的方式保存到本配置对象中
		AddCommandFlag(name string) Config

		// AddEnv 从环境变量中加载指定名称的配置到本对象中，由于环境变量不能重复，因此如果数组类型，就需要有一定规则分割，可以通过参数delimiter指定分隔符
		AddEnv(name string, delimiter ...string) Config

		// Clear 清除本对象中所有配置项
		Clear() Config
	}
)

// NewConfig 创建一个新的Config对象，并且做适当的初始化，由于Config对象中所有成员都是私有的，因此必须依靠该函数来初始化
func NewConfig() Config {
	return &config_impl{}
}
