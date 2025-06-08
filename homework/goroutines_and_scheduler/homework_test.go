package main

import (
	"testing"

	"container/heap"

	"github.com/stretchr/testify/assert"
)

type Task struct {
	Identifier int
	Priority   int
}

type Scheduler struct {
	heap  *Heap
	tasks map[int]*taskHeap
}

func NewScheduler() Scheduler {
	h := make(Heap, 0)
	return Scheduler{
		heap:  &h,
		tasks: make(map[int]*taskHeap),
	}
}

func (s *Scheduler) AddTask(task Task) {
	if _, ok := s.tasks[task.Identifier]; ok {
		return
	}
	newTask := &taskHeap{
		task:  &task,
		index: -1,
	}
	s.tasks[task.Identifier] = newTask
	heap.Push(s.heap, newTask)
}

func (s *Scheduler) ChangeTaskPriority(taskID int, newPriority int) {
	if item, exists := s.tasks[taskID]; exists {
		item.task.Priority = newPriority
		heap.Fix(s.heap, item.index)
	}
}

func (s *Scheduler) GetTask() Task {
	if s.heap.Len() == 0 {
		return Task{}
	}

	item := heap.Pop(s.heap).(*taskHeap)
	delete(s.tasks, item.task.Identifier)
	return *item.task
}

type taskHeap struct {
	task  *Task
	index int
}

type Heap []*taskHeap

func (th Heap) Len() int           { return len(th) }
func (th Heap) Less(i, j int) bool { return th[i].task.Priority > th[j].task.Priority }
func (th Heap) Swap(i, j int) {
	th[i], th[j] = th[j], th[i]
	th[i].index = i
	th[j].index = j
}

func (th *Heap) Push(x any) {
	n := len(*th)
	item := x.(*taskHeap)
	item.index = n
	*th = append(*th, item)
}

func (th *Heap) Pop() any {
	old := *th
	n := len(old)
	item := old[n-1]
	item.index = -1
	*th = old[0 : n-1]
	return item
}
func TestTrace(t *testing.T) {
	task1 := Task{Identifier: 1, Priority: 10}
	task2 := Task{Identifier: 2, Priority: 20}
	task3 := Task{Identifier: 3, Priority: 30}
	task4 := Task{Identifier: 4, Priority: 40}
	task5 := Task{Identifier: 5, Priority: 50}

	scheduler := NewScheduler()
	scheduler.AddTask(task1)
	scheduler.AddTask(task2)
	scheduler.AddTask(task3)
	scheduler.AddTask(task4)
	scheduler.AddTask(task5)

	task := scheduler.GetTask()
	assert.Equal(t, task5, task)

	task = scheduler.GetTask()
	assert.Equal(t, task4, task)

	scheduler.ChangeTaskPriority(1, 100)

	task = scheduler.GetTask()
	assert.Equal(t, Task{Identifier: 1, Priority: 100}, task) // 100 же по идее должно быть

	task = scheduler.GetTask()
	assert.Equal(t, task3, task)
}
