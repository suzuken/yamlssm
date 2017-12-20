# yamlssm

yamlssm is extended YAML format with Amazon Simple System Manager.

## How to use

yamlssm has expand ssm value as macro prefixed by `ssm://`. See example below:

```yaml
foo: ssm://prod.database.name
bar: test
```

`ssm://prod.database.name` should be set on you [System Manager Parameter Store](http://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html). If you set, yamlssm set a value stored in parameter store.

```go
package main

import (
  "fmt"

  "github.com/suzuken/yamlssm"
)

type T struct {
  Foo string `yaml:"foo"`
  Bar string `yaml:"bar"`
}

func main() {
  t := T{}
  if err := yamlssm.Unmarshal([]byte(`
foo: ssm://prod.database.name
bar: test
`), &t)
  err != nil {
    fmt.Println(err)
  }

  fmt.Println(t.Foo) // -> value of prod.database.name on your ssm
}
```

## Notice

* Only supports default encryption key of your account.
* To set aws region, please use `AWS_REGION` environment variables. (ex. `export AWS_REGION='ap-northeast-1'`)
* This behavior based on aws-sdk-go. If you run yamlssm on your EC2 or some kind of instance on AWS, it's use that environment on default.

## References

http://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html

## LICENSE

MIT
