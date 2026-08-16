package forms

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
)

const (
	DATE_FORMAT             = "Y-m-d"
	TIME_FORMAT             = "H:i:S"
	DATETIME_FORMAT         = "Y-m-d H:i:S"
	DISPLAY_DATE_FORMAT     = "m/d/Y"
	DISPLAY_TIME_FORMAT     = "H:i"
	DISPLAY_DATETIME_FORMAT = "m/d/Y H:i"
)

type InputBaseProps struct {
	Label       string
	Type        string
	Name        string
	Placeholder string
	Value       string
	Required    bool
	LabelIcon   string
	IconBefore  string
	IconAfter   string
	Error       string
	HelperText  string
	ReadOnly    bool
	Disabled    bool
	Attrs       []Attr
	Selected    []string
	Options     []SelectOption
	Multiple    bool
}

type SelectOption struct {
	Label    string
	Subtitle string
	Value    string
	Disabled bool
	Attrs    []Attr
}

type Attr struct {
	Name  string
	Value string
	Hide  bool
}

type ButtonBaseProps struct {
	Type      string
	Loading   bool
	StartIcon string
	EndIcon   string
	Disabled  bool
	Label     string
	Attrs     []Attr
}

func (b *ButtonBaseProps) MakeAttrs(moreAttrs ...Attr) templ.Attributes {
	attrs := append(b.Attrs, moreAttrs...)
	return AppendAttrs(attrs, []Attr{
		{Name: "disabled", Hide: !b.Disabled},
		{Name: "type", Value: ifElse(b.Type, "button")},
	}...)
}

func (b *InputBaseProps) MakePlaceholder() string {
	if strings.Contains(b.Type, "date") && strings.Contains(b.Type, "time") {
		return getAttribute(b.Attrs, "data-altformat", "MM/DD/YYYY HH:MM")
	} else if strings.Contains(b.Type, "date") {
		return getAttribute(b.Attrs, "data-altformat", "MM/DD/YYYY")
	} else if strings.Contains(b.Type, "time") {
		return "HH:MM"
	}
	if len(b.Placeholder) > 0 {
		return b.Placeholder
	}
	return "Please enter " + b.Label
}

func (b *InputBaseProps) MakeAttrs(moreAttrs ...Attr) templ.Attributes {
	attrs := append(b.Attrs, moreAttrs...)
	isDateOrTimeField := strings.Contains(b.Type, "date") || strings.Contains(b.Type, "time")
	attrs = append(attrs, []Attr{
		{Name: "disabled", Hide: !b.Disabled},
		{Name: "readonly", Hide: !b.ReadOnly},
		{Name: "required", Hide: !b.Required},
		{Name: "autocomplete", Value: "off"},
		{Name: "name", Value: fmt.Sprintf(`%s`, b.Name)},
		{Name: "id", Value: fmt.Sprintf(`input_id_%s`, b.Name)},
		{Name: "placeholder", Value: b.MakePlaceholder()},
	}...)
	if isDateOrTimeField {
		attrs = append(attrs, []Attr{
			{Name: "type", Value: "text"},
			{Name: "data-flatpicker"},
			{Name: "data-altformat", Value: getAttribute(attrs, "data-altformat", DISPLAY_DATE_FORMAT)},
			{Name: "data-dateformat", Value: getAttribute(attrs, "data-dateformat", DATE_FORMAT)},
		}...)
	} else if strings.EqualFold(b.Type, "textarea") {
	} else if strings.EqualFold(getAttribute(attrs, "type", b.Type), "checkbox") {
		attrs = append(attrs, Attr{Name: "type", Value: "checkbox"})
	} else if strings.EqualFold(getAttribute(attrs, "type", b.Type), "radio") {
		attrs = append(attrs, Attr{Name: "type", Value: "radio"})
	} else {
		attrs = append(attrs, Attr{Name: "type", Value: ifElse(b.Type, "text")})
		attrs = append(attrs, Attr{Name: "value", Value: ifElse(b.Value, "")})
	}

	if strings.EqualFold(b.Type, "password") {
		attrs = append(attrs, Attr{
			Name: "autocomplete", Value: "new-password",
		})
	}
	if isDateOrTimeField {
		attrs = append(attrs, Attr{
			Name: "autocomplete", Value: "bday",
		})
	}
	return AppendAttrs(attrs)
}

func (b *InputBaseProps) IfTrue(show bool, val string) string {
	if show {
		return val
	}
	return ""
}

func AppendAttrs(attrs []Attr, list ...Attr) templ.Attributes {

	result := templ.Attributes{}
	finalAttrs := append(attrs, list...)
	for _, attr := range finalAttrs {
		if !attr.Hide {
			value := attr.Value
			if value == "" {
				//value = attr.Name
			}
			result[attr.Name] = value
		}
	}
	return result
}

func ifElse(val string, returnVal string) string {
	if len(val) > 0 {
		return val
	}
	return returnVal
}

func getAttribute(attrs []Attr, key string, defaultValue string) string {
	for _, attr := range attrs {
		if strings.EqualFold(attr.Name, key) {
			return ifElse(attr.Value, defaultValue)
		}
	}
	return defaultValue
}
