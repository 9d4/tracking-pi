package log

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"sync"
	"time"
)

const defaultMaxAttempt = uint(3)

var q *Queue

func init() {
	q = New()
}

type status struct {
	attempted uint
}

// Queue is log processing queue
type Queue struct {
	m          map[primitive.ObjectID]*status
	maxAttempt uint
	every      time.Duration
	mu         sync.Mutex
	wg         sync.WaitGroup
}

func New() *Queue {
	return &Queue{
		m:          map[primitive.ObjectID]*status{},
		every:      time.Second * 5,
		maxAttempt: 3,
	}
}

func (q *Queue) Add(logID primitive.ObjectID, immediate bool) {
	q.mu.Lock()
	q.m[logID] = &status{}
	q.mu.Unlock()
	if immediate {
		q.wg.Add(1)
		go q.run(logID)
	}
}

func (q *Queue) run(logID primitive.ObjectID) {
	fmt.Println("Queue:TryRun:", logID.Hex())
	if q.m[logID].attempted >= q.maxAttempt {
		q.mu.Lock()
		delete(q.m, logID)
		q.mu.Unlock()
		q.wg.Done()
		return
	}

	fmt.Println("Queue:Processing:", logID.Hex())
	q.m[logID].attempted++
	ProcessLogResult(logID)
	q.wg.Done()
	fmt.Println("Queue:Processed:", logID.Hex())
}

// Run should be called once in go routine
func (q *Queue) Run() {
	for {
		<-time.Tick(q.every)

		for id, _ := range q.m {
			q.wg.Add(1)
			go q.run(id)
		}

		q.wg.Wait()
	}
}

func GetQueue() *Queue {
	return q
}
