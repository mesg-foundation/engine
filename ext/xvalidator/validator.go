package xvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/go-bip39"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/mesg-foundation/engine/ext/xerrors"
	"github.com/mesg-foundation/engine/ext/xnet"
	"github.com/mesg-foundation/engine/hash"
)

const (
	minPort = 1
	maxPort = 65535

	portSeparator = ":"
	envSeparator  = "="
)

// Struct validates a structure using go-playground/validator and more validation fields.
func Struct(s interface{}) error {
	var errs xerrors.Errors
	val, trans := New("")
	if err := val.Struct(s); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errs = append(errs, fmt.Errorf("%s", e.Translate(trans)))
		}
	}
	return errs.ErrorOrNil()
}

// New returns a new instance of 'validate' with more validation fields prefixed with 'prefix'.
func New(prefix string) (*validator.Validate, ut.Translator) {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validate := validator.New()

	validate.RegisterValidation("env", IsEnv)
	validate.RegisterTranslation("env", trans, func(ut ut.Translator) error {
		return ut.Add("env", "{0} must be a valid env variable name", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("portmap", IsPortMapping)
	validate.RegisterTranslation("portmap", trans, func(ut ut.Translator) error {
		return ut.Add("portmap", "{0} must be a valid port mapping (eg: 80 or 80:80)", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("domain", IsDomainName)
	validate.RegisterTranslation("domain", trans, func(ut ut.Translator) error {
		return ut.Add("domain", "{0} must respect domain-style notation (eg: author.name)", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("hash", IsHash)
	validate.RegisterTranslation("hash", trans, func(ut ut.Translator) error {
		return ut.Add("hash", "{0} must be a valid hash", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("accaddress", IsAccAddress)
	validate.RegisterTranslation("accaddress", trans, func(ut ut.Translator) error {
		return ut.Add("accaddress", "{0} must be a valid address", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("coins", IsCoins)
	validate.RegisterTranslation("coins", trans, func(ut ut.Translator) error {
		return ut.Add("coins", "{0} must be a valid coins", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("coin", IsCoin)
	validate.RegisterTranslation("coin", trans, func(ut ut.Translator) error {
		return ut.Add("coin", "{0} must be a valid coin", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("coinsPositiveZero", IsCoinsPositiveOrZero)
	validate.RegisterTranslation("coinsPositiveZero", trans, func(ut ut.Translator) error {
		return ut.Add("coinsPositiveZero", "{0} must be positive or zero coins", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("deccoins", IsDecCoins)
	validate.RegisterTranslation("deccoins", trans, func(ut ut.Translator) error {
		return ut.Add("deccoins", "{0} must be a valid deccoins", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("mnemonic", IsMnemonic)
	validate.RegisterTranslation("mnemonic", trans, func(ut ut.Translator) error {
		return ut.Add("mnemonic", "{0} must be a valid mnemonic", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("bech32accpubkey", IsBech32AccPubKey)
	validate.RegisterTranslation("bech32accpubkey", trans, func(ut ut.Translator) error {
		return ut.Add("bech32accpubkey", "{0} must be a valid bech32accpubkey", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	validate.RegisterValidation("bigint", IsBigInt)
	validate.RegisterTranslation("bigint", trans, func(ut ut.Translator) error {
		return ut.Add("bigint", "{0} must be a valid big int", false)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T(fe.Tag(), fe.Field(), prefix)
		return t
	})

	en_translations.RegisterDefaultTranslations(validate, trans)
	return validate, trans
}

// IsHash validates if given field is valid hash.
// It checks both string and slice of bytes.
func IsHash(fl validator.FieldLevel) bool {
	switch v := fl.Field(); v.Kind() {
	case reflect.String:
		_, err := hash.Decode(v.String())
		return err == nil
	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			// if it's not slice of bytes then break
			break
		}
		_, err := hash.DecodeFromBytes(v.Bytes())
		return err == nil
	}
	return false
}

// IsAccAddress validates if given field is valid cosmos account address.
func IsAccAddress(fl validator.FieldLevel) bool {
	switch v := fl.Field(); v.Kind() {
	case reflect.String:
		_, err := sdk.AccAddressFromBech32(v.String())
		return err == nil
	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			// if it's not slice of bytes then break
			break
		}
		return sdk.VerifyAddressFormat(v.Bytes()) == nil
	}
	return false
}

// IsCoins validates if given field is valid cosmos coins.
func IsCoins(fl validator.FieldLevel) bool {
	_, err := sdk.ParseCoins(fl.Field().String())
	return err == nil
}

// IsCoin validates if given field is valid cosmos coins.
func IsCoin(fl validator.FieldLevel) bool {
	_, err := sdk.ParseCoin(fl.Field().String())
	return err == nil
}

// IsDecCoins validates if given field is valid cosmos coins.
func IsDecCoins(fl validator.FieldLevel) bool {
	_, err := sdk.ParseDecCoins(fl.Field().String())
	return err == nil
}

// IsMnemonic validates if given field is valid cosmos coins.
func IsMnemonic(fl validator.FieldLevel) bool {
	return bip39.IsMnemonicValid(fl.Field().String())
}

// IsBech32AccPubKey validates if given field is valid cosmos coins.
func IsBech32AccPubKey(fl validator.FieldLevel) bool {
	_, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeAccPub, fl.Field().String())
	return err == nil
}

// IsCoinsPositiveOrZero validates if given field is valid cosmos positive or zero coins.
func IsCoinsPositiveOrZero(fl validator.FieldLevel) bool {
	coins, err := sdk.ParseCoins(fl.Field().String())
	if err != nil {
		return false
	}
	return coins.IsAllPositive() || coins.IsZero()
}

// IsDomainName validates if given field is valid domain name.
func IsDomainName(fl validator.FieldLevel) bool {
	return xnet.IsDomainName(fl.Field().String())
}

// IsPortMapping validates if given field is valid port mapping.
func IsPortMapping(fl validator.FieldLevel) bool {
	ports := strings.Split(fl.Field().String(), portSeparator)
	if len(ports) != 1 && len(ports) != 2 {
		return false
	}

	for _, port := range ports {
		i, err := strconv.Atoi(port)
		if err != nil || !(minPort <= i && i <= maxPort) {
			return false
		}
	}
	return true
}

var envNameRegexp = regexp.MustCompile("^[a-zA-Z][a-zA-Z0-9_]*$")

// IsEnv validates if given field is valid env variable declaration.
// The valid formats are:
// - ENV
// - ENV=
// - ENV=value
func IsEnv(fl validator.FieldLevel) bool {
	e := strings.Split(fl.Field().String(), envSeparator)
	return (len(e) == 1 || len(e) == 2) && envNameRegexp.MatchString(e[0])
}

// IsBigInt validates that field can be parsed as a bigint
func IsBigInt(fl validator.FieldLevel) bool {
	_, ok := sdk.NewIntFromString(fl.Field().String())
	return ok
}
