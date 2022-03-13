package scheduler

import (
	"context"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"sigs.k8s.io/controller-runtime/pkg/client"

	schedulingv1alpha1 "github.com/moirai-io/moirai-scheduler/api/v1alpha1"

	moiraiClient "github.com/moirai-io/moirai-scheduler/pkg/client"
)

const (
	Name = "moirai"
)

type Plugin struct {
	client client.WithWatch
}

var _ framework.FilterPlugin = &Plugin{}

// Name returns name of the plugin.
func (p *Plugin) Name() string {
	return Name
}

// New initializes a new plugin and returns it.
func New(obj runtime.Object, handle framework.Handle) (framework.Plugin, error) {
	// return nil, errors.New("not implemented")
	klog.V(3).Info("Creating new moirai plugin")
	c, err := moiraiClient.NewWatchClient()
	if err != nil {
		return nil, err
	}

	go func() {
		watchInterface, err := c.Watch(context.TODO(), &schedulingv1alpha1.QueueList{}, &client.ListOptions{})
		klog.V(3).Info("Watching for queue changes")
		if err != nil {
			klog.Error(err)
			return
		}
		defer watchInterface.Stop()

		for event := range watchInterface.ResultChan() {
			klog.V(3).Info("Received event", event.Object)
		}
	}()

	return &Plugin{
		client: c,
	}, nil
}

func (p *Plugin) Filter(ctx context.Context, _ *framework.CycleState, pod *v1.Pod, node *framework.NodeInfo) *framework.Status {
	klog.V(3).Infof("filter pod: %v, node: %v\n", pod.Name, node.Node().Name)
	return framework.NewStatus(framework.Success, "")
}
