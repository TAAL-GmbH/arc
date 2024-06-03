package k8s_watcher

import (
	"context"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/bitcoin-sv/arc/pkg/blocktx"
	"github.com/bitcoin-sv/arc/pkg/metamorph"
)

const (
	logLevelDefault  = slog.LevelInfo
	metamorphService = "metamorph"
	blocktxService   = "blocktx"
	intervalDefault  = 15 * time.Second
)

type K8sClient interface {
	GetRunningPodNames(ctx context.Context, namespace string, service string) (map[string]struct{}, error)
}

type Watcher struct {
	metamorphClient   metamorph.TransactionMaintainer
	blocktxClient     blocktx.BlocktxClient
	k8sClient         K8sClient
	logger            *slog.Logger
	tickerMetamorph   Ticker
	tickerBlocktx     Ticker
	namespace         string
	waitGroup         *sync.WaitGroup
	shutdownMetamorph context.CancelFunc
	shutdownBlocktx   context.CancelFunc
}

func WithLogger(logger *slog.Logger) func(*Watcher) {
	return func(p *Watcher) {
		p.logger = logger.With(slog.String("service", "k8s-watcher"))
	}
}

type ServerOption func(f *Watcher)

// New The K8s watcher listens to events coming from Kubernetes. If it detects a metamorph pod which was terminated, then it sets records locked by this pod to unlocked. This is a safety measure for the case that metamorph is terminated ungracefully where it misses to unlock its records itself.
func New(metamorphClient metamorph.TransactionMaintainer, blocktxClient blocktx.BlocktxClient, k8sClient K8sClient, namespace string, opts ...ServerOption) *Watcher {
	watcher := &Watcher{
		metamorphClient: metamorphClient,
		blocktxClient:   blocktxClient,
		k8sClient:       k8sClient,

		namespace:       namespace,
		logger:          slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: logLevelDefault})).With(slog.String("service", "k8s-watcher")),
		tickerMetamorph: NewDefaultTicker(intervalDefault),
		tickerBlocktx:   NewDefaultTicker(intervalDefault),
		waitGroup:       &sync.WaitGroup{},
	}
	for _, opt := range opts {
		opt(watcher)
	}

	return watcher
}

type Ticker interface {
	Stop()
	Tick() <-chan time.Time
}

type DefaultTicker struct {
	*time.Ticker
}

func NewDefaultTicker(d time.Duration) DefaultTicker {
	return DefaultTicker{time.NewTicker(d)}
}

func (d DefaultTicker) Tick() <-chan time.Time {
	return d.C
}

func WithMetamorphTicker(t Ticker) func(*Watcher) {
	return func(p *Watcher) {
		p.tickerMetamorph = t
	}
}

func WithBlocktxTicker(t Ticker) func(*Watcher) {
	return func(p *Watcher) {
		p.tickerBlocktx = t
	}
}

func (c *Watcher) Start() error {
	c.watchMetamorph()

	c.watchBlocktx()

	return nil
}

func (c *Watcher) watchBlocktx() {
	ctx, cancel := context.WithCancel(context.Background())
	c.shutdownBlocktx = cancel
	c.waitGroup.Add(1)
	go func() {
		var runningPods map[string]struct{}
		defer c.waitGroup.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case <-c.tickerBlocktx.Tick():
				// Update the list of running pods. Detect those which have been terminated and call them to delete unfinished block processing
				ctx := context.Background()
				runningPodsK8s, err := c.k8sClient.GetRunningPodNames(ctx, c.namespace, blocktxService)
				if err != nil {
					c.logger.Error("failed to get pods", slog.String("err", err.Error()))
					continue
				}

				for podName := range runningPods {
					// Ignore all other services than blocktx
					if !strings.Contains(podName, blocktxService) {
						continue
					}

					_, found := runningPodsK8s[podName]
					if !found {
						// A previously running pod has been terminated => set records locked by this pod unlocked
						err = c.blocktxClient.DelUnfinishedBlockProcessing(ctx, podName)
						if err != nil {
							c.logger.Error("failed to delete unfinished block processing", slog.String("pod-name", podName), slog.String("err", err.Error()))
							continue
						}
						c.logger.Info("Deleted unfinished block processing", slog.String("pod-name", podName))
					}
				}

				runningPods = runningPodsK8s
			}
		}
	}()
}

func (c *Watcher) watchMetamorph() {
	ctx, cancel := context.WithCancel(context.Background())
	c.shutdownMetamorph = cancel
	c.waitGroup.Add(1)
	go func() {
		defer c.waitGroup.Done()
		var runningPods map[string]struct{}

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.tickerMetamorph.Tick():
				// Update the list of running pods. Detect those which have been terminated and unlock records for these pods
				ctx := context.Background()
				runningPodsK8s, err := c.k8sClient.GetRunningPodNames(ctx, c.namespace, metamorphService)
				if err != nil {
					c.logger.Error("failed to get pods", slog.String("err", err.Error()))
					continue
				}

				for podName := range runningPods {
					// Ignore all other services than metamorph
					if !strings.Contains(podName, metamorphService) {
						continue
					}

					_, found := runningPodsK8s[podName]
					if !found {
						// A previously running pod has been terminated => set records locked by this pod unlocked
						resp, err := c.metamorphClient.SetUnlockedByName(ctx, podName)
						if err != nil {
							c.logger.Error("failed to unlock metamorph records", slog.String("pod-name", podName), slog.String("err", err.Error()))
							continue
						}

						c.logger.Info("records unlocked", slog.Int64("rows-affected", resp), slog.String("pod-name", podName))
					}
				}

				runningPods = runningPodsK8s
			}
		}
	}()
}

func (c *Watcher) Shutdown() {
	if c.shutdownMetamorph != nil {
		c.shutdownMetamorph()
	}

	if c.shutdownBlocktx != nil {
		c.shutdownBlocktx()
	}

	c.waitGroup.Wait()
}
