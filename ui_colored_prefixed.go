package cli

import (
	"github.com/fatih/color"
	"os"
)

type ColorAttributes []color.Attribute

// PrefixedUi is an implementation of Ui that prefixes messages.
type ColoredPrefixedUi struct {
	AskPrefix       string
	AskSecretPrefix string

	OutputPrefix    string
	InfoPrefix      string
	ErrorPrefix     string
	WarnPrefix      string

	AskPrefixColorAttributes [2]ColorAttributes
	AskSecretPrefixColorAttributes [2]ColorAttributes

	OutputPrefixColorAttributes [2]ColorAttributes
	InfoPrefixColorAttributes [2]ColorAttributes
	ErrorPrefixColorAttributes [2]ColorAttributes
	WarnPrefixColorAttributes [2]ColorAttributes

	Ui              Ui
}

func NewColoredPrefixedUi() *ColoredPrefixedUi {
	return &ColoredPrefixedUi{
		Ui:&BasicUi{
			Reader:      os.Stdin,
			Writer:      os.Stdout,
			ErrorWriter: os.Stdout,
		},
		AskPrefixColorAttributes:[2]ColorAttributes{
			nil,
			nil,
		},
		AskSecretPrefixColorAttributes:[2]ColorAttributes{
			nil,
			nil,
		},
		OutputPrefixColorAttributes:[2]ColorAttributes{
			nil,
			nil,
		},
		InfoPrefixColorAttributes:[2]ColorAttributes{
			nil,
			nil,
		},
		ErrorPrefixColorAttributes:[2]ColorAttributes{
			nil,
			nil,
		},
		WarnPrefixColorAttributes:[2]ColorAttributes{
			nil,
			nil,
		},
	}
}
func (u *ColoredPrefixedUi) Ask(query string) (string, error) {
	return u.Ask2(query, u.AskPrefixColorAttributes[1], u.AskPrefixColorAttributes[0])
}

func (u *ColoredPrefixedUi) AskSecret(query string) (string, error) {
	return u.AskSecret2(query, u.AskSecretPrefixColorAttributes[1], u.AskSecretPrefixColorAttributes[0])
}

func (u *ColoredPrefixedUi) Error(message string) {
	u.Error2(message, u.ErrorPrefixColorAttributes[1], u.ErrorPrefixColorAttributes[0])
}

func (u *ColoredPrefixedUi) Info(message string) {
	u.Info2(message, u.InfoPrefixColorAttributes[1], u.InfoPrefixColorAttributes[0])
}

func (u *ColoredPrefixedUi) Output(message string) {
	u.Output2(message, u.OutputPrefixColorAttributes[1], u.OutputPrefixColorAttributes[0])
}

func (u *ColoredPrefixedUi) Warn(message string) {
	u.Warn2(message, u.WarnPrefixColorAttributes[1], u.WarnPrefixColorAttributes[0])
}

func (u *ColoredPrefixedUi) Ask2(query string, attrsMsg []color.Attribute, attrsPre []color.Attribute) (string, error) {
	str1 := u.AskPrefix
	if attrsPre != nil && u.AskPrefix != "" {
		str1 = u.colorize2(u.AskPrefix, attrsPre...)
	}
	str2 := query
	if attrsMsg != nil{
		str2 = u.colorize2(query, attrsMsg...)
	}

	return u.Ui.Ask(str1+str2)
}

func (u *ColoredPrefixedUi) AskSecret2(query string, attrsMsg []color.Attribute, attrsPre []color.Attribute) (string, error) {
	str1 := u.AskSecretPrefix
	if attrsPre != nil && u.AskSecretPrefix != "" {
		str1 = u.colorize2(u.AskSecretPrefix, attrsPre...)
	}
	str2 := query
	if attrsMsg != nil{
		str2 = u.colorize2(query, attrsMsg...)
	}

	return u.Ui.AskSecret(str1+str2)
}

func (u *ColoredPrefixedUi) Error2(message string, attrsMsg []color.Attribute, attrsPre []color.Attribute) {
	str1 := u.ErrorPrefix
	if attrsPre != nil && u.ErrorPrefix != "" {
		str1 = u.colorize2(u.ErrorPrefix, attrsPre...)
	}
	str2 := message
	if attrsMsg != nil{
		str2 = u.colorize2(message, attrsMsg...)
	}
	u.Ui.Error(str1 + str2)
}

func (u *ColoredPrefixedUi) Info2(message string, attrsMsg []color.Attribute, attrsPre []color.Attribute) {
	str1 := u.InfoPrefix
	if attrsPre != nil && u.InfoPrefix != ""  {
		str1 = u.colorize2(u.InfoPrefix, attrsPre...)
	}
	str2 := message
	if attrsMsg != nil{
		str2 = u.colorize2(message, attrsMsg...)
	}
	u.Ui.Info(str1 + str2)
}

func (u *ColoredPrefixedUi) Output2(message string, attrsMsg []color.Attribute, attrsPre []color.Attribute) {
	str1 := u.OutputPrefix
	if attrsPre != nil && u.OutputPrefix != "" {
		str1 = u.colorize2(u.OutputPrefix, attrsPre...)
	}
	str2 := message
	if attrsMsg != nil{
		str2 = u.colorize2(message, attrsMsg...)
	}
	u.Ui.Output(str1 + str2)
}

func (u *ColoredPrefixedUi) Warn2(message string, attrsMsg []color.Attribute, attrsPre []color.Attribute) {
	str1 := u.WarnPrefix
	if attrsPre != nil && u.WarnPrefix != "" {
		str1 = u.colorize2(u.WarnPrefix, attrsPre...)
	}
	str2 := message
	if attrsMsg != nil{
		str2 = u.colorize2(message, attrsMsg...)
	}
	u.Ui.Warn(str1 + str2)
}

func (u *ColoredPrefixedUi) colorize2(message string, attr ...color.Attribute) string {
	return color.New(attr...).SprintFunc()(message)
}
