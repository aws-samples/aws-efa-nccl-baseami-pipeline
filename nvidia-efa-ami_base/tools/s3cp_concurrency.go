//Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
//Permission is hereby granted, free of charge, to any person obtaining a copy of this
//software and associated documentation files (the "Software"), to deal in the Software
//without restriction, including without limitation the rights to use, copy, modify,
//merge, publish, distribute, sublicense, and/or sell copies of the Software, and to
//permit persons to whom the Software is furnished to do so.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
//INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
//PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
//HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
//OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
//SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
//
// Author: Amr Ragab amrraga@amazon.com

package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var DEFAULT_THREADS = 1
var MAX_THREADS = 1024
var DEFAULT_PART_SIZE int64 = 5 * 1024 * 1024
var MAX_PART_SIZE int64 = 128 * 1024 * 1024
var objLength int64

func main() {
	if len(os.Args) < 7 {
		fmt.Fprintf(os.Stderr, "usage: %s <upload|download> <file> <bucket> <key> <nthreads> <part_size>\n", os.Args[0])
		os.Exit(1)
	}
	direction := os.Args[1]
	src_path := os.Args[2]
	s3_bucket := os.Args[3]
	s3_key := os.Args[4]
	nthreads, err := strconv.Atoi(os.Args[5])
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: nthreads must be an integer - %v\n", err)
		os.Exit(1)
	}
	part_size, err := strconv.ParseInt(os.Args[6], 10, 64)
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: part_size must be an integer - %v\n", err)
		os.Exit(1)
	}
	part_size *= 1024 * 1024

	if nthreads > MAX_THREADS {
		fmt.Fprintf(os.Stderr, "INFO: Number of threads too high (%d), changing it to max (%d)\n", nthreads, MAX_THREADS)
		nthreads = MAX_THREADS
	}
	if nthreads < 1 {
		fmt.Fprintf(os.Stderr, "INFO: Number of threads too low (%d), changing it to default (%d)\n", nthreads, DEFAULT_THREADS)
		nthreads = DEFAULT_THREADS
	}
	if part_size < DEFAULT_PART_SIZE {
		fmt.Fprintf(os.Stderr, "INFO: Part size too small (%d), changing it to default (%d)\n", part_size, DEFAULT_PART_SIZE)
		part_size = DEFAULT_PART_SIZE
	}
	if strings.HasPrefix(direction, "upload") {
		src_info, err := os.Stat(src_path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
		src_size := src_info.Size()
		body, err := os.Open(src_path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
			os.Exit(1)
		}
		sess := session.Must(session.NewSession())
		uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
			u.PartSize = part_size
			u.Concurrency = nthreads
		})
		upParams := &s3manager.UploadInput{
			Bucket: &s3_bucket,
			Key:    &s3_key,
			Body:   body,
		}
		start := time.Now()
		_, err = uploader.Upload(upParams)
		if err != nil {
			fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
			os.Exit(1)
		}
		elapsed := time.Since(start).Seconds()
		bandwidth := float64(src_size) / elapsed / (1024.0 * 1024.0)
		fmt.Fprintf(os.Stdout, "INFO: copied %d bytes in %f seconds - %f MB/s\n", src_size, elapsed, bandwidth)
		os.Exit(0)
	}
	if strings.HasPrefix(direction, "download") {
		file, err := os.Create(src_path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to create file %q, %v\n", src_path, err)
			os.Exit(1)
		}
		defer file.Close()
		sess := session.Must(session.NewSession())
		downloader := s3manager.NewDownloader(sess, func(d *s3manager.Downloader) {
			d.PartSize = part_size
			d.Concurrency = nthreads
		})
		start := time.Now()
		numBytes, err := downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: &s3_bucket,
				Key:    &s3_key,
			})
		if err != nil {
			fmt.Fprintf(os.Stderr, "FATAL: %v\n", err)
			os.Exit(1)
		}
		elapsed := time.Since(start).Seconds()
		bandwidth := float64(numBytes) / elapsed / (1024.0 * 1024.0)
		fmt.Fprintf(os.Stdout, "INFO: copied %d bytes in %f seconds - %f MB/s\n", numBytes, elapsed, bandwidth)
		os.Exit(0)
	}
}
