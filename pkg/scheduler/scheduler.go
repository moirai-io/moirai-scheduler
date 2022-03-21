package scheduler

import (
	"context"
	"fmt"
	"os"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	moirai "github.com/moirai-io/moirai-scheduler/api/v1alpha1"
	"github.com/moirai-io/moirai-scheduler/pkg/internal"
	"github.com/moirai-io/moirai-scheduler/pkg/manager"
)

const (
	// Name is the name of the plugin used in Registry and configurations.
	Name = "moirai"
)

// Plugin is a scheduling plugin for Moirai.
type Plugin struct {
	client           client.Client
	cache            cache.Cache
	manager          *manager.MoiraiManager
	frameworkHandler framework.Handle
}

// scheduling cycle
var _ framework.PreFilterPlugin = &Plugin{}
var _ framework.FilterPlugin = &Plugin{}
var _ framework.PostFilterPlugin = &Plugin{}

var _ framework.PreScorePlugin = &Plugin{}
var _ framework.ScorePlugin = &Plugin{}

var _ framework.ReservePlugin = &Plugin{}
var _ framework.PermitPlugin = &Plugin{}

// binding cycle
var _ framework.PostBindPlugin = &Plugin{}

var _ framework.EnqueueExtensions = &Plugin{}

// Name returns name of the plugin.
func (p *Plugin) Name() string {
	return Name
}

// New initializes a new plugin and returns it.
func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	klog.Info("Creating new moirai plugin")

	moiraiClient, err := internal.NewClient()
	if err != nil {
		return nil, err
	}

	moiraiCache, err := internal.NewCache()
	if err != nil {
		return nil, err
	}

	go func(ctx context.Context) {
		if err = moiraiCache.Start(ctx); err != nil {
			os.Exit(1)
		}
	}(context.TODO())

	if !moiraiCache.WaitForCacheSync(context.TODO()) {
		err := fmt.Errorf("failed to sync caches")
		klog.Error("failed to sync caches")
		return nil, err
	}

	moiraiManager := manager.NewMoiraiManager(
		handle.ClientSet(),
		moiraiCache,
	)

	return &Plugin{
		client:           moiraiClient,
		cache:            moiraiCache,
		manager:          moiraiManager,
		frameworkHandler: handle,
	}, nil
}

// EventsToRegister returns a series of possible events that may cause a Pod
// failed by this plugin schedulable.
func (p *Plugin) EventsToRegister() []framework.ClusterEvent {
	queueGVK := moirai.GroupVersion.WithKind("queue").String()

	return []framework.ClusterEvent{
		{Resource: framework.Pod, ActionType: framework.Add},
		{Resource: framework.GVK(queueGVK), ActionType: framework.All},
	}
}
