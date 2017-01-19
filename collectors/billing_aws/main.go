package billing_aws

import (
	"archive/zip"
	"bytes"
	"encoding/csv"
	"github.com/danielstutzman/prometheus-custom-metrics/storage/bigquery"
	"github.com/danielstutzman/prometheus-custom-metrics/storage/s3"
	"io"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"time"
)

const GIBIBYTES_TO_BYTES = 1024 * 1024 * 1024
const MIN_BUCKET_BYTES_TO_HAVE_OWN_METRIC = 10 * 1024 * 1024
const OTHER_BUCKET_NAME = "other"

var PER_GIG_MONTH = regexp.MustCompile(
	`^(\$[0-9.]+ per GB - first [0-9]+ TB / month of storage used|\$[0-9.]+ per GB-Month of storage used in Standard-Infrequent Access)$`)

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

type StorageRecord struct {
	usageStartDate string
	usageEndDate   string
	usageQuantity  string
}

func MakeCollector(options *Options) *BillingAwsCollector {
	validateOptions(options)
	bigqueryConn := bigquery.NewBigqueryConnection(&options.Bigquery)

	s3Conn := s3.NewS3Connection(&options.S3)
	bestPath := listBestPathForThisMonth(s3Conn)
	if bestPath == "" {
		log.Fatalf("Couldn't find billing path for today")
	}

	zipDownloader := s3Conn.DownloadPath(bestPath)
	zipBytes, err := ioutil.ReadAll(zipDownloader)
	if err != nil {
		log.Fatalf("Error from ioutil.ReadAll of %s: %s", bestPath, err)
	}

	zipReader, err := zip.NewReader(bytes.NewReader(zipBytes), int64(len(zipBytes)))
	if err != nil {
		log.Fatalf("Error from zip.NewReader: %s", err)
	}

	s3BucketNameToStorageRecord := map[string]StorageRecord{}
	for _, f := range zipReader.File {
		rc, err := f.Open()
		if err != nil {
			log.Fatalf("Error from f.Open() of %s: %s", f.Name, err)
		}

		csvReader := csv.NewReader(rc)
		lineNum := 1
		headings, err := csvReader.Read()
		if err == io.EOF {
			log.Fatalf("Empty file %s inside zip", f.Name)
		}
		if err != nil {
			log.Fatalf("Error from csvReader.Read() line %d: %s", lineNum, err)
		}

		headingToColNum := map[string]int{}
		for colNum, heading := range headings {
			headingToColNum[heading] = colNum
		}
		productNameCol, ok := headingToColNum["ProductName"]
		if !ok {
			log.Fatalf("Can't find ProductName heading in %s", headings)
		}
		operationCol, ok := headingToColNum["Operation"]
		if !ok {
			log.Fatalf("Can't find Operation heading in %s", headings)
		}
		itemDescriptionCol, ok := headingToColNum["ItemDescription"]
		if !ok {
			log.Fatalf("Can't find ItemDescription heading in %s", headings)
		}
		usageStartDateCol, ok := headingToColNum["UsageStartDate"]
		if !ok {
			log.Fatalf("Can't find UsageStartDate heading in %s", headings)
		}
		usageEndDateCol, ok := headingToColNum["UsageEndDate"]
		if !ok {
			log.Fatalf("Can't find UsageEndDate heading in %s", headings)
		}
		usageQuantityCol, ok := headingToColNum["UsageQuantity"]
		if !ok {
			log.Fatalf("Can't find UsageQuantity heading in %s", headings)
		}
		resourceIdCol, ok := headingToColNum["ResourceId"]
		if !ok {
			log.Fatalf("Can't find ResourceId heading in %s", headings)
		}

		for {
			lineNum += 1
			values, err := csvReader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error from csvReader.Read() line %d: %s", lineNum, err)
			}

			if values[productNameCol] == "Amazon Simple Storage Service" &&
				(values[operationCol] == "StandardStorage" ||
					values[operationCol] == "StandardIAStorage") { // IA = Infrequent Access
				if !PER_GIG_MONTH.MatchString(values[itemDescriptionCol]) {
					log.Fatalf("Unexpected ItemDescription %s", values[itemDescriptionCol])
				}
				s3BucketName := values[resourceIdCol]

				oldRecord, exists := s3BucketNameToStorageRecord[s3BucketName]
				if !exists ||
					oldRecord.usageStartDate < values[usageStartDateCol] {
					s3BucketNameToStorageRecord[s3BucketName] = StorageRecord{
						usageStartDate: values[usageStartDateCol],
						usageEndDate:   values[usageEndDateCol],
						usageQuantity:  values[usageQuantityCol],
					}
				}
			}
		}
		rc.Close()
	}

	combinedBucketNameToUsageBytes := map[string]float64{}
	for s3BucketName, record := range s3BucketNameToStorageRecord {
		startTime, err := time.Parse("2006-01-02 15:04:05", record.usageStartDate)
		if err != nil {
			log.Fatalf("Couldn't parse UsageStartDate %s", record.usageStartDate)
		}

		endTime, err := time.Parse("2006-01-02 15:04:05", record.usageEndDate)
		if err != nil {
			log.Fatalf("Couldn't parse UsageEndDate %s", record.usageEndDate)
		}

		usageQuantity, err := strconv.ParseFloat(record.usageQuantity, 64)
		if err != nil {
			log.Fatalf("Couldn't parse UsageQuantity %s", record.usageQuantity)
		}

		dayDuration := endTime.Sub(startTime)
		if dayDuration < 58*time.Minute || dayDuration > 62*time.Minute {
			log.Fatalf("Unexpected non-hour duration %s to %s",
				record.usageStartDate, record.usageEndDate)
		}

		firstOfMonth := startTime.AddDate(0, 0, -startTime.Day()+1)
		lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
		numDaysInMonth := lastOfMonth.Day()
		numBytes := usageQuantity * GIBIBYTES_TO_BYTES * float64(numDaysInMonth)

		combinedBucketName := s3BucketName
		if numBytes < MIN_BUCKET_BYTES_TO_HAVE_OWN_METRIC {
			combinedBucketName = OTHER_BUCKET_NAME
		}
		combinedBucketNameToUsageBytes[combinedBucketName] += numBytes
	}
	log.Printf("%v", combinedBucketNameToUsageBytes)

	return NewBillingAwsCollector(options, bigqueryConn, s3Conn)
}
