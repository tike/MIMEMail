package MIMEMail

import "testing"

type testConf struct {
	*Config
	sender   sender
	receiver receiver
}

type sender struct {
	address string
	key     string
	pass    string
}

type receiver struct {
	address string
	key     string
}

var testConfig *testConf

func getTestConfig(t *testing.T) *testConf {
	if testConfig != nil {
		return testConfig
	}

	t.Skip("\n\033[33;4mYou need to pass a test config by creating \033[31;4maccount_test.go!\033[0m\n\n" +

		"see: \033[32;4maccount_sample_test.go\033[0m for details.\n" +
		"You can copy it to account_test.go (.gitignore'd) and fill in the necessary details")
	return nil
}
