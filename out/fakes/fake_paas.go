// This file was generated by counterfeiter
package fakes

import (
	"github.com/mcb/cf-service-resource/out"
	"sync"
)

type FakePAAS struct {
	LoginStub        func(api string, username string, password string, insecure bool) error
	loginMutex       sync.RWMutex
	loginArgsForCall []struct {
		api      string
		username string
		password string
		insecure bool
	}
	loginReturns struct {
		result1 error
	}
	TargetStub        func(organization string, space string) error
	targetMutex       sync.RWMutex
	targetArgsForCall []struct {
		organization string
		space        string
	}
	targetReturns struct {
		result1 error
	}
	CreateServiceStub  func(service string, plan string, instanceName string) error
	createServiceMutex sync.RWMutex
	createServiceArgsForCall []struct {
		service      string
		plan         string
		instanceName string
	}
	createServiceReturns struct {
		result1 error
	}
	BindServiceStub  func(currentAppName string, instanceName string) error
	bindServiceMutex sync.RWMutex
	bindServiceArgsForCall []struct {
		currentAppName string
		instanceName   string
	}
	bindServiceReturns struct {
		result1 error
	}
	RestageAppStub  func(currentAppName string) error
	restageAppMutex sync.RWMutex
	restageAppArgsForCall []struct {
		currentAppName string
	}
	restageAppReturns struct {
		result1 error
	}
}

func (fake *FakePAAS) Login(api string, username string, password string, insecure bool) error {
	fake.loginMutex.Lock()
	fake.loginArgsForCall = append(fake.loginArgsForCall, struct {
		api      string
		username string
		password string
		insecure bool
	}{api, username, password, insecure})
	fake.loginMutex.Unlock()
	if fake.LoginStub != nil {
		return fake.LoginStub(api, username, password, insecure)
	} else {
		return fake.loginReturns.result1
	}
}

func (fake *FakePAAS) LoginCallCount() int {
	fake.loginMutex.RLock()
	defer fake.loginMutex.RUnlock()
	return len(fake.loginArgsForCall)
}

func (fake *FakePAAS) LoginArgsForCall(i int) (string, string, string, bool) {
	fake.loginMutex.RLock()
	defer fake.loginMutex.RUnlock()
	return fake.loginArgsForCall[i].api, fake.loginArgsForCall[i].username, fake.loginArgsForCall[i].password, fake.loginArgsForCall[i].insecure
}

func (fake *FakePAAS) LoginReturns(result1 error) {
	fake.LoginStub = nil
	fake.loginReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePAAS) Target(organization string, space string) error {
	fake.targetMutex.Lock()
	fake.targetArgsForCall = append(fake.targetArgsForCall, struct {
		organization string
		space        string
	}{organization, space})
	fake.targetMutex.Unlock()
	if fake.TargetStub != nil {
		return fake.TargetStub(organization, space)
	} else {
		return fake.targetReturns.result1
	}
}

func (fake *FakePAAS) TargetCallCount() int {
	fake.targetMutex.RLock()
	defer fake.targetMutex.RUnlock()
	return len(fake.targetArgsForCall)
}

func (fake *FakePAAS) TargetArgsForCall(i int) (string, string) {
	fake.targetMutex.RLock()
	defer fake.targetMutex.RUnlock()
	return fake.targetArgsForCall[i].organization, fake.targetArgsForCall[i].space
}

func (fake *FakePAAS) TargetReturns(result1 error) {
	fake.TargetStub = nil
	fake.targetReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePAAS) CreateService(service string, plan string, instanceName string) error {
	fake.createServiceMutex.Lock()
	fake.createServiceArgsForCall = append(fake.createServiceArgsForCall, struct {
		service      string
		plan         string
		instanceName string
	}{service, plan, instanceName})
	fake.createServiceMutex.Unlock()
	if fake.CreateServiceStub != nil {
		return fake.CreateServiceStub(service, plan, instanceName)
	} else {
		return fake.createServiceReturns.result1
	}
}

func (fake *FakePAAS) CreateServiceCallCount() int {
	fake.createServiceMutex.RLock()
	defer fake.createServiceMutex.RUnlock()
	return len(fake.createServiceArgsForCall)
}

func (fake *FakePAAS) CreateServiceArgsForCall(i int) (string, string, string) {
	fake.createServiceMutex.RLock()
	defer fake.createServiceMutex.RUnlock()
	return fake.createServiceArgsForCall[i].service, fake.createServiceArgsForCall[i].plan, fake.createServiceArgsForCall[i].instanceName
}

func (fake *FakePAAS) CreateServiceReturns(result1 error) {
	fake.CreateServiceStub = nil
	fake.createServiceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePAAS) BindService(currentAppName string, instanceName string) error {
	fake.bindServiceMutex.Lock()
	fake.bindServiceArgsForCall = append(fake.bindServiceArgsForCall, struct {
		currentAppName string
		instanceName   string
	}{currentAppName, instanceName})
	fake.bindServiceMutex.Unlock()
	if fake.BindServiceStub != nil {
		return fake.BindServiceStub(currentAppName, instanceName)
	} else {
		return fake.bindServiceReturns.result1
	}
}

func (fake *FakePAAS) BindServiceCallCount() int {
	fake.bindServiceMutex.RLock()
	defer fake.bindServiceMutex.RUnlock()
	return len(fake.bindServiceArgsForCall)
}

func (fake *FakePAAS) BindServiceArgsForCall(i int) (string, string) {
	fake.bindServiceMutex.RLock()
	defer fake.bindServiceMutex.RUnlock()
	return fake.bindServiceArgsForCall[i].currentAppName, fake.bindServiceArgsForCall[i].instanceName
}

func (fake *FakePAAS) BindServiceReturns(result1 error) {
	fake.BindServiceStub = nil
	fake.bindServiceReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePAAS) RestageApp(currentAppName string) error {
	fake.restageAppMutex.Lock()
	fake.restageAppArgsForCall = append(fake.restageAppArgsForCall, struct {
		currentAppName string
	}{currentAppName})
	fake.restageAppMutex.Unlock()
	if fake.RestageAppStub != nil {
		return fake.RestageAppStub(currentAppName)
	} else {
		return fake.restageAppReturns.result1
	}
}

func (fake *FakePAAS) RestageAppCallCount() int {
	fake.restageAppMutex.RLock()
	defer fake.restageAppMutex.RUnlock()
	return len(fake.restageAppArgsForCall)
}

func (fake *FakePAAS) RestageAppArgsForCall(i int) (string) {
	fake.restageAppMutex.RLock()
	defer fake.restageAppMutex.RUnlock()
	return fake.restageAppArgsForCall[i].currentAppName
}

func (fake *FakePAAS) RestageAppReturns(result1 error) {
	fake.RestageAppStub = nil
	fake.restageAppReturns = struct {
		result1 error
	}{result1}
}

var _ out.PAAS = new(FakePAAS)
