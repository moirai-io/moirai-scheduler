package moirai

import (
	"context"
	"fmt"
	"os"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"

	configv1beta3 "github.com/moirai-io/moirai-scheduler/apis/config/v1beta3"
	"github.com/moirai-io/moirai-scheduler/pkg/internal"
	"github.com/moirai-io/moirai-scheduler/pkg/plugins/moirai/manager"
	"github.com/moirai-io/moirai-scheduler/pkg/utils"
)

const (
	// Name is the name of the plugin used in Registry and configurations.
	Name = "moirai"
)

// Plugin is a scheduling plugin for Moirai.
type Plugin struct {
	sync.RWMutex
	moiraiClient     client.Client
	moiraiCache      cache.Cache
	moiraiManager    *manager.MoiraiManager
	frameworkHandler framework.Handle
}

// queuesort
var _ framework.QueueSortPlugin = &Plugin{}

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
	klog.Info("Creating new Moirai plugin")

	// FIXME: parse args directly
	// args, _ := obj.(*configv1beta3.MoiraiArgs)
	var args configv1beta3.MoiraiArgs
	err := utils.ParsePluginArgs(obj, &args)
	if err != nil {
		return nil, err
	}
	klog.V(5).InfoS("Successfully parse args", "MoiraiArgs", args)

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
		moiraiClient:     moiraiClient,
		moiraiCache:      moiraiCache,
		moiraiManager:    moiraiManager,
		frameworkHandler: handle,
	}, nil
}

// EventsToRegister returns a series of possible events that may cause a Pod
// failed by this plugin schedulable.
func (p *Plugin) EventsToRegister() []framework.ClusterEvent {
	return []framework.ClusterEvent{
		{Resource: framework.Pod, ActionType: framework.Add},
	}
}
