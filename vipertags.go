package vipertags

import (
	"path/filepath"
	"reflect"
	"strings"

	"os"

	"github.com/Sirupsen/logrus"
	"github.com/fatih/structs"
	"github.com/k0kubun/pp"
	"github.com/spf13/viper"
)

func buildConfiguration(st0 interface{}, prefix string) interface{} {
	st := structs.New(st0)
	for _, field := range st.Fields() {
		defaultTagValue := field.Tag("default")
		envTagValue := field.Tag("env")
		configTagValue := field.Tag("config")

		if configTagValue != "" {
			prefix = prefix + configTagValue
		}

		if field.Kind() == reflect.Struct {
			buildConfiguration(field.Value(), prefix)
			continue
		}
		if field.Kind() == reflect.Map {
			pp.Println("map == ", field)
			continue
		}
		if field.Kind() == reflect.Array || field.Kind() == reflect.Slice {
			logrus.Fatal("Not currently working...")
			t := reflect.Indirect(reflect.ValueOf(field.Value()))

			slice := reflect.MakeSlice(t.Type(), t.Len(), t.Len())
			slices := reflect.New(slice.Type())
			slices.Elem().Set(slice)

			for ii := 0; ii < t.Len(); ii++ {
				e := t.Index(ii).Interface()
				elem := buildConfiguration(e, prefix)
				slices.Elem().Index(ii).Set(reflect.ValueOf(elem))
			}
			pp.Println(slices.Interface())
			field.Set(slices.Interface())
			continue
		}
		if defaultTagValue != "" && configTagValue != "" {
			viper.SetDefault(configTagValue, defaultTagValue)
		}
		if defaultTagValue != "" && configTagValue == "" {
			field.Set(defaultTagValue)
		}

		if envTagValue != "" && configTagValue != "" {
			viper.BindEnv(configTagValue, envTagValue)
		}
		if envTagValue != "" && configTagValue == "" {
			if e := os.Getenv(envTagValue); e != "" {
				field.Set(e)
			}
		}
		if configTagValue != "" {
			field.Set(viper.Get(configTagValue))
		}
	}
	return st0
}

func Fill(class interface{}) {
	buildConfiguration(class, "")
}

func Setup(fileType string, prefix string) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("conf")
	viper.AddConfigPath("config")
	viper.SetConfigType(fileType)
	viper.AutomaticEnv()
	viper.SetEnvPrefix(prefix)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}

func FromFile(filename string, prefix string) {
	Setup(strings.Replace(filepath.Ext(filename), ".", "", 1), prefix)
	viper.SetConfigFile(filename)
}
