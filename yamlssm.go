package yamlssm

import (
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

type decoder struct {
	sess *session.Session
	svc  *ssm.SSM
}

func (d *decoder) expand(encrypted string) (string, error) {
	trimed := strings.TrimPrefix(encrypted, "ssm://")

	params := &ssm.GetParameterInput{
		Name:           aws.String(trimed),
		WithDecryption: aws.Bool(true),
	}
	resp, err := d.svc.GetParameter(params)
	if err != nil {
		return "", err
	}
	return *resp.Parameter.Value, nil
}

func (d *decoder) override(out interface{}) {
	val := reflect.ValueOf(out).Elem()
	if !val.IsValid() {
		return
	}
	if val.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		if !f.IsValid() {
			continue
		}
		if f.Kind() != reflect.String {
			continue
		}
		v := f.Interface().(string)
		if strings.HasPrefix(v, "ssm://") {
			actual, err := d.expand(v)
			if err != nil {
				// key hasn't found in ssm. path through.
				continue
			}
			if f.CanSet() {
				// override
				f.SetString(actual)
			}
		}
	}
}

func newDecoder() *decoder {
	sess := session.New()
	svc := ssm.New(sess)
	return &decoder{sess, svc}
}

// Unmarshal works same as gopkg.in/yaml.v2.
//
// After unmarshal YAML files, yamlssm replace value prefixed `ssm://`
// by encrypted value which stored in your System Manager Parameter Store.
func Unmarshal(in []byte, out interface{}) error {
	if err := yaml.Unmarshal(in, out); err != nil {
		return err
	}
	d := newDecoder()
	d.override(out)
	return nil
}
