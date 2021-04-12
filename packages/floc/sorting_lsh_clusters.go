package floc

import (
	"errors"
)

func ApplySortingLsh(sim_hash uint64, cluster_data []byte, kMaxNumberOfBitsInFloc uint8, check_sensiveness bool) (uint64, error) {
	var kExpectedFinalCumulativeSum uint64 = (1 << kMaxNumberOfBitsInFloc);
	var kSortingLshMaxBits uint8 = 7
	var kSortingLshBlockedMask uint8 = 0b1000000
	var kSortingLshSizeMask uint8 = 0b0111111
	var cumulative_sum uint64 = 0
	var index uint64
	
	for index = 0; index < uint64(len(cluster_data)); index++ {
		// TODO implement google::protobuf::io::CodedInputStream::ReadVarint32
		next_combined := uint8(cluster_data[index])
		if (next_combined >> kSortingLshMaxBits) > 0 {
			return 0, errors.New("need implement CodedInputStream::ReadVarint32")
		}
		
		is_blocked := next_combined & kSortingLshBlockedMask
		next := next_combined & kSortingLshSizeMask
		
		if next > kMaxNumberOfBitsInFloc {
			return 0, errors.New("invalid cluster data")
		}
		
		cumulative_sum += (1 << next)

		if cumulative_sum > kExpectedFinalCumulativeSum {
			return 0, errors.New("cumulative_sum overflowed")
		}
		
		if cumulative_sum > sim_hash {
			if check_sensiveness && (is_blocked != 0) {
				return 0, errors.New("blocked")
			}
			return index, nil
		}
		
	}
	
	return 0, errors.New("index not found")
}
