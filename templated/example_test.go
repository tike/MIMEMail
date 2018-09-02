package templated

import (
	"fmt"
)

func Example() {
	cnf := Config{"example", "en_US"}
	ctx := map[string]interface{}{
		"Name":    "Test Customer",
		"Company": "ACME Coorp",
	}

	subj, body, err := Render(&cnf, "foo", ctx)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Subject: %s\n", string(subj))
	fmt.Printf("Body:\n%s\n", string(body))
}
