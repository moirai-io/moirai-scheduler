package controller

import (
	"os"
	"testing"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/e2e-framework/pkg/env"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/envfuncs"

	moirai "github.com/moirai-io/moirai-scheduler/apis/scheduling/v1alpha1"
)

var (
	testEnv          env.Environment
	image            = os.Getenv("IMAGE")
	e2eTestAssetsDir = os.Getenv("E2ETEST_ASSETS_DIR")
)

func TestMain(m *testing.M) {
	utilruntime.Must(moirai.AddToScheme(scheme.Scheme))

	testEnv = env.New()
	kindClusterName := envconf.RandomName("moirai", 16)
	namespace := envconf.RandomName("moirai-controller-ns", 16)

	testEnv.Setup(
		envfuncs.CreateKindClusterWithConfig(kindClusterName, "kindest/node:v1.24.0", "kind-config.yaml"),
		envfuncs.CreateNamespace(namespace),
		// envfuncs.LoadDockerImageToCluster(kindClusterName, image),
		envfuncs.SetupCRDs(e2eTestAssetsDir, "moirai-scheduler.crds.yaml"),
	).Finish(
		envfuncs.DeleteNamespace(namespace),
		envfuncs.TeardownCRDs(e2eTestAssetsDir, "moirai-scheduler.crds.yaml"),
		envfuncs.DestroyKindCluster(kindClusterName),
	)
	os.Exit(testEnv.Run(m))
}
