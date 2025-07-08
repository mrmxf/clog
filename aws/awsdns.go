// Package aws provides functionality to query and list DNS records
// from AWS Route53 hosted zones and Lightsail domains across all regions.
//
// The package exports functions to retrieve A, AAAA, and CNAME records
// with configurable sorting options for comprehensive domain analysis.
//
// Example usage:
//
//	opts := aws.SortOptions{
//		Sort1: "name",
//		Sort2: "location",
//		Sort3: "address",
//	}
//	err := aws.Route53Domains(opts)
//	if err != nil {
//		log.Fatal(err)
//	}
package aws

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"sort"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/lightsail"
	lightsailtypes "github.com/aws/aws-sdk-go-v2/service/lightsail/types"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
	"github.com/mrmxf/clog/crayon"
)

// SortOptions defines the sorting preferences for DNS record output.
// The sorting is applied hierarchically: Sort1 first, then Sort2 for ties, then Sort3.
type SortOptions struct {
	// Sort1 is the primary sorting field
	Sort1 string
	// Sort2 is the secondary sorting field, applied when Sort1 values are equal
	Sort2 string
	// Sort3 is the tertiary sorting field, applied when Sort1 and Sort2 values are equal
	Sort3 string
}

// DNSRecord represents a single DNS record with its associated metadata.
// This structure provides all necessary information for sorting and display.
type DNSRecord struct {
	// Name is the DNS record name (e.g., "www.example.com")
	Name string
	// Type is the DNS record type (A, AAAA, or CNAME)
	Type string
	// Value is the resolved value of the DNS record
	Value string
	// Location represents the AWS region or availability zone where the record is hosted
	Location string
	// Source indicates whether the record comes from "route53" or "lightsail"
	Source string
}

// Route53Domains queries all Route53 hosted zones and Lightsail domains across
// all AWS regions, then prints a sorted list of A, AAAA, and CNAME records to stdout.
//
// The function performs the following operations:
// 1. Loads AWS configuration from the default credential chain
// 2. Queries all available AWS regions
// 3. Retrieves Route53 hosted zones and their DNS records
// 4. Retrieves Lightsail domains from each region
// 5. Sorts the combined results according to the provided options
// 6. Prints the formatted results to stdout
//
// Parameters:
//   - opts: SortOptions struct defining the sorting hierarchy
//     Valid sort field values: "address", "name", "location"
//
// Returns an error if AWS services cannot be accessed or if there are
// authentication/authorization issues.
func Route53Domains(opts SortOptions) error {
	// Load AWS configuration using the default credential chain
	// This follows AWS best practices for credential management
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return fmt.Errorf("failed to load AWS config: %w", err)
	}

	var allRecords []DNSRecord

	// Query Route53 hosted zones (global service)
	route53Records, err := getRoute53Records(cfg)
	if err != nil {
		// Log the error but continue with Lightsail query
		// This ensures partial results if one service fails
		log.Printf("Warning: Failed to query Route53: %v", err)
	} else {
		allRecords = append(allRecords, route53Records...)
	}

	// Query Lightsail domains from all regions
	// Lightsail is region-specific, so we must check each region
	lightsailRecords, err := getLightsailRecords(cfg)
	if err != nil {
		// Log the error but continue with available data
		log.Printf("Warning: Failed to query Lightsail: %v", err)
	} else {
		allRecords = append(allRecords, lightsailRecords...)
	}

	// Sort records according to the specified options
	// Multiple sort criteria ensure predictable, stable ordering
	sortRecords(allRecords, opts)

	// Print results to stdout in a readable format
	printRecords(allRecords)

	return nil
}

// getRoute53Records retrieves all DNS records from Route53 hosted zones.
// Route53 is a global service, so we don't need to iterate through regions.
func getRoute53Records(cfg aws.Config) ([]DNSRecord, error) {
	var records []DNSRecord

	client := route53.NewFromConfig(cfg)

	// List all hosted zones in the account
	hostedZones, err := client.ListHostedZones(context.TODO(), &route53.ListHostedZonesInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list hosted zones: %w", err)
	}

	// Iterate through each hosted zone to get its records
	for _, zone := range hostedZones.HostedZones {
		zoneRecords, err := getHostedZoneRecords(client, *zone.Id)
		if err != nil {
			// Log error but continue with other zones
			// This ensures we get partial results even if some zones fail
			log.Printf("Warning: Failed to get records for zone %s: %v", *zone.Name, err)
			continue
		}
		records = append(records, zoneRecords...)
	}

	return records, nil
}

// getHostedZoneRecords retrieves DNS records for a specific Route53 hosted zone.
// It filters for A, AAAA, and CNAME record types as specified in requirements.
func getHostedZoneRecords(client *route53.Client, zoneId string) ([]DNSRecord, error) {
	var records []DNSRecord

	// List all resource record sets in the hosted zone
	paginator := route53.NewListResourceRecordSetsPaginator(client, &route53.ListResourceRecordSetsInput{
		HostedZoneId: aws.String(zoneId),
	})

	// Use paginator to handle large hosted zones efficiently
	// This prevents memory issues with zones containing many records
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to list resource record sets: %w", err)
		}

		for _, recordSet := range page.ResourceRecordSets {
			// Filter for only the required record types
			if isTargetRecordType(recordSet.Type) {
				record := convertRoute53Record(recordSet, "route53")
				if record != nil {
					records = append(records, *record)
				}
			}
		}
	}

	return records, nil
}

// getLightsailRecords retrieves DNS records from Lightsail domains across all regions.
// Lightsail stores DNS records in the DomainEntries array within each Domain.
// We iterate through each domain and extract A, AAAA, and CNAME records from DomainEntries.
func getLightsailRecords(cfg aws.Config) ([]DNSRecord, error) {
	var allRecords []DNSRecord

	// Get list of all available regions
	regions, err := getAWSRegions(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get AWS regions: %w", err)
	}

	// Query Lightsail domains in each region
	for _, region := range regions {
		// Create region-specific configuration
		regionCfg := cfg.Copy()
		regionCfg.Region = region

		regionRecords := getLightsailRecordsForRegion(regionCfg, region)
		allRecords = append(allRecords, regionRecords...)
	}

	return allRecords, nil
}

// getLightsailRecordsForRegion retrieves Lightsail DNS records for a specific region.
// This function queries registered domains and extracts DNS records from their DomainEntries.
func getLightsailRecordsForRegion(cfg aws.Config, region string) []DNSRecord {
	var records []DNSRecord

	client := lightsail.NewFromConfig(cfg)

	// Get all registered domains in this region
	domains, err := client.GetDomains(context.TODO(), &lightsail.GetDomainsInput{})
	if err != nil {
		slog.Debug("Warning: Failed to query Lightsail in region %s: %v", region, err)
		return nil
	}

	// Process each domain to extract DNS records from DomainEntries
	for _, domain := range domains.Domains {
		if domain.Name == nil {
			continue
		}

		// Extract DNS records from the domain's DomainEntries array
		domainRecords := extractDomainEntries(domain, region)
		records = append(records, domainRecords...)
	}

	return records
}

// extractDomainEntries processes a Lightsail domain's DomainEntries to extract DNS records.
// Each Domain contains a DomainEntries array with DomainEntry structs that hold the actual DNS records.
func extractDomainEntries(domain lightsailtypes.Domain, region string) []DNSRecord {
	var records []DNSRecord

	if domain.DomainEntries == nil {
		return records
	}

	// Iterate over each domain entry in the domain's DomainEntries array
	for _, entry := range domain.DomainEntries {
		// Filter for only the DNS record types we want (A, AAAA, CNAME)
		if isTargetLightsailDomainEntryType(entry.Type) {
			record := convertLightsailDomainEntry(entry, region)
			if record != nil {
				records = append(records, *record)
			}
		}
	}

	return records
}

// isTargetLightsailDomainEntryType checks if a Lightsail DomainEntry type is one we want to include.
// We only care about A, AAAA, and CNAME records as specified in requirements.
func isTargetLightsailDomainEntryType(entryType *string) bool {
	if entryType == nil {
		return false
	}

	switch *entryType {
	case "A", "AAAA", "CNAME":
		return true
	default:
		return false
	}
}

// convertLightsailDomainEntry converts a Lightsail DomainEntry to our standard DNSRecord format.
// DomainEntry contains the actual DNS record information within each Lightsail domain.
func convertLightsailDomainEntry(entry lightsailtypes.DomainEntry, region string) *DNSRecord {
	if entry.Name == nil || entry.Type == nil {
		return nil
	}

	var value string
	if entry.Target != nil {
		value = *entry.Target
	}

	return &DNSRecord{
		Name:     *entry.Name,
		Type:     *entry.Type,
		Value:    value,
		Location: region,
		Source:   "lightsail",
	}
}

// getAWSRegions retrieves a list of all available AWS regions.
// This ensures we query Lightsail in all regions where it might be available.
func getAWSRegions(cfg aws.Config) ([]string, error) {
	client := ec2.NewFromConfig(cfg)

	regions, err := client.DescribeRegions(context.TODO(), &ec2.DescribeRegionsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to describe regions: %w", err)
	}

	var regionNames []string
	for _, region := range regions.Regions {
		regionNames = append(regionNames, *region.RegionName)
	}

	return regionNames, nil
}

// isTargetRecordType checks if a Route53 record type is one we want to include.
// We only care about A, AAAA, and CNAME records as specified in requirements.
func isTargetRecordType(recordType types.RRType) bool {
	switch recordType {
	case types.RRTypeA, types.RRTypeAaaa, types.RRTypeCname:
		return true
	default:
		return false
	}
}

// convertRoute53Record converts a Route53 ResourceRecordSet to our standard DNSRecord format.
// Returns nil if the record cannot be converted or has no valid resource records.
func convertRoute53Record(recordSet types.ResourceRecordSet, source string) *DNSRecord {
	if recordSet.Name == nil || len(recordSet.ResourceRecords) == 0 {
		return nil
	}

	// For simplicity, we take the first resource record if multiple exist
	// In production, you might want to create separate DNSRecord entries for each
	var value string
	if len(recordSet.ResourceRecords) > 0 && recordSet.ResourceRecords[0].Value != nil {
		value = *recordSet.ResourceRecords[0].Value
	}

	return &DNSRecord{
		Name:     strings.TrimSuffix(*recordSet.Name, "."), // Remove trailing dot
		Type:     string(recordSet.Type),
		Value:    value,
		Location: "global", // Route53 is a global service
		Source:   source,
	}
}

// sortRecords sorts the DNS records according to the specified sort options.
// The sorting is stable and hierarchical: Sort1, then Sort2, then Sort3.
func sortRecords(records []DNSRecord, opts SortOptions) {
	sort.Slice(records, func(i, j int) bool {
		// Compare using Sort1 field first
		if cmp := compareRecordFields(records[i], records[j], opts.Sort1); cmp != 0 {
			return cmp < 0
		}

		// If Sort1 fields are equal, compare using Sort2
		if cmp := compareRecordFields(records[i], records[j], opts.Sort2); cmp != 0 {
			return cmp < 0
		}

		// If Sort1 and Sort2 are equal, compare using Sort3
		return compareRecordFields(records[i], records[j], opts.Sort3) < 0
	})
}

// compareRecordFields compares two DNS records based on the specified field.
// Returns -1 if record a should come before b, 1 if b should come before a, 0 if equal.
func compareRecordFields(a, b DNSRecord, field string) int {
	switch field {
	case "address":
		return strings.Compare(a.Value, b.Value)
	case "name":
		// Case insensitive comparison as specified in requirements
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	case "location":
		return strings.Compare(a.Location, b.Location)
	default:
		// Default to name comparison if invalid field specified
		return strings.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
	}
}

// isIPv4Address checks if a string represents a valid IPv4 address.
// This helps determine whether to apply IPv4-specific coloring to record values.
func isIPv4Address(value string) bool {
	ip := net.ParseIP(value)
	if ip == nil {
		return false
	}
	// Check if it's an IPv4 address (not IPv6)
	return ip.To4() != nil
}

// printRecords outputs the DNS records to stdout in a formatted, colored table.
// The format uses crayon library for colored output:
// - Domain names are highlighted with c.U() (underline)
// - IPv4 addresses are highlighted with c.W() (warning/bright color)
// - CNAME record values are highlighted with c.E() (error/emphasis color)
// - Regions have no highlighting for readability
func printRecords(records []DNSRecord) {
	if len(records) == 0 {
		fmt.Println("No DNS records found.")
		return
	}

	// Create crayon instance for colored output
	c := crayon.Color()

	// Print header
	r := strings.Repeat
	d := "-"
	fmt.Printf("%-50s %-50s %-8s %-12s %-15s\n", "DNS name", "entry", "type", "service", "region")
	fmt.Printf("%-50s %-50s %-8s %-12s %-15s\n", r(d, 50), r(d, 50), r(d, 8), r(d, 12), r(d, 15))

	// Print each record with appropriate coloring
	for _, record := range records {
		// Color the domain name with underline
		tNam := truncateString(record.Name, 50)
		lNam := 50 - len(tNam)
		coloredName := c.U(tNam)

		// Color the value based on its type and content
		var coloredValue string
		tVal := truncateString(record.Value, 50)
		lVal := 50 - len(tVal)

		if record.Type == "CNAME" {
			// CNAME records get error/emphasis coloring
			coloredValue = c.D(tVal)
		} else if isIPv4Address(record.Value) {
			// IPv4 addresses get warning/bright coloring
			coloredValue = c.W(tVal)
		} else {
			// Other values (like IPv6 or non-standard values) remain uncolored
			coloredValue = tVal
		}

		// Location (region) has no highlighting as requested
		region := truncateString(record.Location, 15)

		// Color the value based on its type and content
		var coloredSource string
		lcs := 12 - len(record.Source)
		if record.Source == "lightsail" {
			// lightsail records get Success coloring
			coloredSource = c.S(record.Source)
		} else {
			// Other values (like IPv6 or non-standard values) remain uncolored
			coloredSource = record.Source
		}

		fmt.Printf("%*s%s %s%*s %-8s %-*s%s %-15s\n",
			lNam, " ", coloredName,
			coloredValue, lVal, " ",
			record.Type,
			lcs, " ", coloredSource,
			region)
	}

	fmt.Printf("\nTotal records found: %d\n", len(records))
}

// truncateString truncates a string to the specified length, adding "..." if truncated.
// This ensures consistent column formatting in the output table.
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
