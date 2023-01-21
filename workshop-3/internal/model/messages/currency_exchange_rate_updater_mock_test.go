package messages

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i gitlab.ozon.dev/go/classroom-4/teachers/homework/internal/model/messages.CurrencyExchangeRateUpdater -o ./currency_exchange_rate_updater_mock_test.go -n CurrencyExchangeRateUpdaterMock

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	"time"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// CurrencyExchangeRateUpdaterMock implements CurrencyExchangeRateUpdater
type CurrencyExchangeRateUpdaterMock struct {
	t minimock.Tester

	funcUpdateCurrencyExchangeRatesOn          func(ctx context.Context, time time.Time) (err error)
	inspectFuncUpdateCurrencyExchangeRatesOn   func(ctx context.Context, time time.Time)
	afterUpdateCurrencyExchangeRatesOnCounter  uint64
	beforeUpdateCurrencyExchangeRatesOnCounter uint64
	UpdateCurrencyExchangeRatesOnMock          mCurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOn
}

// NewCurrencyExchangeRateUpdaterMock returns a mock for CurrencyExchangeRateUpdater
func NewCurrencyExchangeRateUpdaterMock(t minimock.Tester) *CurrencyExchangeRateUpdaterMock {
	m := &CurrencyExchangeRateUpdaterMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.UpdateCurrencyExchangeRatesOnMock = mCurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOn{mock: m}
	m.UpdateCurrencyExchangeRatesOnMock.callArgs = []*CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnParams{}

	return m
}

type mCurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOn struct {
	mock               *CurrencyExchangeRateUpdaterMock
	defaultExpectation *CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnExpectation
	expectations       []*CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnExpectation

	callArgs []*CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnParams
	mutex    sync.RWMutex
}

// CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnExpectation specifies expectation struct of the CurrencyExchangeRateUpdater.UpdateCurrencyExchangeRatesOn
type CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnExpectation struct {
	mock    *CurrencyExchangeRateUpdaterMock
	params  *CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnParams
	results *CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnResults
	Counter uint64
}

// CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnParams contains parameters of the CurrencyExchangeRateUpdater.UpdateCurrencyExchangeRatesOn
type CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnParams struct {
	ctx  context.Context
	time time.Time
}

// CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnResults contains results of the CurrencyExchangeRateUpdater.UpdateCurrencyExchangeRatesOn
type CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnResults struct {
	err error
}

// Expect sets up expected params for CurrencyExchangeRateUpdater.UpdateCurrencyExchangeRatesOn
func (mmUpdateCurrencyExchangeRatesOn *mCurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOn) Expect(ctx context.Context, time time.Time) *mCurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOn {
	if mmUpdateCurrencyExchangeRatesOn.mock.funcUpdateCurrencyExchangeRatesOn != nil {
		mmUpdateCurrencyExchangeRatesOn.mock.t.Fatalf("CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn mock is already set by Set")
	}

	if mmUpdateCurrencyExchangeRatesOn.defaultExpectation == nil {
		mmUpdateCurrencyExchangeRatesOn.defaultExpectation = &CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnExpectation{}
	}

	mmUpdateCurrencyExchangeRatesOn.defaultExpectation.params = &CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnParams{ctx, time}
	for _, e := range mmUpdateCurrencyExchangeRatesOn.expectations {
		if minimock.Equal(e.params, mmUpdateCurrencyExchangeRatesOn.defaultExpectation.params) {
			mmUpdateCurrencyExchangeRatesOn.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmUpdateCurrencyExchangeRatesOn.defaultExpectation.params)
		}
	}

	return mmUpdateCurrencyExchangeRatesOn
}

// Inspect accepts an inspector function that has same arguments as the CurrencyExchangeRateUpdater.UpdateCurrencyExchangeRatesOn
func (mmUpdateCurrencyExchangeRatesOn *mCurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOn) Inspect(f func(ctx context.Context, time time.Time)) *mCurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOn {
	if mmUpdateCurrencyExchangeRatesOn.mock.inspectFuncUpdateCurrencyExchangeRatesOn != nil {
		mmUpdateCurrencyExchangeRatesOn.mock.t.Fatalf("Inspect function is already set for CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn")
	}

	mmUpdateCurrencyExchangeRatesOn.mock.inspectFuncUpdateCurrencyExchangeRatesOn = f

	return mmUpdateCurrencyExchangeRatesOn
}

// Return sets up results that will be returned by CurrencyExchangeRateUpdater.UpdateCurrencyExchangeRatesOn
func (mmUpdateCurrencyExchangeRatesOn *mCurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOn) Return(err error) *CurrencyExchangeRateUpdaterMock {
	if mmUpdateCurrencyExchangeRatesOn.mock.funcUpdateCurrencyExchangeRatesOn != nil {
		mmUpdateCurrencyExchangeRatesOn.mock.t.Fatalf("CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn mock is already set by Set")
	}

	if mmUpdateCurrencyExchangeRatesOn.defaultExpectation == nil {
		mmUpdateCurrencyExchangeRatesOn.defaultExpectation = &CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnExpectation{mock: mmUpdateCurrencyExchangeRatesOn.mock}
	}
	mmUpdateCurrencyExchangeRatesOn.defaultExpectation.results = &CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnResults{err}
	return mmUpdateCurrencyExchangeRatesOn.mock
}

//Set uses given function f to mock the CurrencyExchangeRateUpdater.UpdateCurrencyExchangeRatesOn method
func (mmUpdateCurrencyExchangeRatesOn *mCurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOn) Set(f func(ctx context.Context, time time.Time) (err error)) *CurrencyExchangeRateUpdaterMock {
	if mmUpdateCurrencyExchangeRatesOn.defaultExpectation != nil {
		mmUpdateCurrencyExchangeRatesOn.mock.t.Fatalf("Default expectation is already set for the CurrencyExchangeRateUpdater.UpdateCurrencyExchangeRatesOn method")
	}

	if len(mmUpdateCurrencyExchangeRatesOn.expectations) > 0 {
		mmUpdateCurrencyExchangeRatesOn.mock.t.Fatalf("Some expectations are already set for the CurrencyExchangeRateUpdater.UpdateCurrencyExchangeRatesOn method")
	}

	mmUpdateCurrencyExchangeRatesOn.mock.funcUpdateCurrencyExchangeRatesOn = f
	return mmUpdateCurrencyExchangeRatesOn.mock
}

// When sets expectation for the CurrencyExchangeRateUpdater.UpdateCurrencyExchangeRatesOn which will trigger the result defined by the following
// Then helper
func (mmUpdateCurrencyExchangeRatesOn *mCurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOn) When(ctx context.Context, time time.Time) *CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnExpectation {
	if mmUpdateCurrencyExchangeRatesOn.mock.funcUpdateCurrencyExchangeRatesOn != nil {
		mmUpdateCurrencyExchangeRatesOn.mock.t.Fatalf("CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn mock is already set by Set")
	}

	expectation := &CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnExpectation{
		mock:   mmUpdateCurrencyExchangeRatesOn.mock,
		params: &CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnParams{ctx, time},
	}
	mmUpdateCurrencyExchangeRatesOn.expectations = append(mmUpdateCurrencyExchangeRatesOn.expectations, expectation)
	return expectation
}

// Then sets up CurrencyExchangeRateUpdater.UpdateCurrencyExchangeRatesOn return parameters for the expectation previously defined by the When method
func (e *CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnExpectation) Then(err error) *CurrencyExchangeRateUpdaterMock {
	e.results = &CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnResults{err}
	return e.mock
}

// UpdateCurrencyExchangeRatesOn implements CurrencyExchangeRateUpdater
func (mmUpdateCurrencyExchangeRatesOn *CurrencyExchangeRateUpdaterMock) UpdateCurrencyExchangeRatesOn(ctx context.Context, time time.Time) (err error) {
	mm_atomic.AddUint64(&mmUpdateCurrencyExchangeRatesOn.beforeUpdateCurrencyExchangeRatesOnCounter, 1)
	defer mm_atomic.AddUint64(&mmUpdateCurrencyExchangeRatesOn.afterUpdateCurrencyExchangeRatesOnCounter, 1)

	if mmUpdateCurrencyExchangeRatesOn.inspectFuncUpdateCurrencyExchangeRatesOn != nil {
		mmUpdateCurrencyExchangeRatesOn.inspectFuncUpdateCurrencyExchangeRatesOn(ctx, time)
	}

	mm_params := &CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnParams{ctx, time}

	// Record call args
	mmUpdateCurrencyExchangeRatesOn.UpdateCurrencyExchangeRatesOnMock.mutex.Lock()
	mmUpdateCurrencyExchangeRatesOn.UpdateCurrencyExchangeRatesOnMock.callArgs = append(mmUpdateCurrencyExchangeRatesOn.UpdateCurrencyExchangeRatesOnMock.callArgs, mm_params)
	mmUpdateCurrencyExchangeRatesOn.UpdateCurrencyExchangeRatesOnMock.mutex.Unlock()

	for _, e := range mmUpdateCurrencyExchangeRatesOn.UpdateCurrencyExchangeRatesOnMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmUpdateCurrencyExchangeRatesOn.UpdateCurrencyExchangeRatesOnMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmUpdateCurrencyExchangeRatesOn.UpdateCurrencyExchangeRatesOnMock.defaultExpectation.Counter, 1)
		mm_want := mmUpdateCurrencyExchangeRatesOn.UpdateCurrencyExchangeRatesOnMock.defaultExpectation.params
		mm_got := CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnParams{ctx, time}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmUpdateCurrencyExchangeRatesOn.t.Errorf("CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmUpdateCurrencyExchangeRatesOn.UpdateCurrencyExchangeRatesOnMock.defaultExpectation.results
		if mm_results == nil {
			mmUpdateCurrencyExchangeRatesOn.t.Fatal("No results are set for the CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn")
		}
		return (*mm_results).err
	}
	if mmUpdateCurrencyExchangeRatesOn.funcUpdateCurrencyExchangeRatesOn != nil {
		return mmUpdateCurrencyExchangeRatesOn.funcUpdateCurrencyExchangeRatesOn(ctx, time)
	}
	mmUpdateCurrencyExchangeRatesOn.t.Fatalf("Unexpected call to CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn. %v %v", ctx, time)
	return
}

// UpdateCurrencyExchangeRatesOnAfterCounter returns a count of finished CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn invocations
func (mmUpdateCurrencyExchangeRatesOn *CurrencyExchangeRateUpdaterMock) UpdateCurrencyExchangeRatesOnAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmUpdateCurrencyExchangeRatesOn.afterUpdateCurrencyExchangeRatesOnCounter)
}

// UpdateCurrencyExchangeRatesOnBeforeCounter returns a count of CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn invocations
func (mmUpdateCurrencyExchangeRatesOn *CurrencyExchangeRateUpdaterMock) UpdateCurrencyExchangeRatesOnBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmUpdateCurrencyExchangeRatesOn.beforeUpdateCurrencyExchangeRatesOnCounter)
}

// Calls returns a list of arguments used in each call to CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmUpdateCurrencyExchangeRatesOn *mCurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOn) Calls() []*CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnParams {
	mmUpdateCurrencyExchangeRatesOn.mutex.RLock()

	argCopy := make([]*CurrencyExchangeRateUpdaterMockUpdateCurrencyExchangeRatesOnParams, len(mmUpdateCurrencyExchangeRatesOn.callArgs))
	copy(argCopy, mmUpdateCurrencyExchangeRatesOn.callArgs)

	mmUpdateCurrencyExchangeRatesOn.mutex.RUnlock()

	return argCopy
}

// MinimockUpdateCurrencyExchangeRatesOnDone returns true if the count of the UpdateCurrencyExchangeRatesOn invocations corresponds
// the number of defined expectations
func (m *CurrencyExchangeRateUpdaterMock) MinimockUpdateCurrencyExchangeRatesOnDone() bool {
	for _, e := range m.UpdateCurrencyExchangeRatesOnMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.UpdateCurrencyExchangeRatesOnMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterUpdateCurrencyExchangeRatesOnCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcUpdateCurrencyExchangeRatesOn != nil && mm_atomic.LoadUint64(&m.afterUpdateCurrencyExchangeRatesOnCounter) < 1 {
		return false
	}
	return true
}

// MinimockUpdateCurrencyExchangeRatesOnInspect logs each unmet expectation
func (m *CurrencyExchangeRateUpdaterMock) MinimockUpdateCurrencyExchangeRatesOnInspect() {
	for _, e := range m.UpdateCurrencyExchangeRatesOnMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.UpdateCurrencyExchangeRatesOnMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterUpdateCurrencyExchangeRatesOnCounter) < 1 {
		if m.UpdateCurrencyExchangeRatesOnMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn")
		} else {
			m.t.Errorf("Expected call to CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn with params: %#v", *m.UpdateCurrencyExchangeRatesOnMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcUpdateCurrencyExchangeRatesOn != nil && mm_atomic.LoadUint64(&m.afterUpdateCurrencyExchangeRatesOnCounter) < 1 {
		m.t.Error("Expected call to CurrencyExchangeRateUpdaterMock.UpdateCurrencyExchangeRatesOn")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *CurrencyExchangeRateUpdaterMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockUpdateCurrencyExchangeRatesOnInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *CurrencyExchangeRateUpdaterMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *CurrencyExchangeRateUpdaterMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockUpdateCurrencyExchangeRatesOnDone()
}
