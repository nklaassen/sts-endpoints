package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"golang.org/x/exp/maps"

	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/service/sts"
)

// product returns the cartesion product of given lists
func product[T any](lists ...[]T) [][]T {
	result := [][]T{{}}
	for _, pool := range lists {
		var newResult [][]T
		for _, x := range result {
			for _, y := range pool {
				temp := append([]T{}, x...)
				temp = append(temp, y)
				newResult = append(newResult, temp)
			}
		}
		result = newResult
	}
	return result
}

func stsEndpoints() []string {
	nullOption := func(*endpoints.Options) {}

	// Generate all combinations of endpoint resolution options with one choice
	// from each layer.
	layers := [][]func(*endpoints.Options){
		{
			endpoints.StrictMatchingOption,
		},
		{
			nullOption,
			endpoints.STSRegionalEndpointOption,
		},
		{
			nullOption,
			endpoints.UseFIPSEndpointOption,
			endpoints.UseDualStackEndpointOption,
		},
	}
	allOptions := product(layers...)

	endpointsSet := make(map[string]struct{})
	for _, partition := range endpoints.DefaultPartitions() {
		for region := range partition.Regions() {
			for _, opts := range allOptions {
				endpoint, err := partition.EndpointFor(sts.ServiceName, region, opts...)
				if err != nil {
					// Likely there's no fips or dualstack endpoint for this region.
					continue
				}

				endpointsSet[strings.TrimPrefix(endpoint.URL, "https://")] = struct{}{}
			}
		}
	}

	endpointsSlice := maps.Keys(endpointsSet)
	sort.Strings(endpointsSlice)
	return endpointsSlice
}

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "--help" || os.Args[1] != "--go-list" || len(os.Args) > 2) {
		fmt.Printf("usage: %s [--go-list]\n", os.Args[0])
		return
	}
	formatGoList := len(os.Args) == 2 && os.Args[1] == "--go-list"

	for _, e := range stsEndpoints() {
		if formatGoList {
			fmt.Println(`"` + e + `",`)
		} else {
			fmt.Println(e)
		}
	}
}
