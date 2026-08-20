package forms

import (
	"fmt"
	"strings"

	"github.com/a-h/templ"
	"github.com/ashish9868/rapidbackend/utils"
)

const (
	DATE_FORMAT             = "Y-m-d"
	TIME_FORMAT             = "H:i:S"
	DATETIME_FORMAT         = "Y-m-d H:i:S"
	DISPLAY_DATE_FORMAT     = "m/d/Y"
	DISPLAY_TIME_FORMAT     = "H:i"
	DISPLAY_DATETIME_FORMAT = "m/d/Y H:i"
)

type FormField struct {
	ID          string
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

type InputBaseProps struct {
	ID              string
	Varient         string
	Loading         bool
	Label           string
	Type            string
	Name            string
	Placeholder     string
	Value           string
	Required        bool
	LabelIcon       string
	IconBefore      string
	IconAfter       string
	Error           string
	HelperText      string
	ReadOnly        bool
	Disabled        bool
	Attrs           []Attr
	Selected        []string
	Options         []SelectOption
	Multiple        bool
	CustomComponent templ.Component
	ItemSize        int
	Signals         templ.Attributes
}

func (f *InputBaseProps) Render(formId string) templ.Component {
	switch f.Type {
	case "select":
		return Select(formId, *f)
	case "file":
		return FileInput(formId, *f)
	case "switch":
		return Switch(formId, *f)
	case "button", "submit", "reset":
		return Button(formId, *f)
	case "custom":
		return f.CustomComponent
	}
	return Input(formId, *f)
}

func (f *InputBaseProps) BuildID(formId string) string {
	return fmt.Sprintf("%s_%s", formId, utils.Coalesce(f.ID, f.Name))
}

func (f *InputBaseProps) BuildSignalAccessVar(formId string, key string) string {
	return fmt.Sprintf("$%s_%s", f.BuildID(formId), key)
}

func (f *InputBaseProps) BuildSignalVar(formId string, key string) string {
	return fmt.Sprintf("data-signals:%s_%s", f.BuildID(formId), key)
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
	if b.Type == "button" {
		return AppendAttrs(attrs, []Attr{
			{Name: "disabled", Hide: !b.Disabled},
			{Name: "type", Value: utils.Coalesce(b.Type, "button")},
		}...)
	}
	isDateOrTimeField := strings.Contains(b.Type, "date") || strings.Contains(b.Type, "time")
	attrs = append(attrs, []Attr{
		{Name: "disabled", Hide: !b.Disabled},
		{Name: "readonly", Hide: !b.ReadOnly},
		{Name: "required", Hide: !b.Required},
		{Name: "autocomplete", Value: "off"},
		{Name: "name", Value: fmt.Sprintf(`%s`, b.Name)},
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
		attrs = append(attrs, Attr{Name: "type", Value: utils.Coalesce(b.Type, "text")})
		attrs = append(attrs, Attr{Name: "value", Value: utils.Coalesce(b.Value, "")})
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

func getAttribute(attrs []Attr, key string, defaultValue string) string {
	for _, attr := range attrs {
		if strings.EqualFold(attr.Name, key) {
			return utils.Coalesce(attr.Value, defaultValue)
		}
	}
	return defaultValue
}
