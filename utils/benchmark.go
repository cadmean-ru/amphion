package utils

import (
	"fmt"
	"time"
)

// A type for benchmarking time performance of code sections.
type Benchmark struct {
	name      string
	startTime time.Time
	result    time.Duration
	started   bool
	count     int
}

func (b *Benchmark) GetName() string {
	return b.name
}

// Starts timer if not already started.
func (b *Benchmark) Start() {
	if b.started {
		return
	}

	b.startTime = time.Now()
	b.started = true
}

// Stops timer if started.
func (b *Benchmark) Stop() {
	if !b.started {
		return
	}

	b.result = time.Now().Sub(b.startTime)
	b.count += 1
	b.started = false
}

// Returns the duration of the last run timer.
func (b *Benchmark) GetResult() time.Duration {
	return b.result
}

// Returns the duration in seconds of the last run timer.
func (b *Benchmark) GetSeconds() float64 {
	return b.result.Seconds()
}

// Returns the duration in milliseconds of the last run timer.
func (b *Benchmark) GetMillis() int64 {
	return b.result.Milliseconds()
}

// Returns the number of tests run with this instance.
func (b *Benchmark) GetCount() int {
	return b.count
}

// Returns string representation.
func (b *Benchmark) ToString() string {
	return fmt.Sprintf("Benchmark %s: %dms", b.name, b.GetMillis())
}

// Prints the result of the last run.
func (b *Benchmark) Print() {
	fmt.Println(b.ToString())
}

func NewBenchmark(name string) *Benchmark {
	return &Benchmark{
		name: name,
	}
}