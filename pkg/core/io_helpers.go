package core

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"io"
)

func UnmarshalYaml(filename string, obj interface{}) error {
	v, err := ParseYaml(filename)
	if err != nil {
		return err
	}
	err = v.Unmarshal(obj)
	if err != nil {
		return err
	}
	return nil
}

func MustUnmarshalYaml(filename string, obj interface{}) {
	err := UnmarshalYaml(filename, obj)
	if err != nil {
		panic(err)
	}
}

func ParseYaml(filename string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(filename)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}

func Obj2Proto(o interface{}, m proto.Message) error {
	bs, err := json.Marshal(o)
	if err != nil {
		return err
	}

	err = protojson.Unmarshal(bs, m)
	if err != nil {
		return err
	}

	return nil
}

type YamlRefresher interface {
	RefreshFile(*viper.Viper) error
}

func ParseYamlFile(filename string, watch bool, fn YamlRefresher) error {
	v, err := ParseYaml(filename)
	if err != nil {
		return err
	}

	return YamlParserV(v, watch, fn)
}

func ParseYamlFolder(path string, watch bool, fn YamlRefresher) error {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("master")
	if err := v.ReadInConfig(); err != nil {
		return err
	}
	return YamlParserV(v, watch, fn)
}

func YamlParserV(v *viper.Viper, watch bool, p YamlRefresher) error {
	err := p.RefreshFile(v)
	if err != nil {
		return err
	}

	if watch {
		v.WatchConfig()
		v.OnConfigChange(func(in fsnotify.Event) {
			filename := v.ConfigFileUsed()
			zap.S().Infow("config file modified", "filename", filename)
			if (in.Op & (fsnotify.Write | fsnotify.Create)) == 0 {
				return
			}
			zap.S().Infow("config parsing file", "filename", filename)

			err := p.RefreshFile(v)
			if err != nil {
				zap.S().Errorw("failed to refresh config",
					"filename", filename,
					"error", err)
			}
		})
	}

	return nil
}

func ZapGlobalLevel(l zapcore.Level) {
	c := zap.NewDevelopmentConfig()
	c.Level = zap.NewAtomicLevelAt(l)
	log, err := c.Build()
	if err != nil {
		panic(err)
	}

	zap.ReplaceGlobals(log)
}

func CloseLogErr(c io.Closer) {
	err := c.Close()
	if err != nil {
		zap.S().Error("Close failed",
			"item", c,
			"err", err)
	}
}

type ArrayFlags struct {
	data []string
}

func (i *ArrayFlags) Get() []string {
	return i.data
}

func (i *ArrayFlags) String() string {
	return fmt.Sprintf("%v", i.data)
}

func (i *ArrayFlags) Set(value string) error {
	i.data = append(i.data, value)
	return nil
}

func MultiStringFlag(name string, usage string) *ArrayFlags {
	var myFlags ArrayFlags
	flag.Var(&myFlags, name, usage)
	return &myFlags
}
