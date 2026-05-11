// Package output renders the raw JSON returned by the Razorpay API into
// whichever presentation format the user has configured. The API layer is
// always JSON-on-the-wire; this package is the translation hop.
//
// To add a new format (e.g. XML, msgpack, CSV), implement Formatter and
// call Register("<name>", f) — typically from a package init(). The CLI's
// `configure` command will then accept the new name automatically.
package output

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"sync"

	toml "github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// DefaultFormat is the format used when none is configured.
const DefaultFormat = "json"

// Formatter renders a JSON byte buffer into a target output format.
// Implementations must be safe for concurrent use.
type Formatter interface {
	Format(w io.Writer, data []byte) error
}

// FormatterFunc adapts a function into a Formatter.
type FormatterFunc func(w io.Writer, data []byte) error

// Format implements Formatter.
func (f FormatterFunc) Format(w io.Writer, data []byte) error { return f(w, data) }

var (
	mu       sync.RWMutex
	registry = map[string]Formatter{}
)

// Register adds a Formatter under the given (case-insensitive) name.
// Re-registering an existing name overwrites it.
func Register(name string, f Formatter) {
	mu.Lock()
	defer mu.Unlock()
	registry[normalise(name)] = f
}

// Get returns the Formatter registered under name (case-insensitive).
func Get(name string) (Formatter, bool) {
	mu.RLock()
	defer mu.RUnlock()
	f, ok := registry[normalise(name)]
	return f, ok
}

// IsRegistered reports whether the given name has a Formatter.
func IsRegistered(name string) bool {
	_, ok := Get(name)
	return ok
}

// Names returns the registered format names, sorted, for use in help text
// and validation error messages.
func Names() []string {
	mu.RLock()
	defer mu.RUnlock()
	out := make([]string, 0, len(registry))
	for k := range registry {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Print writes data to stdout using the formatter named by `format`. If
// `format` is empty or unknown, it falls back to JSON and reports the
// problem to stderr.
func Print(format string, data []byte) {
	fprint(os.Stdout, os.Stderr, format, data)
}

// fprint is the testable seam for Print.
func fprint(stdout, stderr io.Writer, format string, data []byte) {
	if normalise(format) == "" {
		format = DefaultFormat
	}
	f, ok := Get(format)
	if !ok {
		fmt.Fprintf(stderr, "warning: unknown output format %q, falling back to %s\n", format, DefaultFormat)
		f, _ = Get(DefaultFormat)
	}
	if err := f.Format(stdout, data); err != nil {
		fmt.Fprintf(stderr, "warning: failed to render as %s (%v), falling back to %s\n", format, err, DefaultFormat)
		fallback, _ := Get(DefaultFormat)
		_ = fallback.Format(stdout, data)
	}
}

func normalise(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

func init() {
	Register("json", FormatterFunc(formatJSON))
	Register("yaml", FormatterFunc(formatYAML))
	Register("toml", FormatterFunc(formatTOML))
}

func formatJSON(w io.Writer, data []byte) error {
	var v any
	if err := json.Unmarshal(data, &v); err != nil {
		// Not JSON — preserve the payload verbatim rather than dropping it.
		_, werr := w.Write(data)
		return werr
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func formatYAML(w io.Writer, data []byte) error {
	v, err := decodeJSON(data)
	if err != nil {
		return err
	}
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(v); err != nil {
		_ = enc.Close()
		return err
	}
	return enc.Close()
}

func formatTOML(w io.Writer, data []byte) error {
	v, err := decodeJSON(data)
	if err != nil {
		return err
	}
	// TOML's top level must be a table. API list responses already are (they
	// come back as {entity, count, items}), but if a future endpoint returns
	// a bare array or scalar, wrap it so the encode succeeds.
	cleaned := stripNil(v)
	root, ok := cleaned.(map[string]any)
	if !ok {
		root = map[string]any{"value": cleaned}
	}
	return toml.NewEncoder(w).Encode(root)
}

// decodeJSON unmarshals via a json.Decoder with UseNumber, then narrows
// each json.Number back to int64 or float64. Without this, every JSON
// number becomes float64 and a Unix timestamp like 1778486691 renders as
// "1.778486691e+09" in YAML / TOML — fine arithmetically, wrong visually.
func decodeJSON(data []byte) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, fmt.Errorf("payload is not valid JSON: %w", err)
	}
	return narrowNumbers(v), nil
}

func narrowNumbers(v any) any {
	switch x := v.(type) {
	case map[string]any:
		for k, vv := range x {
			x[k] = narrowNumbers(vv)
		}
		return x
	case []any:
		for i, vv := range x {
			x[i] = narrowNumbers(vv)
		}
		return x
	case json.Number:
		if i, err := x.Int64(); err == nil {
			return i
		}
		if f, err := x.Float64(); err == nil {
			return f
		}
		return x.String()
	default:
		return v
	}
}

// stripNil recursively removes nil entries from maps and slices. TOML and
// some YAML consumers reject explicit nulls; this keeps the output clean
// without altering non-null fields.
func stripNil(v any) any {
	switch x := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(x))
		for k, vv := range x {
			sv := stripNil(vv)
			if sv == nil {
				continue
			}
			out[k] = sv
		}
		return out
	case []any:
		out := make([]any, 0, len(x))
		for _, vv := range x {
			sv := stripNil(vv)
			if sv == nil {
				continue
			}
			out = append(out, sv)
		}
		return out
	default:
		return v
	}
}
