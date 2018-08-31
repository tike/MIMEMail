package templated

import "testing"

func TestRender(t *testing.T) {
	cnf := Config{"example", "en_US"}
	ctx := map[string]interface{}{
		"Name":    "Test Customer",
		"Company": "ACME Coorp",
	}

	subj, body, err := Render(&cnf, "foo", ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", string(subj))
	t.Logf("%s", string(body))
}
