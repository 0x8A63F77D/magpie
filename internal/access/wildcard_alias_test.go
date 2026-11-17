package access

import "testing"

func TestWildcardAlias(t *testing.T) {
	as, err := CleanAliases(map[string]string{"*": "a/*", "fast": "b/gpt-5"})
	if err != nil {
		t.Fatal(err)
	}
	who := Identity{Aliases: as}
	for _, c := range []struct {
		in, want string
		ok       bool
	}{
		{"claude-opus-5-5", "a/claude-opus-5-5", true},
		{"FAST", "b/gpt-5", true},
		{"b/claude-opus-5-5", "", false},
		{"", "", false},
	} {
		got, ok := who.AliasFor(c.in)
		if got != c.want || ok != c.ok {
			t.Errorf("AliasFor(%q) = %q, %v; want %q, %v", c.in, got, ok, c.want, c.ok)
		}
	}
	if name, ok := who.AliasOf("a/*"); ok {
		t.Errorf("AliasOf(a/*) = %q; the * name lists nothing", name)
	}
	for _, bad := range []string{"a/x", "*", "a/b/*", "/*", "re:x/*"} {
		if _, err := CleanAliases(map[string]string{"*": bad}); err == nil {
			t.Errorf("CleanAliases(* = %q) took it", bad)
		}
	}
	if _, err := CleanAliases(map[string]string{"gpt*": "a/x"}); err == nil {
		t.Error("a name with * other than * itself was taken")
	}
}
