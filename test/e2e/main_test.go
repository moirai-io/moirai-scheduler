package e2e

import (
	"context"
	"os"
	"testing"

	moirai "github.com/moirai-io/moirai-scheduler/api/v1alpha1"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/klog/v2"

	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"
)

var (
	testenv env.Environment
	image   = os.Getenv("IMAGE")
)

func TestMain(m *testing.M) {
	utilruntime.Must(moirai.AddToScheme(scheme.Scheme))

	testenv = env.NewWithConfig(envconf.New())
	// Create KinD Cluster
	kindClusterName := envconf.RandomName("moirai", 16)
	namespace := envconf.RandomName("moirai-ns", 16)

	testenv.Setup(
		envfuncs.CreateKindClusterWithConfig(kindClusterName, "kindest/node:v1.22.5", "kind-config.yaml"),
		envfuncs.CreateNamespace(namespace),
		envfuncs.LoadDockerImageToCluster(kindClusterName, image),
		deployManifest(namespace),
	).Finish(
		envfuncs.DeleteNamespace(namespace),
		envfuncs.DestroyKindCluster(kindClusterName),
	)
	os.Exit(testenv.Run(m))
}

func deployManifest(namespace string) env.Func {
	return func(ctx context.Context, cfg *envconf.Config) (context.Context, error) {
		_, err := cfg.NewClient()
		if err != nil {
			klog.ErrorS(err, "Failed to create new Client")
			return ctx, err
		}

		return ctx, nil
	}
}
