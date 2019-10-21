// Code generated by counterfeiter. DO NOT EDIT.
package v7fakes

import (
	"sync"

	"code.cloudfoundry.org/cli/actor/v7action"
	v7 "code.cloudfoundry.org/cli/command/v7"
)

type FakeSharePrivateDomainActor struct {
	SharePrivateDomainStub        func(string, string) (v7action.Warnings, error)
	sharePrivateDomainMutex       sync.RWMutex
	sharePrivateDomainArgsForCall []struct {
		arg1 string
		arg2 string
	}
	sharePrivateDomainReturns struct {
		result1 v7action.Warnings
		result2 error
	}
	sharePrivateDomainReturnsOnCall map[int]struct {
		result1 v7action.Warnings
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSharePrivateDomainActor) SharePrivateDomain(arg1 string, arg2 string) (v7action.Warnings, error) {
	fake.sharePrivateDomainMutex.Lock()
	ret, specificReturn := fake.sharePrivateDomainReturnsOnCall[len(fake.sharePrivateDomainArgsForCall)]
	fake.sharePrivateDomainArgsForCall = append(fake.sharePrivateDomainArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("SharePrivateDomain", []interface{}{arg1, arg2})
	fake.sharePrivateDomainMutex.Unlock()
	if fake.SharePrivateDomainStub != nil {
		return fake.SharePrivateDomainStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.sharePrivateDomainReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeSharePrivateDomainActor) SharePrivateDomainCallCount() int {
	fake.sharePrivateDomainMutex.RLock()
	defer fake.sharePrivateDomainMutex.RUnlock()
	return len(fake.sharePrivateDomainArgsForCall)
}

func (fake *FakeSharePrivateDomainActor) SharePrivateDomainCalls(stub func(string, string) (v7action.Warnings, error)) {
	fake.sharePrivateDomainMutex.Lock()
	defer fake.sharePrivateDomainMutex.Unlock()
	fake.SharePrivateDomainStub = stub
}

func (fake *FakeSharePrivateDomainActor) SharePrivateDomainArgsForCall(i int) (string, string) {
	fake.sharePrivateDomainMutex.RLock()
	defer fake.sharePrivateDomainMutex.RUnlock()
	argsForCall := fake.sharePrivateDomainArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeSharePrivateDomainActor) SharePrivateDomainReturns(result1 v7action.Warnings, result2 error) {
	fake.sharePrivateDomainMutex.Lock()
	defer fake.sharePrivateDomainMutex.Unlock()
	fake.SharePrivateDomainStub = nil
	fake.sharePrivateDomainReturns = struct {
		result1 v7action.Warnings
		result2 error
	}{result1, result2}
}

func (fake *FakeSharePrivateDomainActor) SharePrivateDomainReturnsOnCall(i int, result1 v7action.Warnings, result2 error) {
	fake.sharePrivateDomainMutex.Lock()
	defer fake.sharePrivateDomainMutex.Unlock()
	fake.SharePrivateDomainStub = nil
	if fake.sharePrivateDomainReturnsOnCall == nil {
		fake.sharePrivateDomainReturnsOnCall = make(map[int]struct {
			result1 v7action.Warnings
			result2 error
		})
	}
	fake.sharePrivateDomainReturnsOnCall[i] = struct {
		result1 v7action.Warnings
		result2 error
	}{result1, result2}
}

func (fake *FakeSharePrivateDomainActor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.sharePrivateDomainMutex.RLock()
	defer fake.sharePrivateDomainMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeSharePrivateDomainActor) recordInvocation(key string, args []interface{}) {
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

var _ v7.SharePrivateDomainActor = new(FakeSharePrivateDomainActor)