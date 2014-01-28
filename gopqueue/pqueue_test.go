// This package provides a priority queue implementation and
// scaffold interfaces.
//
// Copyright (C) 2011 by Krzysztof Kowalik <chris@nu7hat.ch>
package pqueue

import (
	"math/rand"
	"testing"
	"time"
)

type DummyTask struct {
	priority int
}

func NewDummyTask(p int) *DummyTask {
	return &DummyTask{priority: p}
}

func (dt *DummyTask) Less(other interface{}) bool {
	return dt.priority < other.(*DummyTask).priority
}

func TestNewQueue(t *testing.T) {
	q := New(100)
	if q.Limit != 100 {
		t.Errorf("Expected to set queue limit on create")
	}
}

func TestEnqueueAndDequeue(t *testing.T) {
	q := New(0)
	for _, x := range []int{1, 3, 4, 2, 7, 3} {
		q.Enqueue(NewDummyTask(x))
	}
	if q.Len() != 6 {
		t.Errorf("Expected to enqueue all the items")
	}
	for _, x := range []int{1, 2, 3, 3, 4, 7} {
		task := q.Dequeue().(*DummyTask)
		if task.priority != x {
			t.Errorf("Expected priority to be %d, given %d", x, task.priority)
		}
	}
	if q.Len() != 0 {
		t.Errorf("Expected to dequeue all the items")
	}
}

func TestWaitForDequeue(t *testing.T) {
	q := New(0)
	dequeued := false
	go func() {
		if q.Dequeue() != nil {
			dequeued = true
		}
	}()
	<-time.After(1e9)
	q.Enqueue(NewDummyTask(1))
	<-time.After(1e2)
	if !dequeued {
		t.Errorf("Expected to wait for dequeue")
	}
}

func TestIsEmpty(t *testing.T) {
	q := New(0)
	if !q.IsEmpty() {
		t.Errorf("Expected queue to be empty")
	}
	for _, x := range []int{1, 2, 3, 4} {
		q.Enqueue(NewDummyTask(x))
	}
	if q.IsEmpty() {
		t.Errorf("Expected queue to not be empty")
	}
}

func TestLimit(t *testing.T) {
	q := New(10)
	var err error
	for i := 0; i < 20; i += 1 {
		err = q.Enqueue(NewDummyTask(i))
	}
	if err == nil || err.Error() != "Queue limit reached" {
		t.Errorf("Expected to reach queue limit")
	}
	if q.Len() != 10 {
		t.Errorf("Expected to enqueue only 10 items, %d enqueued", q.Len())
	}
}

func BenchmarkEnqueue(b *testing.B) {
	b.StopTimer()
	q := New(0)
	b.StartTimer()
	for i := 0; i < 200000; i += 1 {
		q.Enqueue(NewDummyTask(rand.Intn(10)))
	}
}

func BenchmarkMultiEnqueue(b *testing.B) {
	b.StopTimer()
	q := New(0)
	done := make(chan bool)
	b.StartTimer()
	for i := 0; i < 4; i += 1 {
		go func() {
			for j := 0; j < 50000; j += 1 {
				q.Enqueue(NewDummyTask(rand.Intn(10)))
			}
			done <- true
		}()
	}
	for i := 0; i < 4; i += 1 {
		<-done
	}
}

func BenchmarkDequeue(b *testing.B) {
	b.StopTimer()
	q := New(0)
	b.StartTimer()
	go func() {
		for i := 0; i < 200000; i += 1 {
			q.Enqueue(NewDummyTask(rand.Intn(10)))
		}
		q.Enqueue(NewDummyTask(1000000))
	}()
	for {
		task := q.Dequeue().(*DummyTask)
		if task.priority == 1000000 {
			break
		}
	}
}

func BenchmarkMultiDequeue(b *testing.B) {
	b.StopTimer()
	q := New(0)
	done := make(chan bool)
	b.StartTimer()
	go func() {
		for i := 0; i < 200000; i += 1 {
			q.Enqueue(NewDummyTask(rand.Intn(10)))
		}
		q.Enqueue(NewDummyTask(1000000))
	}()
	for i := 0; i < 4; i += 1 {
		go func() {
			for {
				task := q.Dequeue().(*DummyTask)
				if task.priority == 1000000 {
					done <- true
					break
				}
			}
		}()
	}
	<-done
}
