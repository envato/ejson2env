package ejson2env

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestReadSecrets(t *testing.T) {
	var err error

	_, err = ReadAndExtractEnv("testdata/bad.ejson", "./key", TestKeyValue)
	if nil == err {
		t.Fatal("failed to fail when loading a broken ejson file")
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("error should be \"no such file or directory\": %s", err)
	}
}

func TestReadAndExportEnv(t *testing.T) {
	outputBuffer := new(bytes.Buffer)
	output = outputBuffer

	// ensure that output returns to os.Stdout
	defer func() {
		output = os.Stdout
	}()

	tests := []struct {
		name           string
		exportFunc     ExportFunction
		ejsonFile      string
		expectedOutput string
	}{
		{
			name:           "ExportEnv",
			exportFunc:     ExportEnv,
			ejsonFile:      "testdata/test-expected-usage.ejson",
			expectedOutput: "export test_key='test value'\n",
		},
		{
			name:           "ExportQuiet",
			exportFunc:     ExportQuiet,
			ejsonFile:      "testdata/test-expected-usage.ejson",
			expectedOutput: "test_key='test value'\n",
		},
		{
			name:           "ExportEnvTrimUnderscore",
			exportFunc:     TrimLeadingUnderscoreExportWrapper(ExportEnv),
			ejsonFile:      "testdata/test-leading-underscore-env-key.ejson",
			expectedOutput: "export test_key='test value'\n",
		},
		{
			name:           "ExportEnvNoTrimUnderscore",
			exportFunc:     ExportEnv,
			ejsonFile:      "testdata/test-leading-underscore-env-key.ejson",
			expectedOutput: "export _test_key='test value'\n",
		},
	}

	for _, test := range tests {
		err := ReadAndExportEnv(test.ejsonFile, "./key", TestKeyValue, test.exportFunc)
		if nil != err {
			t.Errorf("testing %s failed: %s", test.name, err)
			continue
		}

		actualOutput := outputBuffer.String()

		if test.expectedOutput != actualOutput {
			t.Error(formatInvalid(actualOutput, test.expectedOutput))
		}
		outputBuffer.Reset()
	}
}

func TestReadAndExportEnvWithBadEjson(t *testing.T) {
	var err error

	outputBuffer := new(bytes.Buffer)
	output = outputBuffer

	// ensure that output returns to os.Stdout
	defer func() {
		output = os.Stdout
	}()

	err = ReadAndExportEnv("bad.ejson", "./key", TestKeyValue, ExportEnv)
	if nil == err {
		t.Fatal("failed to fail when loading a broken ejson file")
	}
	if !strings.Contains(err.Error(), "no such file or directory") {
		t.Errorf("error should be \"no such file or directory\": %s", err)
	}
}
