package render

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"github.com/valyala/fasthttp"
	"html/template"
)

// Engine is the generic interface for all responses.
type Engine interface {
	Render(*fasthttp.RequestCtx, interface{}) error
}

// Head defines the basic ContentType and Status fields.
type Head struct {
	ContentType string
	Status      int
}

// Data built-in renderer.
type Data struct {
	Head
}

// HTML built-in renderer.
type HTML struct {
	Head
	Name      string
	Templates *template.Template
}

// JSON built-in renderer.
type JSON struct {
	Head
	Indent        bool
	UnEscapeHTML  bool
	Prefix        []byte
	StreamingJSON bool
}

// JSONP built-in renderer.
type JSONP struct {
	Head
	Indent   bool
	Callback string
}

// Text built-in renderer.
type Text struct {
	Head
}

// XML built-in renderer.
type XML struct {
	Head
	Indent bool
	Prefix []byte
}

// Write outputs the header content.
func (h Head) Write(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Set(ContentType, h.ContentType)
	ctx.SetStatusCode(h.Status)
}

// Render a data response.
func (d Data) Render(ctx *fasthttp.RequestCtx, v interface{}) error {
	c := string(ctx.Request.Header.Peek(ContentType))
	w := ctx.Response.BodyWriter()
	if c != "" {
		d.Head.ContentType = c
	}

	d.Head.Write(ctx)
	w.Write(v.([]byte))
	return nil
}

// Render a HTML response.
func (h HTML) Render(ctx *fasthttp.RequestCtx, binding interface{}) error {
	// Retrieve a buffer from the pool to write to.
	out := bufPool.Get()
	err := h.Templates.ExecuteTemplate(out, h.Name, binding)
	if err != nil {
		return err
	}
	w := ctx.Response.BodyWriter()
	h.Head.Write(ctx)
	out.WriteTo(w)

	// Return the buffer to the pool.
	bufPool.Put(out)
	return nil
}

// Render a JSON response.
func (j JSON) Render(ctx *fasthttp.RequestCtx, v interface{}) error {
	if j.StreamingJSON {
		return j.renderStreamingJSON(ctx, v)
	}

	var result []byte
	var err error

	if j.Indent {
		result, err = json.MarshalIndent(v, "", "  ")
		result = append(result, '\n')
	} else {
		result, err = json.Marshal(v)
	}
	if err != nil {
		return err
	}

	// Unescape HTML if needed.
	if j.UnEscapeHTML {
		result = bytes.Replace(result, []byte("\\u003c"), []byte("<"), -1)
		result = bytes.Replace(result, []byte("\\u003e"), []byte(">"), -1)
		result = bytes.Replace(result, []byte("\\u0026"), []byte("&"), -1)
	}
	w := ctx.Response.BodyWriter()
	// JSON marshaled fine, write out the result.
	j.Head.Write(ctx)
	if len(j.Prefix) > 0 {
		w.Write(j.Prefix)
	}
	w.Write(result)
	return nil
}

func (j JSON) renderStreamingJSON(ctx *fasthttp.RequestCtx, v interface{}) error {
	j.Head.Write(ctx)
	w := ctx.Response.BodyWriter()
	if len(j.Prefix) > 0 {
		w.Write(j.Prefix)
	}

	return json.NewEncoder(w).Encode(v)
}

// Render a JSONP response.
func (j JSONP) Render(ctx *fasthttp.RequestCtx, v interface{}) error {
	var result []byte
	var err error

	if j.Indent {
		result, err = json.MarshalIndent(v, "", "  ")
	} else {
		result, err = json.Marshal(v)
	}
	if err != nil {
		return err
	}
	w := ctx.Response.BodyWriter()

	// JSON marshaled fine, write out the result.
	j.Head.Write(ctx)
	w.Write([]byte(j.Callback + "("))
	w.Write(result)
	w.Write([]byte(");"))

	// If indenting, append a new line.
	if j.Indent {
		w.Write([]byte("\n"))
	}
	return nil
}

// Render a text response.
func (t Text) Render(ctx *fasthttp.RequestCtx, v interface{}) error {
	c := string(ctx.Request.Header.Peek(ContentType))
	if c != "" {
		t.Head.ContentType = c
	}
	w := ctx.Response.BodyWriter()
	t.Head.Write(ctx)
	w.Write([]byte(v.(string)))
	return nil
}

// Render an XML response.
func (x XML) Render(ctx *fasthttp.RequestCtx, v interface{}) error {
	var result []byte
	var err error

	if x.Indent {
		result, err = xml.MarshalIndent(v, "", "  ")
		result = append(result, '\n')
	} else {
		result, err = xml.Marshal(v)
	}
	if err != nil {
		return err
	}

	// XML marshaled fine, write out the result.
	x.Head.Write(ctx)
	w := ctx.Response.BodyWriter()
	if len(x.Prefix) > 0 {
		w.Write(x.Prefix)
	}
	w.Write(result)
	return nil
}