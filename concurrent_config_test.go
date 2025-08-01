package goconfig_test

import (
	"github.com/gojekfarm/goconfig/v2"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestConcurrentGetValue(t *testing.T) {
	baseConfig := goconfig.NewBaseConfig()
	baseConfig.Load()

	const numGoroutines = 100
	const numOperations = 1000
	var wg sync.WaitGroup

	// Test concurrent access to the same key
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				value := baseConfig.GetValue("foo")
				if value != "bar" {
					t.Errorf("Expected 'bar', got '%s'", value)
				}
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentGetValueDifferentKeys(t *testing.T) {
	baseConfig := goconfig.NewBaseConfig()
	baseConfig.Load()

	const numGoroutines = 50
	var wg sync.WaitGroup

	// Test concurrent access to different keys
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			if goroutineID%2 == 0 {
				value := baseConfig.GetValue("foo")
				if value != "bar" {
					t.Errorf("Expected 'bar', got '%s'", value)
				}
			} else {
				value := baseConfig.GetValue("new_relic_app_name")
				if value != "foo" {
					t.Errorf("Expected 'foo', got '%s'", value)
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestConcurrentGetIntValue(t *testing.T) {
	baseConfig := goconfig.NewBaseConfig()
	baseConfig.Load()

	const numGoroutines = 100
	const numOperations = 500
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				value := baseConfig.GetIntValue("someInt")
				if value != 1 {
					t.Errorf("Expected 1, got %d", value)
				}
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentGetOptionalValue(t *testing.T) {
	baseConfig := goconfig.NewBaseConfig()
	baseConfig.Load()

	const numGoroutines = 50
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Test existing key
			value := baseConfig.GetOptionalValue("foo", "default")
			if value != "bar" {
				t.Errorf("Expected 'bar', got '%s'", value)
			}

			// Test non-existing key with default
			value = baseConfig.GetOptionalValue("nonexistent", "default")
			if value != "default" {
				t.Errorf("Expected 'default', got '%s'", value)
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentGetOptionalIntValue(t *testing.T) {
	baseConfig := goconfig.NewBaseConfig()
	baseConfig.Load()

	const numGoroutines = 50
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Test existing key
			value := baseConfig.GetOptionalIntValue("someInt", 999)
			if value != 1 {
				t.Errorf("Expected 1, got %d", value)
			}

			// Test non-existing key with default
			value = baseConfig.GetOptionalIntValue("nonexistentInt", 999)
			if value != 999 {
				t.Errorf("Expected 999, got %d", value)
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentGetFeature(t *testing.T) {
	baseConfig := goconfig.NewBaseConfig()
	baseConfig.Load()

	const numGoroutines = 50
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Test feature that is true
			value := baseConfig.GetFeature("someFeature")
			if !value {
				t.Errorf("Expected true, got %v", value)
			}

			// Test feature that is false
			value = baseConfig.GetFeature("someOtherFeature")
			if value {
				t.Errorf("Expected false, got %v", value)
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentMixedOperations(t *testing.T) {
	baseConfig := goconfig.NewBaseConfig()
	baseConfig.Load()

	const numGoroutines = 20
	const numOperations = 100
	var wg sync.WaitGroup

	// Test mixed concurrent operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				switch goroutineID % 5 {
				case 0:
					baseConfig.GetValue("foo")
				case 1:
					baseConfig.GetIntValue("someInt")
				case 2:
					baseConfig.GetOptionalValue("foo", "default")
				case 3:
					baseConfig.GetOptionalIntValue("someInt", 999)
				case 4:
					baseConfig.GetFeature("someFeature")
				}
			}
		}(i)
	}

	wg.Wait()
}

func TestNoDeadlockScenario(t *testing.T) {
	baseConfig := goconfig.NewBaseConfig()
	baseConfig.Load()

	const numGoroutines = 100
	var wg sync.WaitGroup
	timeout := time.After(10 * time.Second)
	done := make(chan bool)

	// Start goroutines
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Rapidly access different methods that use the same mutex
			for j := 0; j < 100; j++ {
				baseConfig.GetValue("foo")
				baseConfig.GetIntValue("someInt")
				baseConfig.GetOptionalValue("nonexistent", "default")
				baseConfig.GetOptionalIntValue("nonexistentInt", 999)
				baseConfig.GetFeature("someFeature")
			}
		}()
	}

	// Wait for completion or timeout
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		// Test completed successfully, no deadlock
	case <-timeout:
		t.Fatal("Test timed out - possible deadlock detected")
	}
}

func TestRaceConditionDetection(t *testing.T) {
	baseConfig := goconfig.NewBaseConfig()
	baseConfig.Load()

	const numGoroutines = 50
	var wg sync.WaitGroup

	// Create high contention scenario
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				// All goroutines access the same key to create contention
				baseConfig.GetValue("foo")
				runtime.Gosched() // Yield to increase chance of race conditions
			}
		}()
	}

	wg.Wait()
}

func TestConcurrentCacheEviction(t *testing.T) {
	// Test scenario where cache might be cleared while being accessed
	baseConfig := goconfig.BaseConfig{}
	baseConfig.Load()

	const numGoroutines = 20
	var wg sync.WaitGroup

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			if goroutineID == 0 {
				// One goroutine periodically clears cache
				for j := 0; j < 10; j++ {
					time.Sleep(50 * time.Millisecond)
					baseConfig.Load()
				}
			} else {
				// Other goroutines continuously access values
				for j := 0; j < 100; j++ {
					baseConfig.GetValue("foo")
					time.Sleep(10 * time.Millisecond)
				}
			}
		}(i)
	}

	wg.Wait()
}
