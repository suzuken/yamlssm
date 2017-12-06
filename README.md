# yamlssm

yamlssm is extended YAML format with Amazon Simple System Manager.

## How to use

yamlssm has expand ssm value as macro prefixed by `ssm://`. See example below:

```yaml
foo: ssm://prod.database.name
bar: test
```

`ssm://prod.database.name` should be set on you (System Manager Parameter Store)[http://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html]. If you set, yamlssm set a value stored in parameter store.

```go
type T struct {
	Foo string
	Bar string
}
var t T

if err := yamlssm.Unmarshal([]byte[`
foo: ssm://prod.database.name
bar: test
`], &t); err != nil {
    // ...
}

fmt.Print(t.Foo) // -> value of prod.database.name on your ssm
```

## References

http://docs.aws.amazon.com/systems-manager/latest/userguide/systems-manager-paramstore.html

## LICENSE

MIT
