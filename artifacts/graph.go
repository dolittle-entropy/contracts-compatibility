package artifacts

import (
	"fmt"
	"os"
	"sync"
)

// Graph represents a structure containing released versions of the Runtime and SDKs with their corresponding Contracts dependency
type Graph struct {
	Runtime Releases            `json:"runtime"`
	SDKs    map[string]Releases `json:"sdk"`
}

// CreateGraphFor creates a new Graph by resolving all releases for the provided Runtime and SDKs ReleaseListResolver
func CreateGraphFor(runtime *ReleaseListResolver, sdks map[string]*ReleaseListResolver) *Graph {
	graph := &Graph{
		SDKs: make(map[string]Releases),
	}

	wg := sync.WaitGroup{}
	lock := sync.Mutex{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		releases, err := runtime.ListAndResolve()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to list Runtime releases", err)
			return
		}

		lock.Lock()
		defer lock.Unlock()
		graph.Runtime = releases
	}()

	for sdk, resolver := range sdks {
		wg.Add(1)
		go func(sdk string, resolver *ReleaseListResolver) {
			defer wg.Done()
			releases, err := resolver.ListAndResolve()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to list", sdk, "SDK releases", err)
				return
			}

			lock.Lock()
			defer lock.Unlock()
			graph.SDKs[sdk] = releases
		}(sdk, resolver)
	}

	wg.Wait()
	return graph
}
