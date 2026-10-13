package headlessui

import (
	"context"
	"strconv"
	"strings"

	"github.com/a-h/templ"
)

type checkableKeyType struct{}

var checkableKey = checkableKeyType{}

type checkableState struct {
	checked       bool
	indeterminate bool
}

func checkableFromCtx(ctx context.Context) (checkableState, bool) {
	s, ok := ctx.Value(checkableKey).(checkableState)
	return s, ok
}

func checkableContext(ctx context.Context, s checkableState) context.Context {
	return context.WithValue(ctx, checkableKey, s)
}

type checkable struct {
	attrs    templ.Attributes
	input    templ.Attributes
	required templ.Attributes
	state    checkableState
}

func checkableValue(value string) string {
	if value == "" {
		return "on"
	}
	return value
}

func ariaChecked(checked, indeterminate bool) string {
	if indeterminate {
		return "mixed"
	}
	return strconv.FormatBool(checked)
}

func checkableData(component string, f FieldProps, indeterminate bool) string {
	config := make([]string, 0, 6)
	config = append(config, "checked: "+strconv.FormatBool(f.checked))
	if indeterminate {
		config = append(config, "indeterminate: true")
	}
	if f.Name != "" {
		config = append(config, "name: "+jsString(f.Name), "value: "+jsString(checkableValue(f.Value)))
	}
	if f.Form != "" {
		config = append(config, "form: "+jsString(f.Form))
	}
	if f.Disabled {
		config = append(config, "disabled: true")
	}
	return component + "({ " + strings.Join(config, ", ") + " })"
}

func checkableButton(b Base, role, component string, f FieldProps, inField, indeterminate bool, internal Events) checkable {
	own := make(templ.Attributes, 20)
	own["type"] = "button"
	own["role"] = role
	own["x-data"] = checkableData(component, f, indeterminate)
	own["x-modelable"] = "checked"
	own["aria-checked"] = ariaChecked(f.checked, indeterminate)
	own[":aria-checked"] = "ariaChecked()"
	own[":data-checked"] = "checked"
	if f.checked {
		own["data-checked"] = true
	}
	if role == "checkbox" {
		own[":data-indeterminate"] = "indeterminate"
		if indeterminate {
			own["data-indeterminate"] = true
		}
	}
	if inField {
		b.ID = f.ID
		if by := f.describedBy(); by != "" {
			own["aria-describedby"] = by
		}
	}
	if f.Form != "" {
		own["form"] = f.Form
	}
	if f.Name != "" && f.checked {
		own["name"] = f.Name
		own["value"] = checkableValue(f.Value)
	}
	if f.Disabled {
		own["disabled"] = true
	}
	if f.Invalid {
		own["aria-invalid"] = "true"
	}
	required := f.Required && (role == "radio" || !f.option)
	if required && role != "radio" {
		own["aria-required"] = "true"
	}

	c := checkable{
		attrs: b.mergeEvents(f.state(own), internal),
		state: checkableState{checked: f.checked, indeterminate: indeterminate},
	}
	if f.Name != "" && f.checked {
		c.input = templ.Attributes{"type": "hidden", "name": f.Name, "value": checkableValue(f.Value), "data-hui-checked": true}
		if f.Form != "" {
			c.input["form"] = f.Form
		}
		if f.Disabled {
			c.input["disabled"] = true
		}
	}
	if required {
		c.required = validationAttrs(b.ID, f.checked, f.Form, f.Disabled)
	}
	return c
}

func checkableIndicatorAttrs(ctx context.Context, b Base, show string, visible func(checkableState) bool) templ.Attributes {
	own := templ.Attributes{"aria-hidden": "true", "x-show": show}
	transition(own, b.Attrs, "x-transition")
	if s, ok := checkableFromCtx(ctx); ok && !visible(s) {
		own["style"] = "display: none"
	}
	return b.merge(own)
}
