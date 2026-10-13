package headlessui

import (
	"context"
	"regexp"
	"testing"

	"github.com/a-h/templ"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileUpload(t *testing.T) {
	t.Run("the input is a native file input without JavaScript", func(t *testing.T) {
		html := testRender(t, FileUploadInput(FileUploadInputProps{Name: "avatar", Accept: "image/*", Form: "profile"}))

		assert.Equal(t, `<input accept="image/*" form="profile" name="avatar" type="file">`, html)
	})

	t.Run("a dropzone fills the input", func(t *testing.T) {
		html := testRender(t, FileUpload(FileUploadProps{ID: "docs"}),
			FileUploadInput(FileUploadInputProps{Name: "docs", Multiple: true}),
			FileUploadDropzone(FileUploadDropzoneProps{}),
		)
		html = regexp.MustCompile(`(?s)<script>.*?</script>`).ReplaceAllString(html, "")

		assert.Contains(t, html, `<div :data-dragging="dragging" id="docs" x-data="huiFileUpload">`)
		assert.Contains(t, html, `@change="changed($event)" data-hui-file="docs" multiple name="docs" type="file"`)
		assert.Contains(t, html, `@drop.prevent="drop($event)"`)
		assert.Contains(t, html, `@click="browse($event)"`)
	})

	t.Run("Field is the single source of truth", func(t *testing.T) {
		html := testRender(t, Field(FieldProps{Name: "cv", Required: true, Error: "Upload your CV"}),
			Label(LabelProps{}),
			FileUploadInput(FileUploadInputProps{Name: "own"}),
		)

		assert.Contains(t, html, `id="cv" name="cv" required type="file"`)
		assert.Contains(t, html, `aria-describedby="cv-error"`)
		assert.NotContains(t, html, "own")
	})

	t.Run("a dropzone outside FileUpload is a plain element", func(t *testing.T) {
		html := testRender(t, FileUploadDropzone(FileUploadDropzoneProps{Class: "border-dashed"}))

		assert.Equal(t, `<div class="border-dashed"></div>`, html)
	})

	t.Run("escape hatches survive the internal wiring", func(t *testing.T) {
		ctx, _ := fileUploadContext(context.Background(), FileUploadProps{Base: Base{ID: "docs"}})

		input := fileUploadInputAttrs(ctx, FileUploadInputProps{Base: Base{
			Class:  "sr-only",
			Attrs:  templ.Attributes{"hx-post": "/upload", "hx-encoding": "multipart/form-data", "hx-trigger": "change"},
			Events: Events{OnChange: "track()"},
		}})
		zone := fileUploadDropzoneAttrs(ctx, FileUploadDropzoneProps{Base: Base{Attrs: templ.Attributes{"@drop.prevent": "log()"}}})

		assert.Equal(t, "multipart/form-data", input["hx-encoding"])
		assert.Equal(t, "changed($event); track()", input["@change"])
		assert.Equal(t, "drop($event); log()", zone["@drop.prevent"])
	})

	t.Run("the shipped script is valid and self registering", func(t *testing.T) {
		require.NotEmpty(t, fileUploadJS)
		assert.Contains(t, fileUploadJS, "window.Alpine.plugin")
		assert.NotContains(t, fileUploadJS, "export")
	})

	t.Run("lists the chosen file names", func(t *testing.T) {
		ctx, _ := fileUploadContext(context.Background(), FileUploadProps{Base: Base{ID: "docs"}})

		html := testRenderIn(t, ctx, FileUploadList(FileUploadListProps{Class: "files"}), testWithChildren(FileUploadFile(FileUploadFileProps{Class: "file"}), FileUploadFileName(FileUploadFileNameProps{})))

		assert.Contains(t, html, `<ul class="files" style="display: none" x-show="files.length">`)
		assert.Contains(t, html, `<template x-for="(file, index) in files" :key="index"><li class="file"><span x-text="file"></span></li></template>`)
	})

	t.Run("the list parts render plainly outside a file upload", func(t *testing.T) {
		name := gofakeit.Word()

		html := testRender(t, FileUploadFile(FileUploadFileProps{}), FileUploadFileName(FileUploadFileNameProps{}))
		plain := testRender(t, FileUploadFileName(FileUploadFileNameProps{}), templ.Raw(name))

		assert.Equal(t, `<li><span></span></li>`, html)
		assert.Equal(t, `<span>`+name+`</span>`, plain)
	})
}
