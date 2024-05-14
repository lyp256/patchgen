package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	serializeryaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var decUnstructured = serializeryaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)

func usage() {
	_, exe := filepath.Split(os.Args[0])
	fmt.Fprintf(os.Stderr, "usage: %s [raw file] [update file]\n", exe)
	pflag.PrintDefaults()
	os.Exit(2)
}

func main() {
	t := pflag.StringP("type", "t", "strategic", "merge|strategic")
	pflag.Parse()

	if len(pflag.Args()) != 2 {
		usage()
	}
	raw, err := getObject(pflag.Arg(0))
	if err != nil {
		panic(err)
		return
	}
	update, err := getObject(pflag.Arg(1))
	if err != nil {
		panic(err)
		return
	}
	var patch client.Patch
	switch strings.ToLower(*t) {
	case "merge":
		patch = client.MergeFrom(raw)
	case "strategic":
		patch = client.StrategicMergeFrom(raw)
	default:
		logrus.Errorf("unknown type: %s", *t)
		usage()
	}
	data, err := patch.Data(update)
	if err != nil {
		panic(err)
	}
	_, _ = os.Stdout.Write(data)
}

func getObject(filename string) (client.Object, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	meta, _, err := decUnstructured.Decode(data, nil, &metaV1.PartialObjectMetadata{})
	if err != nil {
		return nil, err
	}

	out, err := scheme.Scheme.New(meta.GetObjectKind().GroupVersionKind())
	if err != nil {
		return nil, err
	}
	out, _, err = decUnstructured.Decode(data, nil, out)
	if err != nil {
		return nil, err
	}
	return out.(client.Object), err
}
