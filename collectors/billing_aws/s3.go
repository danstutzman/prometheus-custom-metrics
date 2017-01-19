package billing_aws

import (
	"github.com/danielstutzman/prometheus-custom-metrics/storage/s3"
	"io/ioutil"
	"log"
	"regexp"
	"time"
)

// returns "" if can't find any match
func listBestPathForThisMonth(s3 *s3.S3Connection) string {
	todayPaths := listPathsForThisMonth(s3, 0)
	if len(todayPaths) == 0 {
		ydayPaths := listPathsForThisMonth(s3, -1)
		if len(ydayPaths) == 0 {
			return ""
		} else if len(ydayPaths) == 1 {
			return ydayPaths[0]
		} else {
			log.Fatalf("Found >1 billing paths for yday's month: %v", todayPaths)
			return "" // unreachable
		}
	} else if len(todayPaths) == 1 {
		return todayPaths[0]
	} else {
		log.Fatalf("Found >1 billing paths for today's month: %v", todayPaths)
		return "" // unreachable
	}
}

func listPathsForThisMonth(s3 *s3.S3Connection, dayOffset int) []string {
	paths := []string{}
	yyyymm := regexp.MustCompile(
		"[0-9]+-aws-billing-detailed-line-items-with-resources-and-tags-" +
			time.Now().UTC().AddDate(0, 0, dayOffset).Format("2006-01") + "\\.csv\\.zip")
	for _, path := range s3.ListPaths() {
		if yyyymm.MatchString(path) {
			paths = append(paths, path)
		}
	}
	return paths
}

func (collector *BillingAwsCollector) downloadBucketNameToNumBytes() map[string]float64 {
	bestPath := listBestPathForThisMonth(collector.s3)
	if bestPath == "" {
		log.Fatalf("Couldn't find billing path for today")
	}

	zipDownloader := collector.s3.DownloadPath(bestPath)
	zipBytes, err := ioutil.ReadAll(zipDownloader)
	if err != nil {
		log.Fatalf("Error from ioutil.ReadAll of %s: %s", bestPath, err)
	}

	return readBucketSizesFromZip(zipBytes)
}
