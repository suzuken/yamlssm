package yamlssm

import (
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	yaml "gopkg.in/yaml.v2"
)

// toPtr retrun string pointer
func toPtr(s string) *string {
	return &s
}

// unmarshal is Test function for func Unmarshal()
// because use SSMmock.
func unmarshal(in []byte, out interface{}) error {
	if err := yaml.Unmarshal(in, out); err != nil {
		return err
	}

	mock := &mockSSMClient{}
	d := newTestssmDecrypter(mock)
	return d.override(out)
}

// mockSSMClient stores SSM interface for mock
type mockSSMClient struct {
	ssmiface.SSMAPI
}

// newTestssmDecrypter returns a new ssmDecrypter for mock.
func newTestssmDecrypter(mock ssmiface.SSMAPI) *ssmDecrypter {
	return &ssmDecrypter{
		svc: mock,
	}
}

// GetParameter returns "decrypted" that is Decrypted SSM parameter.
func (m *mockSSMClient) GetParameter(i *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	parameter := &ssm.Parameter{
		Value: aws.String("decrypted"),
	}

	return &ssm.GetParameterOutput{
		Parameter: parameter,
	}, nil
}

func TestSSMUnmarshal(t *testing.T) {
	cases := []struct {
		value    string
		expected interface{}
	}{
		// struct
		{
			"a: a\nb: ssm://encrypt_parameter",
			&struct {
				A, B string
			}{A: "a", B: "decrypted"},
		},
	}

	for _, c := range cases {
		v := reflect.ValueOf(c.expected).Type()

		var out interface{}
		switch v.Kind() {
		case reflect.Ptr:
			out = reflect.New(v.Elem()).Interface()
		default:
			t.Fatalf("missing case for %s", v)
		}

		if err := unmarshal([]byte(c.value), out); err != nil {
			t.Fatalf("failed unmarshal: %s", err)
		}

		switch c.expected.(type) {
		case string:
			c.expected = toPtr(c.expected.(string))
		}

		if !reflect.DeepEqual(c.expected, out) {
			t.Errorf("want %s got %s", c.expected, out)
		}
	}
}
