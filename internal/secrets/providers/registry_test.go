package providers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/user/envchain/internal/secrets/providers"
)

type fakeProvider struct {
	name   string
	values map[string]string
}

func (f *fakeProvider) Name() string { return f.name }
func (f *fakeProvider) Resolve(_ context.Context, ref string) (string, error) {
	v, ok := f.values[ref]
	if !ok {
		return "", errors.New("not found: " + ref)
	}
	return v, nil
}

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := providers.NewRegistry()
	r.Register(&fakeProvider{name: "test", values: map[string]string{}})

	p, err := r.Get("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Name() != "test" {
		t.Errorf("expected 'test', got %q", p.Name())
	}
}

func TestRegistry_Get_UnknownProvider(t *testing.T) {
	r := providers.NewRegistry()
	_, err := r.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown provider")
	}
}

func TestRegistry_Resolve_Dispatches(t *testing.T) {
	r := providers.NewRegistry()
	r.Register(&fakeProvider{
		name:   "fake",
		values: map[string]string{"mykey": "myvalue"},
	})

	val, err := r.Resolve(context.Background(), "fake", "mykey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "myvalue" {
		t.Errorf("expected 'myvalue', got %q", val)
	}
}

func TestRegistry_Names(t *testing.T) {
	r := providers.NewRegistry()
	r.Register(&fakeProvider{name: "a", values: nil})
	r.Register(&fakeProvider{name: "b", values: nil})

	names := r.Names()
	if len(names) != 2 {
		t.Errorf("expected 2 names, got %d", len(names))
	}
}

func TestRegistry_DuplicateRegister_Panics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on duplicate registration")
		}
	}()
	reg := providers.NewRegistry()
	reg.Register(&fakeProvider{name: "dup", values: nil})
	reg.Register(&fakeProvider{name: "dup", values: nil})
}
