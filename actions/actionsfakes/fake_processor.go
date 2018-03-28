// Code generated by counterfeiter. DO NOT EDIT.
package actionsfakes

import (
	"sync"

	"github.com/pivotal-cf/reconfigure-pipeline/actions"
)

type FakeProcessor struct {
	ProcessStub        func(config string) string
	processMutex       sync.RWMutex
	processArgsForCall []struct {
		config string
	}
	processReturns struct {
		result1 string
	}
	processReturnsOnCall map[int]struct {
		result1 string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeProcessor) Process(config string) string {
	fake.processMutex.Lock()
	ret, specificReturn := fake.processReturnsOnCall[len(fake.processArgsForCall)]
	fake.processArgsForCall = append(fake.processArgsForCall, struct {
		config string
	}{config})
	fake.recordInvocation("Process", []interface{}{config})
	fake.processMutex.Unlock()
	if fake.ProcessStub != nil {
		return fake.ProcessStub(config)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.processReturns.result1
}

func (fake *FakeProcessor) ProcessCallCount() int {
	fake.processMutex.RLock()
	defer fake.processMutex.RUnlock()
	return len(fake.processArgsForCall)
}

func (fake *FakeProcessor) ProcessArgsForCall(i int) string {
	fake.processMutex.RLock()
	defer fake.processMutex.RUnlock()
	return fake.processArgsForCall[i].config
}

func (fake *FakeProcessor) ProcessReturns(result1 string) {
	fake.ProcessStub = nil
	fake.processReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeProcessor) ProcessReturnsOnCall(i int, result1 string) {
	fake.ProcessStub = nil
	if fake.processReturnsOnCall == nil {
		fake.processReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.processReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeProcessor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.processMutex.RLock()
	defer fake.processMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeProcessor) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ actions.Processor = new(FakeProcessor)
