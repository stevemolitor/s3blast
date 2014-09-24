package main

import (
        "fmt"
        "os"
        "sort"
        "launchpad.net/goamz/s3"
        "launchpad.net/goamz/aws"
)

const (
        LIST_ITEM_FORMAT = "%16s %10d\t%s\n"
        LIST_DIR_FORMAT = "%35s\t%s\n"
)

// Listing of S3 bucket, including "directories" (common prefixes)
// and items (keys).
type Listing struct{
        directories []string
        items []s3.Key
}

// Return sorted, uniq list of common prefixes ("directories") of Listing.
func (l *Listing) Directories() ([]string) {
        dirs := uniq(l.directories)
        sort.Strings(dirs)
        return dirs
}

// Return items (keys) of Listing.
func (l *Listing) Items() ([]s3.Key) {
        return l.items
}

// Concatenate one Listing to another, returning a new Listing
// containing the union of both.
func (l *Listing) Concat(o *Listing) (*Listing) {
        return &Listing{
                append(l.directories, o.directories...),
                append(l.items, o.items...)}
}

// Get the next marker strings, for use when paginating S3 list
// results.  S3 limits the result of a list call 1000 entries. The
// value to be passed as the next "marker" is the greater of either
// the last common prefix returned, or the last key returned.
func getNextMarker(result *s3.ListResp) string {
        prefixesLen := len(result.CommonPrefixes)
        contentsLen := len(result.Contents)

        if prefixesLen == 0 {
                return result.Contents[contentsLen-1].Key
        } else if contentsLen == 0 {
                return result.CommonPrefixes[prefixesLen-1]
        } else if result.Contents[contentsLen-1].Key > result.CommonPrefixes[prefixesLen-1] {
                return result.Contents[contentsLen-1].Key
        } else {
                return result.CommonPrefixes[prefixesLen-1]
        }
}

// Get Listing of bucket, starting with prefix, starting with marker.
func GetListing(bucketName string, prefix string, marker string) (listing *Listing, err error) {
        auth, err := aws.EnvAuth()
        if err != nil {
                fmt.Println("Error getting aws credentials", err)
                os.Exit(1)
        }

        s3 := s3.New(auth, aws.USEast)
        bucket := s3.Bucket(bucketName)
        result, err := bucket.List(prefix, "/", marker, 0)
        if err != nil {
                return nil, err
        }

        listing = &Listing{result.CommonPrefixes, result.Contents}

        if (result.IsTruncated) {
                nextMarker := getNextMarker(result)
                rest, err := GetListing(bucketName, prefix, nextMarker)
                if err != nil {
                        return nil, err
                }
                listing = listing.Concat(rest)
        } 
        
        return listing, nil
}

// List items and common prefixes, in bucket, starting with prefix,
// and print results to standard output.  If prefix is "", print all
// items in bucket.
func List(bucketName string, prefix string) {
        listing, err := GetListing(bucketName, prefix, "")
        if (err != nil) {
                fmt.Println("Error listing bucket -", err)
                os.Exit(1)
        }
        
        for _, pre := range listing.Directories() {
                fmt.Printf(LIST_DIR_FORMAT, "DIR", pre)
        }

        for _, item := range listing.Items() {
                fmt.Printf(LIST_ITEM_FORMAT, item.LastModified, item.Size, item.Key)
        }
}
