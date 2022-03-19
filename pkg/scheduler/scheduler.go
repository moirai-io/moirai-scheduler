package scheduler

import (
	"github.com/moirai-io/moirai-scheduler/pkg/internal"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// Name is the name of the plugin used in Registry and configurations.
	Name = "moirai"
)

// Plugin is a scheduling plugin for Moirai.
type Plugin struct {
	moiraiClient     client.Client
	moiraiCache      cache.Cache
	frameworkHandler framework.Handle
}

var _ framework.PreFilterPlugin = &Plugin{}
var _ framework.FilterPlugin = &Plugin{}
var _ framework.PostFilterPlugin = &Plugin{}
var _ framework.PreScorePlugin = &Plugin{}
var _ framework.ScorePlugin = &Plugin{}

// Name returns name of the plugin.
func (p *Plugin) Name() string {
	return Name
}

// New initializes a new plugin and returns it.
func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	klog.V(3).Info("Creating new moirai plugin")

	moiraiClient, err := internal.NewClient()
	if err != nil {
		return nil, err
	}

	// moiraiCache, err := internal.NewCache()
	// if err != nil {
	// 	return nil, err
	// }

	// err = moiraiCache.Start(context.TODO())
	// if err != nil {
	// 	return nil, err
	// }

	// if !moiraiCache.WaitForCacheSync(context.TODO()) {
	// 	err := fmt.Errorf("failed to sync caches")
	// 	klog.Error("failed to sync caches")
	// 	return nil, err
	// }

	return &Plugin{
		moiraiClient:     moiraiClient,
		moiraiCache:      nil,
		frameworkHandler: handle,
	}, nil
}

func (p *Plugin) EventsToRegister() []framework.ClusterEvent {
	return []framework.ClusterEvent{
		{Resource: framework.Pod, ActionType: framework.All},
	}
}
