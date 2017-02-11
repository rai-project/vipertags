package vipertags

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var testYamlFile = "config.yaml"

var testJsonFile = "config.json"

var testYaml = `test:
    hostname: "127.0.0.1"
    port: 3126
    overridden: ""
`

func TestMain(m *testing.M) {
	os.Setenv("CONF_TEST_OVERRIDDEN", "test")
	bytes := []byte(testYaml)
	err := ioutil.WriteFile(testYamlFile, bytes, 0777)
	if err != nil {
		panic("cannot create temp yaml config file needed for testing")
	}
	defer os.Remove(testYamlFile)
	m.Run()
}

func ExampleGetConfig() {
	// ./config.yaml
	// test:
	//    hostname: "127.0.0.1"
	//    port: 3126
	type StringConfig struct {
		SomeHostname string `config:"test.hostname"`
		Port         int    `config:"test.port"`
		Overriden    string `config:"test.overridden"`
	}
	type InvalidConfig struct {
		Invalid string `config:"some.invalid.config"`
	}
	type InvalidMapConfig struct{}
	c := StringConfig{}
	Setup("yaml", "CONF") // Or Setup("json")
	Fill(&c)
	fmt.Println(c.SomeHostname)
	fmt.Println(c.Port)
	fmt.Println(c.Overriden)
	i := InvalidConfig{}
	Fill(&i)
	fmt.Printf(i.Invalid)
	// Output:
	// 127.0.0.1
	// 3126
	// test
}
