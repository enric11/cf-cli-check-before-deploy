// Code generated by counterfeiter. DO NOT EDIT.
package rpcfakes

import (
	"io"
	"sync"

	"code.cloudfoundry.org/cli/plugin/rpc"
)

type FakeOutputCapture struct {
	SetOutputBucketStub        func(io.Writer)
	setOutputBucketMutex       sync.RWMutex
	setOutputBucketArgsForCall []struct {
		arg1 io.Writer
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeOutputCapture) SetOutputBucket(arg1 io.Writer) {
	fake.setOutputBucketMutex.Lock()
	fake.setOutputBucketArgsForCall = append(fake.setOutputBucketArgsForCall, struct {
		arg1 io.Writer
	}{arg1})
	fake.recordInvocation("SetOutputBucket", []interface{}{arg1})
	fake.setOutputBucketMutex.Unlock()
	if fake.SetOutputBucketStub != nil {
		fake.SetOutputBucketStub(arg1)
	}
}

func (fake *FakeOutputCapture) SetOutputBucketCallCount() int {
	fake.setOutputBucketMutex.RLock()
	defer fake.setOutputBucketMutex.RUnlock()
	return len(fake.setOutputBucketArgsForCall)
}

func (fake *FakeOutputCapture) SetOutputBucketCalls(stub func(io.Writer)) {
	fake.setOutputBucketMutex.Lock()
	defer fake.setOutputBucketMutex.Unlock()
	fake.SetOutputBucketStub = stub
}

func (fake *FakeOutputCapture) SetOutputBucketArgsForCall(i int) io.Writer {
	fake.setOutputBucketMutex.RLock()
	defer fake.setOutputBucketMutex.RUnlock()
	argsForCall := fake.setOutputBucketArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeOutputCapture) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.setOutputBucketMutex.RLock()
	defer fake.setOutputBucketMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeOutputCapture) recordInvocation(key string, args []interface{}) {
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

var _ rpc.OutputCapture = new(FakeOutputCapture)
