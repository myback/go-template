package tmpl_functions

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestToYaml_UnsupportedNestedMapKey(t *testing.T) {
	expected := ``
	vals := map[string]interface{}{
		"foo": map[interface{}]interface{}{
			"bar": "BAR",
		},
	}
	actual, err := toYaml(vals)
	fmt.Println(actual)
	if err == nil {
		t.Fatalf("expected error but got none")
	} else if err.Error() != "error marshaling into JSON: json: unsupported type: map[interface {}]interface {}" {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected result: expected=%v, actual=%v", expected, actual)
	}
}

func TestToYaml(t *testing.T) {
	expected := `foo:
  bar: BAR
`
	vals := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "BAR",
		},
	}
	actual, err := toYaml(vals)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected result: expected=%v, actual=%v", expected, actual)
	}
}

func TestFromYaml(t *testing.T) {
	raw := `foo:
  bar: BAR
`
	expected := map[string]interface{}{
		"foo": map[string]interface{}{
			"bar": "BAR",
		},
	}
	actual, err := fromYaml(raw)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected result: expected=%v, actual=%v", expected, actual)
	}
}

func TestFromYamlToJson(t *testing.T) {
	input := `foo:
  bar: BAR
`
	want := `{"foo":{"bar":"BAR"}}`

	m, err := fromYaml(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := json.Marshal(m)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if d := cmp.Diff(want, string(got)); d != "" {
		t.Errorf("unexpected result: want (-), got (+):\n%s", d)
	}
}

func TestSetValueAtPath_OneComponent(t *testing.T) {
	input := map[string]interface{}{
		"foo": "",
	}
	expected := map[string]interface{}{
		"foo": "FOO",
	}
	actual, err := setValueAtPath("foo", "FOO", input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected result: expected=%v, actual=%v", expected, actual)
	}
}

func TestSetValueAtPath_TwoComponents(t *testing.T) {
	input := map[string]interface{}{
		"foo": map[interface{}]interface{}{
			"bar": "",
		},
	}
	expected := map[string]interface{}{
		"foo": map[interface{}]interface{}{
			"bar": "FOO_BAR",
		},
	}
	actual, err := setValueAtPath("foo.bar", "FOO_BAR", input)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected result: expected=%v, actual=%v", expected, actual)
	}
}

func TestTpl(t *testing.T) {
	text := `foo: {{ .foo }}
`
	expected := `foo: FOO
`
	actual, err := tpl(text, map[string]interface{}{"foo": "FOO"})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("unexpected result: expected=%v, actual=%v", expected, actual)
	}
}

func TestRequired(t *testing.T) {
	type args struct {
		warn string
		val  interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr bool
	}{
		{
			name:    "required val is nil",
			args:    args{warn: "This value is required", val: nil},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "required val is empty string",
			args:    args{warn: "This value is required", val: ""},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "required val is existed",
			args:    args{warn: "This value is required", val: "foo"},
			want:    "foo",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		testCase := tt
		t.Run(testCase.name, func(t *testing.T) {
			got, err := required(testCase.args.warn, testCase.args.val)
			if (err != nil) != testCase.wantErr {
				t.Errorf("Required() error = %v, wantErr %v", err, testCase.wantErr)
				return
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("Required() got = %v, want %v", got, testCase.want)
			}
		})
	}
}

type EmptyStruct struct {
}

func TestGetStruct(t *testing.T) {
	type Foo struct{ Bar string }

	obj := struct{ Foo }{Foo{Bar: "Bar"}}

	v1, err := get("Foo.Bar", obj)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v1 != "Bar" {
		t.Errorf("unexpected value for path Foo.Bar in %v: expected=Bar, actual=%v", obj, v1)
	}

	_, err = get("Foo.baz", obj)

	if err == nil {
		t.Errorf("expected error but was not occurred")
	}

	_, err = get("foo", EmptyStruct{})

	if err == nil {
		t.Errorf("expected error but was not occurred")
	}
}

func TestGetMap(t *testing.T) {
	obj := map[string]interface{}{"Foo": map[string]interface{}{"Bar": "Bar"}}

	v1, err := get("Foo.Bar", obj)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v1 != "Bar" {
		t.Errorf("unexpected value for path Foo.Bar in %v: expected=Bar, actual=%v", obj, v1)
	}

	_, err = get("Foo.baz", obj)

	if err == nil {
		t.Errorf("expected error but was not occurred")
	}
}

func TestGetMapPtr(t *testing.T) {
	obj := map[string]interface{}{"Foo": map[string]interface{}{"Bar": "Bar"}}
	objPrt := &obj

	v1, err := get("Foo.Bar", objPrt)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v1 != "Bar" {
		t.Errorf("unexpected value for path Foo.Bar in %v: expected=Bar, actual=%v", objPrt, v1)
	}

	_, err = get("Foo.baz", objPrt)

	if err == nil {
		t.Errorf("expected error but was not occurred")
	}
}

func TestGet_Default(t *testing.T) {
	obj := map[string]interface{}{"Foo": map[string]interface{}{}, "foo": 1}

	v1, err := get("Foo.Bar", "Bar", obj)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v1 != "Bar" {
		t.Errorf("unexpected value for path Foo.Bar in %v: expected=Bar, actual=%v", obj, v1)
	}

	v2, err := get("Baz", "Baz", obj)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v2 != "Baz" {
		t.Errorf("unexpected value for path Baz in %v: expected=Baz, actual=%v", obj, v2)
	}

	_, err = get("foo.Bar", "fooBar", obj)

	if err == nil {
		t.Errorf("expected error but was not occurred")
	}
}

func TestGetOrNilStruct(t *testing.T) {
	type Foo struct{ Bar string }

	obj := struct{ Foo }{Foo{Bar: "Bar"}}

	v1, err := getOrNil("Foo.Bar", obj)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v1 != "Bar" {
		t.Errorf("unexpected value for path Foo.Bar in %v: expected=Bar, actual=%v", obj, v1)
	}

	v2, err := getOrNil("Foo.baz", obj)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v2 != nil {
		t.Errorf("unexpected value for path Foo.baz in %v: expected=nil, actual=%v", obj, v2)
	}
}

func TestGetOrNilMap(t *testing.T) {
	obj := map[string]interface{}{"Foo": map[string]interface{}{"Bar": "Bar"}}

	v1, err := getOrNil("Foo.Bar", obj)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v1 != "Bar" {
		t.Errorf("unexpected value for path Foo.Bar in %v: expected=Bar, actual=%v", obj, v1)
	}

	v2, err := getOrNil("Foo.baz", obj)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if v2 != nil {
		t.Errorf("unexpected value for path Foo.baz in %v: expected=nil, actual=%v", obj, v2)
	}
}
