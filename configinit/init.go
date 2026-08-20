package configinit

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/go-playground/validator/v10"
	"github.com/mcuadros/go-defaults"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

func InitConfig[T any](path string) (*T, error) {
	var conf T
	defaults.SetDefaults(&conf) // This sets the default values

	v := viper.NewWithOptions(viper.ExperimentalBindStruct())
	v.AutomaticEnv()

	if path != "" {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); err != nil {
			return &conf, fmt.Errorf("error reading .env file: %w", err)
		}
	} else {
		// Viper currently have an experimental feature called BindStruct that
		// basically does the same. We can use it once it's in a stable release.
		bindEnvs(conf)
	}

	err := viper.Unmarshal(&conf, viper.DecodeHook(
		mapstructure.ComposeDecodeHookFunc(
			keyDecodeHook,
			bigIntDecodeHook,
			addressDecodeHook,
			mapstructure.StringToTimeDurationHookFunc(),
		),
	))
	if err != nil {
		return &conf, fmt.Errorf("error unmarshalling config: %w", err)
	}

	return &conf, validator.New().Struct(conf)
}

func bindEnvs(conf interface{}) {
	val := reflect.TypeOf(conf)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		envTag := field.Tag.Get("mapstructure") // Get the mapstructure tag

		if envTag != "" {
			_ = viper.BindEnv(envTag)
		}
	}
}

func bigIntDecodeHook(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	if t != reflect.TypeOf(&big.Int{}) {
		return data, nil
	}

	// An empty string decodes to the zero value (nil), leaving required-ness to
	// the validator.
	if data.(string) == "" {
		return (*big.Int)(nil), nil
	}

	target, ok := new(big.Int).SetString(data.(string), 10)
	if !ok {
		return data, errors.New("wrong data format for field")
	}

	return target, nil
}

func addressDecodeHook(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	if t != reflect.TypeOf(common.Address{}) {
		return data, nil
	}

	// An empty string decodes to the zero address, leaving required-ness to the
	// validator.
	if data.(string) == "" {
		return common.Address{}, nil
	}

	if !common.IsHexAddress(data.(string)) {
		return data, errors.New("wrong data format for field")
	}

	return common.HexToAddress(data.(string)), nil
}

func keyDecodeHook(f reflect.Type, t reflect.Type, data interface{}) (interface{}, error) {
	if f.Kind() != reflect.String {
		return data, nil
	}

	if t != reflect.TypeOf(&ecdsa.PrivateKey{}) {
		return data, nil
	}

	privateKey, err := crypto.HexToECDSA(data.(string))
	if err != nil {
		return data, errors.New("wrong data format for field")
	}

	return privateKey, nil
}
