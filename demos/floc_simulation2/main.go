package main

import (
	"fmt"
	"github.com/shigeki/floc_simulator/packages/floc"
	"log"
	"strconv"
)

//
// simulation1 calculate diff of cohortId when only one domain difference in two histories of n domains.
// user1_history = {"example0.com", "example1.com", "example2.com", "example3.com", "example4.com"}
// user2_history = {"example0.com", "example1.com", "example2.com", "example3.com", "example5.com"}
//
// diff =  cohortId_user1 - cohortId_user2
//
var kMaxNumberOfBitsInFloc uint8 = 50

func getCohortId(domain_list []string, sorting_lsh_cluster_data []byte) (uint64, error) {
	check_sensitiveness := false
	sim_hash := floc.SimHashString(domain_list, kMaxNumberOfBitsInFloc)
	cohortId, err := floc.ApplySortingLsh(sim_hash, sorting_lsh_cluster_data, kMaxNumberOfBitsInFloc, check_sensitiveness)
	return cohortId, err
}

func main() {
	base_num_domains := 1000
	fmt.Println("# of domain, cohortId1, cohortId2, cohort diff")
	for n := 0; n <= base_num_domains; n++ {
		var domainlist1 []string
		var domainlist2 []string
		
		for i := 0; i < base_num_domains-n; i++ {
			j := strconv.Itoa(i)
			domainlist1 = append(domainlist1, "example" + j + ".com")
			domainlist2 = append(domainlist2, "example" + j + ".com")
		}
		for i := base_num_domains-n; i < base_num_domains; i++ {
			domainlist1 = append(domainlist1, "example" + strconv.Itoa(i) + ".com")
			domainlist2 = append(domainlist2, "example" + strconv.Itoa(2*base_num_domains - i - 1) + ".com")
		}
		sorting_lsh_cluster_data, err := floc.SetUpClusterData()
		if err != nil {
			log.Fatal(err)
		}
		
		cohortId1, err := getCohortId(domainlist1, sorting_lsh_cluster_data)
		if err != nil {
			log.Fatal(err)
		}
		cohortId2, err := getCohortId(domainlist2, sorting_lsh_cluster_data)
		if err != nil {
			log.Fatal(err)
		}
		diff := (int64)(cohortId1 - cohortId2)
		fmt.Println(n, ",", cohortId1, ",", cohortId2, ",", diff)
	}
}
