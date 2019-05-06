# Kubesecrets

Kubesecrets will read an input Helm values YAML file, traverse its values, transform them, and write the transformed output to a different file.

Currently the only supported transformations are encoding and decoding leaf / scalar values (`bool`, `int`, `float`, `string`) to and from base64. Non-string types are first stringified before they are encoded. This allows you to manage your secrets in plain text and easily encode them as Kubernetes expects.

The output is written to standard output. If multiple inputs are specified, they are separated by `---` on its own line in the ouput. If no input files are specified, input is read from standard in.

## Installing Kubesecrets

```
$ cd cmd
$ go install kubesecrets.go
```

## Using Kubesecrets

```
$ kubesecrets enc -verbose secrets.yaml
enc: secrets.yaml
```

Refer to the [`example`](https://github.com/bww/kubesecrets/tree/master/example) directory for examples of [input](https://github.com/bww/kubesecrets/blob/master/example/secrets.yaml) and [encoded](https://github.com/bww/kubesecrets/blob/master/example/secrets.enc.yaml) files.
