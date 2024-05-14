# Patchgen

generate kubernetes resource PATCH

## Install

```shell
go install github.com/lyp256/patchgen # go version 1.16+

```

## Usage

```shell
patchgen [raw file] [update file]
  -t, --type string   merge|strategic (default "merge")
 ```

### example

```shell
> patchgen .\raw.json .\new.json
{"spec":{"template":{"spec":{"volumes":[{"hostPath":{"path":"/proc","type":""},"name":"proc"},{"hostPath":{"path":"/sys","type":""},"name":"sys"},{"hostPath
":{"path":"/","type":""},"name":"root"},{"configMap":{"name":"node-export-web-config"},"name":"web-config"}]}}}}
```
