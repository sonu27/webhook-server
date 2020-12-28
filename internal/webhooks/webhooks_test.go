package webhooks

import (
	"reflect"
	"testing"
)

func TestWebhooks_Add(t *testing.T) {
	w := new(Webhooks)
	w.Add("http://example.com", "test")

	var tests = []struct {
		name      string
		want, got interface{}
	}{
		{"size", 1, len(w.data)},
		{"url", "http://example.com", w.data[0].Url},
		{"token", "test", w.data[0].Token},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.got != test.want {
				t.Errorf("got %d, want %d", test.got, test.want)
			}
		})
	}
}

func TestWebhooks_Get(t *testing.T) {
	w := new(Webhooks)
	w.Add("http://example.com", "test")
	w.Add("http://example2.com", "test2")

	want := []webhook{{
		Url:   "http://example.com",
		Token: "test",
	}, {
		Url:   "http://example2.com",
		Token: "test2",
	}}

	if got := w.Get(); !reflect.DeepEqual(got, want) {
		t.Errorf("Get() = %v, want %v", got, want)
	}
}

func TestWebhooks_GetReturnsClone(t *testing.T) {
	w := new(Webhooks)
	w.Add("http://example.com", "test")

	want := []webhook{{
		Url:   "http://example.com",
		Token: "test",
	}}

	diff := w.Get()
	diff[0].Token = "diff"

	if got := w.Get(); !reflect.DeepEqual(got, want) {
		t.Errorf("Get() = %v, want %v", got, want)
	}
}
