package main

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"log"
	"os"
	"github.com/shigeki/floc_simulator/packages/floc"
)

//
// floc_simulator is caluculate CohortId with using host lists and SortingLshClusters.
// This needs a json file of host list for history data.
//
var kMaxNumberOfBitsInFloc uint8 = 50

func main() {
	var kFlocIdMinimumHistoryDomainSizeRequired int = 7
	if len(os.Args) != 2 {
		log.Fatal("[Usage] floc_simulator host_list.json")
	}
	domain_file := os.Args[1]
	history_data, err := ioutil.ReadFile(domain_file)
	if err != nil {
		log.Fatal(err)
	}
	var host_list []string
	if err := json.Unmarshal(history_data, &host_list); err != nil {
		log.Fatal(err);
	}
	if (len(host_list) < kFlocIdMinimumHistoryDomainSizeRequired) {
		log.Fatal("FLoC needs more than %d domains. Current %d", kFlocIdMinimumHistoryDomainSizeRequired, len(host_list))
	}

	var domain_list []string
	for _, host := range host_list {
		eTLDPlusOne, err := publicsuffix.EffectiveTLDPlusOne(host)
		if (err != nil) {
			log.Fatal(err)
		}
		domain_list = append(domain_list, eTLDPlusOne)
	}
	fmt.Println("domain_list:", domain_list)

	// cluster data comes from ~/Library/Application\ Support/Google/Chrome\ Canary/Floc/1.0.6/ in MacOS
	var cluster_file = "../../Floc/1.0.6/SortingLshClusters"
	sorting_lsh_cluster_data, err := ioutil.ReadFile(cluster_file)
	if err != nil {
		log.Fatal(err)
	}

	sim_hash := floc.SimHashString(domain_list, kMaxNumberOfBitsInFloc)
	fmt.Println("sim_hash:", sim_hash)
	cohortId, err := floc.ApplySortingLsh(sim_hash, sorting_lsh_cluster_data, kMaxNumberOfBitsInFloc)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println("cohortId:", cohortId)
}
