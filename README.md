# Kubesecrets

Kubesecrets will read an input Helm values YAML file, traverse its values, transform them, and write the transformed output to a different file.

Currently the only supported transformation is encoding scalar values (`bool`, `int`, `float`, `string`) in base64. This allows you to manage your secrets in plain text and easily encode them as Kubernetes expects.

The output file is named the same as the input file with the extension `.enc.yaml` instead of `.yaml`.

## Installing Kubesecrets

```
$ cd cmd
$ go install kubesecrets.go
```

## Using Kubesecrets

```
$ kubesecrets enc -verbose secrets.yaml
enc: secrets.yaml -> secrets.enc.yaml
```

Refer to the [`example`](https://github.com/bww/kubesecrets/tree/master/example) directory for examples of [input](https://github.com/bww/kubesecrets/blob/master/example/secrets.yaml) and [encoded](https://github.com/bww/kubesecrets/blob/master/example/secrets.enc.yaml) files.
