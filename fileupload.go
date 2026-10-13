package headlessui

import (
	"context"
	_ "embed"

	"github.com/a-h/templ"
)

//go:embed alpine/dist/fileupload.min.js
var fileUploadJS string

var fileUploadScript = templ.NewOnceHandle(templ.WithComponent(templ.Raw("<script>" + fileUploadJS + "</script>")))

func FileUploadScript() templ.Component {
	return fileUploadScript.Once()
}

type FileUploadProps struct {
	Base
	Disabled bool
}

type FileUploadInputProps struct {
	Base
	Name     string
	Accept   string
	Form     string
	Multiple bool
	Disabled bool
	Required bool
	Invalid  bool
}

type FileUploadDropzoneProps struct {
	Base
}

type FileUploadListProps struct {
	Base
}

type FileUploadFileProps struct {
	Base
}

type FileUploadFileNameProps struct {
	Base
}

type fileUploadKeyType struct{}

var fileUploadKey = fileUploadKeyType{}

func fileUploadContext(ctx context.Context, p FileUploadProps) (context.Context, FileUploadProps) {
	if p.ID == "" {
		p.ID = nextID(ctx, "file-upload")
	}
	return context.WithValue(ctx, fileUploadKey, p.ID), p
}

func fileUploadAttrs(p FileUploadProps) templ.Attributes {
	own := make(templ.Attributes, 4)
	own["x-data"] = "huiFileUpload"
	own[":data-dragging"] = "dragging"
	if p.Disabled {
		own["data-disabled"] = true
	}
	return p.merge(own)
}

func fileUploadInputAttrs(ctx context.Context, p FileUploadInputProps) templ.Attributes {
	own := make(templ.Attributes, 14)
	own["type"] = "file"

	f, inField := controlFrom(ctx, FieldProps{
		Name:     p.Name,
		Form:     p.Form,
		Disabled: p.Disabled,
		Required: p.Required,
		Invalid:  p.Invalid,
	})
	f.Readonly = false
	if inField {
		p.ID = f.ID
		if by := f.describedBy(); by != "" {
			own["aria-describedby"] = by
		}
	}
	if f.Name != "" {
		own["name"] = f.Name
	}
	if f.Form != "" {
		own["form"] = f.Form
	}
	if p.Accept != "" {
		own["accept"] = p.Accept
	}
	if p.Multiple {
		own["multiple"] = true
	}
	if f.Disabled {
		own["disabled"] = true
	}
	if f.Required {
		own["required"] = true
	}
	if f.Invalid {
		own["aria-invalid"] = "true"
	}

	if id, ok := ctx.Value(fileUploadKey).(string); ok {
		own["data-hui-file"] = id
		return p.mergeEvents(f.state(own), Events{OnChange: "changed($event)"})
	}
	return p.merge(f.state(own))
}

func fileUploadDropzoneAttrs(ctx context.Context, p FileUploadDropzoneProps) templ.Attributes {
	if _, ok := ctx.Value(fileUploadKey).(string); !ok {
		return p.merge(nil)
	}
	own := make(templ.Attributes, 8)
	own[":data-dragging"] = "dragging"
	own["@dragenter.prevent"] = chain("enter($event)", attrString(p.Attrs["@dragenter.prevent"]))
	own["@dragover.prevent"] = chain("over($event)", attrString(p.Attrs["@dragover.prevent"]))
	own["@dragleave"] = chain("leave($event)", attrString(p.Attrs["@dragleave"]))
	own["@drop.prevent"] = chain("drop($event)", attrString(p.Attrs["@drop.prevent"]))
	return p.mergeEvents(own, Events{OnClick: "browse($event)"})
}

func fileUploadListAttrs(ctx context.Context, p FileUploadListProps) templ.Attributes {
	if _, ok := ctx.Value(fileUploadKey).(string); !ok {
		return p.merge(nil)
	}
	return p.merge(templ.Attributes{"x-show": "files.length", "style": "display: none"})
}

func fileUploadFileAttrs(ctx context.Context, p FileUploadFileProps) (templ.Attributes, bool) {
	_, ok := ctx.Value(fileUploadKey).(string)
	return p.merge(nil), ok
}

func fileUploadFileNameAttrs(ctx context.Context, p FileUploadFileNameProps) templ.Attributes {
	if _, ok := ctx.Value(fileUploadKey).(string); !ok {
		return p.merge(nil)
	}
	return p.merge(templ.Attributes{"x-text": "file"})
}
