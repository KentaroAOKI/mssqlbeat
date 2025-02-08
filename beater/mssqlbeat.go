package beater

import (
	"fmt"
	"sync"
	"time"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/KentaroAOKI/mssqlbeat/config"
)

// mssqlbeat configuration.
type mssqlbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
	mu     sync.Mutex
}

// New creates an instance of mssqlbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("error reading config file: %v", err)
	}

	bt := &mssqlbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts mssqlbeat.
func (bt *mssqlbeat) Run(b *beat.Beat) error {
	logp.Info("mssqlbeat is running! Hit CTRL-C to stop it.")

	var wg sync.WaitGroup
	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		// Create enabled input array.
		enabled_inputs := enabledArray(bt.config.Inputs)

		// Create input array according to number of threads.
		scheduled_inputs := chunkArray(enabled_inputs, bt.config.Threads)

		// Execute a query on the SQL server and publish the output data.
		for _, inputs := range scheduled_inputs {
			for thread_no, input := range inputs {
				wg.Add(1)
				go func(thread_no int, input *config.Input) {
					bt.PublishMssqlData(b, input, thread_no)
					wg.Done()
				}(thread_no, &input)
			}
			wg.Wait()
		}
	}
}

// Stop stops mssqlbeat.
func (bt *mssqlbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
