package set_ttls

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

// WorkTuple defines the object that contains work
type WorkTuple func() (uint64, []string)

// ProcessedTuple defines the object that contains info about a succfully processed batch
type ProcessedTuple func() (uint64, int64)

// Updater contains the logic for updating ttls.
type Updater struct {
	*Config
	wgAll     sync.WaitGroup
	wgTTL     sync.WaitGroup
	rdb       *redis.Client
	done      chan struct{}
	abort     chan error
	work      chan WorkTuple
	processed chan ProcessedTuple
}

// SetTTLs sets ttls.
func (u *Updater) SetTTLs(ctx context.Context) {
	u.wgTTL.Add(u.Workers)
	u.wgAll.Add(u.Workers + 2)

	go u.waitWorker(ctx)
	go u.logWorker(ctx)
	go u.dispatchWorker(ctx)
	for i := 0; i < u.Workers; i++ {
		go u.ttlWorker(ctx, i)
	}

	signals := make(chan os.Signal)
	signal.Notify(signals, os.Interrupt)
	signal.Notify(signals, os.Kill)

	select {
	case sig := <-signals:
		log.Fatalln("Terminating due to signal:", sig)
	case <-u.done:
		log.Println("Procesing complete")
	}
}

func (u *Updater) waitWorker(ctx context.Context) {
	u.wgTTL.Wait()
	close(u.processed)
	u.wgAll.Wait()
	close(u.done)
}

func (u *Updater) logWorker(ctx context.Context) {
	defer u.wgAll.Done()

	ticker := time.NewTicker(u.LogPeriod)
	defer ticker.Stop()

	var cursor uint64
	var count int64
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			log.Println("Last offset:", cursor, "Total keys updated:", count)
		case task, more := <-u.processed:
			if !more {
				log.Println("Total keys updated:", count)
				return
			}

			tCursor, keyCount := task()
			cursor = tCursor
			count += keyCount
		}
	}
}

func (u *Updater) dispatchWorker(ctx context.Context) {
	defer u.wgAll.Done()
	defer close(u.work)

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	cursor := u.Cursor
	rate := time.Second / time.Duration(u.RateLimit)
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-u.abort:
			log.Println("Error occurred while process a batch of keys:", err)
			return
		default:
			start := time.Now()
			keys, tCursor, err := u.rdb.Scan(ctx, cursor, u.Match, u.BatchSize).Result()
			if err != nil {
				log.Println("Could not read next batch of keys:", err)
			}

			if len(keys) != 0 {
				u.work <- deriveTupleWork(tCursor, keys)

				elapsed := time.Now().Sub(start)
				timeToTake := rate * time.Duration(len(keys))
				if timeToTake > elapsed {
					time.Sleep(timeToTake - elapsed)
				}
			}

			if tCursor == 0 {
				log.Println("All offsets dispatched for processing.")
				return
			}
			cursor = tCursor
		}
	}
}

func (u *Updater) ttlWorker(ctx context.Context, id int) {
	defer u.wgAll.Done()
	defer u.wgTTL.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case task, more := <-u.work:
			if !more {
				return
			}

			cursor, keys := task()
			for _, key := range keys {
				if _, err := u.rdb.Expire(ctx, key, u.TTL).Result(); err != nil {
					u.abort <- err
					return
				}
			}

			u.processed <- deriveTupleProcessed(cursor, int64(len(keys)))
		}
	}
}

// NewUpdater returns a new Updater
func NewUpdater(cfg *Config, rdb *redis.Client) *Updater {
	return &Updater{
		Config:    cfg,
		rdb:       rdb,
		done:      make(chan struct{}),
		abort:     make(chan error),
		work:      make(chan WorkTuple),
		processed: make(chan ProcessedTuple),
	}
}
