package controller

import (
	"context"
	"os"
	"testing"

	"sigs.k8s.io/e2e-framework/klient/decoder"
	"sigs.k8s.io/e2e-framework/klient/k8s/resources"
	"sigs.k8s.io/e2e-framework/pkg/envconf"
	"sigs.k8s.io/e2e-framework/pkg/features"

	moirai "github.com/rudeigerc/moirai/apis/scheduling/v1alpha1"
)

func TestCRDSetup(t *testing.T) {
	feature := features.New("Custom Controller").
		Setup(func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			r, err := resources.New(c.Client().RESTConfig())
			if err != nil {
				t.Fail()
			}
			moirai.AddToScheme(r.GetScheme())
			r.WithNamespace(c.Namespace())
			decoder.DecodeEachFile(
				ctx, os.DirFS("testdata"), "scheduling_v1alpha1_queue.yaml",
				decoder.CreateHandler(r),
				decoder.MutateNamespace(c.Namespace()),
			)
			return ctx
		}).
		Assess("Check If Resource created", func(ctx context.Context, t *testing.T, c *envconf.Config) context.Context {
			r, err := resources.New(c.Client().RESTConfig())
			if err != nil {
				t.Fail()
			}
			r.WithNamespace(c.Namespace())
			moirai.AddToScheme(r.GetScheme())
			queue := &moirai.Queue{}
			err = r.Get(ctx, "sample-queue", c.Namespace(), queue)
			if err != nil {
				t.Fail()
			}
			return ctx
		}).Feature()

	testEnv.Test(t, feature)
}
