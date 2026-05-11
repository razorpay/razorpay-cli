package output

import (
	"bytes"
	"io"
	"sort"
	"strings"
	"testing"
)

// sample is a representative Razorpay collection response: integer
// timestamps, string IDs, null fields, nested arrays. Each formatter test
// asserts behaviour against this same payload.
const sample = `{
  "entity": "collection",
  "count": 2,
  "items": [
    {"id": "cust_A", "created_at": 1778486691, "gstin": null, "amount": 12.5},
    {"id": "cust_B", "created_at": 1778486682, "gstin": null, "amount": 0}
  ]
}`

func render(t *testing.T, format, payload string) string {
	t.Helper()
	f, ok := Get(format)
	if !ok {
		t.Fatalf("format %q not registered", format)
	}
	var buf bytes.Buffer
	if err := f.Format(&buf, []byte(payload)); err != nil {
		t.Fatalf("format %q failed: %v", format, err)
	}
	return buf.String()
}

// ---- JSON ----

func TestJSONFormat_pretty_prints_with_two_space_indent(t *testing.T) {
	out := render(t, "json", sample)
	// Pretty-printed JSON has indented keys and the entity tag verbatim.
	if !strings.Contains(out, `"entity": "collection"`) {
		t.Fatalf("json output missing entity tag:\n%s", out)
	}
	if !strings.Contains(out, "\n  \"") {
		t.Fatalf("json output is not indented:\n%s", out)
	}
}

func TestJSONFormat_passthrough_on_invalid_json(t *testing.T) {
	// The JSON formatter is intentionally permissive — if the payload isn't
	// valid JSON it writes the bytes verbatim instead of erroring, so the
	// user still sees whatever the API returned.
	var buf bytes.Buffer
	f, _ := Get("json")
	if err := f.Format(&buf, []byte("not really json")); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if buf.String() != "not really json" {
		t.Fatalf("json passthrough failed: got %q", buf.String())
	}
}

// ---- YAML ----

func TestYAMLFormat_renders_known_keys(t *testing.T) {
	out := render(t, "yaml", sample)
	for _, want := range []string{"entity: collection", "count: 2", "id: cust_A", "id: cust_B"} {
		if !strings.Contains(out, want) {
			t.Fatalf("yaml output missing %q:\n%s", want, out)
		}
	}
}

func TestYAMLFormat_keeps_integers_as_integers(t *testing.T) {
	// Regression: json.Unmarshal widens every number to float64, which
	// yaml.v3 then renders as "1.778486691e+09". UseNumber + narrowNumbers
	// keeps the timestamp printable.
	out := render(t, "yaml", sample)
	if !strings.Contains(out, "created_at: 1778486691") {
		t.Fatalf("yaml widened integer timestamp:\n%s", out)
	}
	if strings.Contains(out, "e+09") || strings.Contains(out, "e+9") {
		t.Fatalf("yaml emitted scientific notation:\n%s", out)
	}
}

// ---- TOML ----

func TestTOMLFormat_renders_table_and_array_of_tables(t *testing.T) {
	out := render(t, "toml", sample)
	if !strings.Contains(out, "entity = 'collection'") {
		t.Fatalf("toml output missing entity row:\n%s", out)
	}
	if !strings.Contains(out, "[[items]]") {
		t.Fatalf("toml output missing items array:\n%s", out)
	}
	if !strings.Contains(out, "created_at = 1778486691") {
		t.Fatalf("toml widened integer timestamp:\n%s", out)
	}
}

func TestTOMLFormat_strips_null_values(t *testing.T) {
	// TOML has no concept of null; null fields are dropped rather than
	// triggering an encoder error.
	out := render(t, "toml", sample)
	if strings.Contains(out, "gstin") {
		t.Fatalf("toml output should have stripped null `gstin`:\n%s", out)
	}
}

func TestTOMLFormat_wraps_non_object_roots(t *testing.T) {
	// TOML's top level must be a table. Wrapping under "value" lets us
	// serialise bare arrays / scalars that some endpoints might return.
	out := render(t, "toml", `[1, 2, 3]`)
	if !strings.Contains(out, "value = [1, 2, 3]") {
		t.Fatalf("toml did not wrap bare array:\n%s", out)
	}
}

// ---- registry / dispatch ----

func TestPrint_unknown_format_falls_back_to_json_with_warning(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fprint(&stdout, &stderr, "xml", []byte(sample))

	if !strings.Contains(stderr.String(), `unknown output format "xml"`) {
		t.Fatalf("expected stderr warning, got: %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), `"entity": "collection"`) {
		t.Fatalf("expected JSON fallback on stdout, got: %q", stdout.String())
	}
}

func TestPrint_empty_format_uses_default(t *testing.T) {
	var stdout, stderr bytes.Buffer
	fprint(&stdout, &stderr, "", []byte(sample))

	if stderr.Len() != 0 {
		t.Fatalf("did not expect stderr output for empty format: %q", stderr.String())
	}
	if !strings.Contains(stdout.String(), `"entity": "collection"`) {
		t.Fatalf("expected default JSON output, got: %q", stdout.String())
	}
}

func TestRegister_is_extensible(t *testing.T) {
	const name = "test-upper"
	t.Cleanup(func() {
		// Best-effort cleanup so later tests aren't polluted; Names() picks
		// up whatever is registered at the time it's called.
		mu.Lock()
		delete(registry, name)
		mu.Unlock()
	})

	Register(name, FormatterFunc(func(w io.Writer, _ []byte) error {
		_, err := io.WriteString(w, "HELLO")
		return err
	}))

	got, ok := Get(strings.ToUpper(name)) // case-insensitive lookup
	if !ok {
		t.Fatalf("Get did not find newly registered formatter")
	}
	var buf bytes.Buffer
	if err := got.Format(&buf, []byte(sample)); err != nil {
		t.Fatalf("custom formatter errored: %v", err)
	}
	if buf.String() != "HELLO" {
		t.Fatalf("custom formatter output mismatch: got %q", buf.String())
	}
}

func TestNames_includes_builtins_sorted(t *testing.T) {
	got := Names()

	// Build a sorted copy of what we expect to be present. Other tests may
	// have registered extras; we only care that the built-ins are there
	// and that the slice is sorted.
	wantPresent := []string{"json", "toml", "yaml", "yml"}
	for _, n := range wantPresent {
		found := false
		for _, g := range got {
			if g == n {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Names() missing built-in %q; got %v", n, got)
		}
	}

	if !sort.StringsAreSorted(got) {
		t.Fatalf("Names() not sorted: %v", got)
	}
}
