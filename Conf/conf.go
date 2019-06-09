package Conf

import (
	"io/ioutil"
	"fmt"
	"gopkg.in/yaml.v2"
)

//配置参数结构体
type ConfSturct struct {
	ProjectPath   string `yaml:"projectPath"`
	IsUseLR1Build bool   `yaml:"isUseLR1Build"`
	GrammarFile   string `yaml:"grammarFile"`
	ProgramFile   string `yaml:"programFile"`
	LR1TableFile  string `yaml:"LR1TableFile"`
}

/*
读取配置文件
参数:void
返回:*ConfSturct
*/
func (c *ConfSturct) GetConf() *ConfSturct {
	yamlFile, err := ioutil.ReadFile("conf.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		fmt.Println(err.Error())
	}
	return c
}
