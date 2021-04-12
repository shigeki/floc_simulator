package floc

import (
	"encoding/json"
	"errors"
	"golang.org/x/net/publicsuffix"
	"io/ioutil"
	"os"
)

var kFlocIdMinimumHistoryDomainSizeRequired int = 7

func SetUp() ([]string, []byte, error) {
	var domain_list []string
	var sorting_lsh_cluster_data []byte
	
	if len(os.Args) != 2 {
		err := errors.New("[Usage] floc_simulator host_list.json")
		return domain_list, sorting_lsh_cluster_data, err
	}
	domain_file := os.Args[1]
	history_data, err := ioutil.ReadFile(domain_file)
	if err != nil {
		return domain_list, sorting_lsh_cluster_data, err
	}
	var host_list []string
	if err := json.Unmarshal(history_data, &host_list); err != nil {
		return domain_list, sorting_lsh_cluster_data, err
	}
	if (len(host_list) < kFlocIdMinimumHistoryDomainSizeRequired) {
		err := errors.New("Too small host list. Need more than 7 domains.")
		return domain_list, sorting_lsh_cluster_data, err
	}

	for _, host := range host_list {
		eTLDPlusOne, err := publicsuffix.EffectiveTLDPlusOne(host)
		if (err != nil) {
			return domain_list, sorting_lsh_cluster_data, err
		}
		domain_list = append(domain_list, eTLDPlusOne)
	}

	// cluster data comes from ~/Library/Application\ Support/Google/Chrome\ Canary/Floc/1.0.6/ in MacOS
	var cluster_file = "../../Floc/1.0.6/SortingLshClusters"
	sorting_lsh_cluster_data, err = ioutil.ReadFile(cluster_file)
	if err != nil {
		return domain_list, sorting_lsh_cluster_data, err
	}
	return domain_list, sorting_lsh_cluster_data, nil
}
