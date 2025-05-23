package main

import (
	"context"
	"crypto/rand"
	"flag"
	"fmt"
	"sort"
	"time"

	"github.com/mailgun/groupcache/v2"
	"github.com/sirupsen/logrus"
)

type contextKey string

const valueLengthKey contextKey = "valueLength"

var (
	logLevel string
	logger   = logrus.New()
	group    *groupcache.Group
)

func main() {
	batchSize := flag.Int("batch-size", 100000, "number of operations per batch")
	batchCount := flag.Int("batch-count", 20, "number of times to repeat the batch")
	valueLength := flag.Int("value-length", 1024, "length of random bytes for each value")
	cacheSize := flag.Int("cache-size", 64, "cache size in megabytes")
	flag.StringVar(&logLevel, "log-level", "info", "log level (debug, info, warn, error)")
	flag.Parse()

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logger.Fatalf("Invalid log level: %v", err)
	}
	logger.SetLevel(level)

	group = groupcache.NewGroup("cache1", int64(*cacheSize)<<20, groupcache.GetterFunc(
		func(ctx context.Context, key string, dest groupcache.Sink) error {
			logger.Debug("Get Called, cache miss")
			valueLen := ctx.Value(valueLengthKey).(int)
			value := make([]byte, valueLen)
			if _, err := rand.Read(value); err != nil {
				logger.Errorf("Failed to generate random value - %v", err)
				return err
			}
			if err := dest.SetBytes(value, time.Now().Add(10*time.Minute)); err != nil {
				logger.Errorf("Failed to set cache value for key '%s' - %v", key, err)
				return err
			}
			return nil
		},
	))

	type batchResult struct {
		number int
		time   time.Duration
	}
	batchTimes := make([]batchResult, 0, *batchCount)
	start := time.Now()
	for batchNum := 1; batchNum <= *batchCount; batchNum++ {
		logger.Infof("Starting batch %d/%d", batchNum, *batchCount)
		batchStart := time.Now()
		for i := 1; i <= *batchSize; i++ {
			var value groupcache.ByteView
			ctx := context.WithValue(context.Background(), valueLengthKey, *valueLength)
			group.Get(ctx, fmt.Sprintf("key-%d", i), groupcache.ByteViewSink(&value))
			logger.Debugf("Value length: %d", value.Len())
		}
		batchTime := time.Since(batchStart)
		batchTimes = append(batchTimes, batchResult{batchNum, batchTime})
		logger.Infof("Batch %d/%d completed in %s", batchNum, *batchCount, batchTime)
	}
	totalTime := time.Since(start)

	sort.Slice(batchTimes, func(i, j int) bool {
		return batchTimes[i].time < batchTimes[j].time
	})

	fmt.Printf("\nTiming Statistics:\n")
	fmt.Printf("Total time for %d batches: %s\n", *batchCount, totalTime)
	fmt.Printf("Fastest batch: #%d (%s)\n", batchTimes[0].number, batchTimes[0].time)
	fmt.Printf("Slowest batch: #%d (%s)\n", batchTimes[len(batchTimes)-1].number, batchTimes[len(batchTimes)-1].time)
	fmt.Printf("25th percentile: %s\n", batchTimes[len(batchTimes)/4].time)
	fmt.Printf("50th percentile: %s\n", batchTimes[len(batchTimes)/2].time)
	fmt.Printf("75th percentile: %s\n", batchTimes[len(batchTimes)*3/4].time)

	fmt.Printf("\nGroupCache Statistics:\n")
	fmt.Printf("Gets: %d\n", group.Stats.Gets.Get())
	fmt.Printf("CacheHits: %d\n", group.Stats.CacheHits.Get())
	fmt.Printf("PeerLoads: %d\n", group.Stats.PeerLoads.Get())
	fmt.Printf("PeerErrors: %d\n", group.Stats.PeerErrors.Get())
	fmt.Printf("Loads: %d\n", group.Stats.Loads.Get())
	fmt.Printf("LoadsDeduped: %d\n", group.Stats.LoadsDeduped.Get())
	fmt.Printf("LocalLoads: %d\n", group.Stats.LocalLoads.Get())
	fmt.Printf("LocalLoadErrs: %d\n", group.Stats.LocalLoadErrs.Get())
	fmt.Printf("ServerRequests: %d\n", group.Stats.ServerRequests.Get())

	mainCache := group.CacheStats(groupcache.MainCache)
	fmt.Printf("\nMain Cache Statistics:\n")
	fmt.Printf("Items: %d\n", mainCache.Items)
	fmt.Printf("Bytes: %d\n", mainCache.Bytes)
	fmt.Printf("Gets: %d\n", mainCache.Gets)
	fmt.Printf("Hits: %d\n", mainCache.Hits)
	fmt.Printf("Evictions: %d\n", mainCache.Evictions)

	hotCache := group.CacheStats(groupcache.HotCache)
	fmt.Printf("\nHot Cache Statistics:\n")
	fmt.Printf("Items: %d\n", hotCache.Items)
	fmt.Printf("Bytes: %d\n", hotCache.Bytes)
	fmt.Printf("Gets: %d\n", hotCache.Gets)
	fmt.Printf("Hits: %d\n", hotCache.Hits)
	fmt.Printf("Evictions: %d\n", hotCache.Evictions)
}
