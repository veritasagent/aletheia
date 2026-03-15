package aletheia

import (
	"context"
	"sync"
)

type Watcher struct {
	events  chan MisinformationEvent
	handler *Handler

	mu      sync.RWMutex
	running bool
	cancel  context.CancelFunc
	wg      sync.WaitGroup
}

func NewWatcher(handler *Handler) *Watcher {
	return &Watcher{
		events:  make(chan MisinformationEvent, 256),
		handler: handler,
	}
}

func (w *Watcher) Emit(event MisinformationEvent) {
	w.mu.RLock()
	running := w.running
	w.mu.RUnlock()
	if !running {
		return
	}

	select {
	case w.events <- event:
	default:
		// Drop the oldest queued claim to preserve real-time feel.
		select {
		case <-w.events:
		default:
		}
		select {
		case w.events <- event:
		default:
		}
	}
}

func (w *Watcher) Start() {
	w.mu.Lock()
	if w.running {
		w.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	w.cancel = cancel
	w.running = true
	w.mu.Unlock()

	w.wg.Add(1)
	go w.loop(ctx)
}

func (w *Watcher) Stop() {
	w.mu.Lock()
	if !w.running {
		w.mu.Unlock()
		return
	}
	w.running = false
	cancel := w.cancel
	w.cancel = nil
	w.mu.Unlock()

	if cancel != nil {
		cancel()
	}
	w.wg.Wait()
}

func (w *Watcher) loop(ctx context.Context) {
	defer w.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-w.events:
			classified := Classify(event)
			score, label := Score(classified)
			classified.RiskScore = score
			classified.RiskLabel = label
			w.handler.Handle(classified)
		}
	}
}
