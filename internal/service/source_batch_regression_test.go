package service

import (
	"reflect"
	"testing"

	"texttool/internal/config"
	"texttool/internal/resource"
	"texttool/internal/store"
	"texttool/pkg/logger"
)

func TestTextBatchReleasesEachSourceBeforeOpeningNext(t *testing.T) {
	svc := New(store.NewMemoryStore(), logger.NewLevel(logger.LevelError), &config.Config{})
	factory := resource.NewFactory(2, map[string]string{"a": "alpha", "b": "beta", "c": "gamma"})
	got, err := svc.ReadTextBatch(factory, []string{"a", "b", "c"})
	if err != nil || !reflect.DeepEqual(got, []string{"alpha", "beta", "gamma"}) {
		t.Fatalf("valid text batch exhausted source handles before the batch finished: got=%v err=%v open=%d", got, err, factory.OpenCount())
	}
	if factory.OpenCount() != 0 {
		t.Fatalf("completed text batch left source handles open: %d", factory.OpenCount())
	}
	if _, err := svc.ReadTextBatch(factory, []string{"a", "missing"}); err == nil {
		t.Fatal("missing source unexpectedly succeeded")
	}
	if factory.OpenCount() != 0 {
		t.Fatalf("failed text batch leaked source capacity into the next batch: %d", factory.OpenCount())
	}
}
