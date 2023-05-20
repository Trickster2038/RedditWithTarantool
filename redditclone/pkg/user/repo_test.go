package user

import (
	"fmt"
	"sync"
	"testing"
)

// just test Register with -race option
func TestRegisterWithRace(t *testing.T) {

	rep := NewMockMemoryRepo()

	wg := &sync.WaitGroup{}
	for i := 0; i < 2; i++ {
		username := fmt.Sprintf("user%d", i)
		wg.Add(1)
		go func() {
			defer wg.Done()
			rep.Register(username, "passwordisok")
		}()
	}
	wg.Wait()
}
