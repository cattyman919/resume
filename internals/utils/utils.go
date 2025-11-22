package utils

import "sync"

// Helper function to do wg.Add(1) and wg.Done() like wg.Go in the newer Go version
func Go(wg *sync.WaitGroup, f func()) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		f()
	}()
}
