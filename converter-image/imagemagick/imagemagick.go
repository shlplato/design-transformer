// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// [START run_imageproc_handler_setup]

// Package imagemagick contains an example of using ImageMagick to process a
// file uploaded to Cloud Storage.
package imagemagick

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"cloud.google.com/go/storage"
)

// Global API clients used across function invocations.
var (
	storageClient *storage.Client
)

func init() {
	// Declare a separate err variable to avoid shadowing the client variables.
	var err error

	storageClient, err = storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("storage.NewClient: %v", err)
	}

}

// [END run_imageproc_handler_setup]

// [START run_imageproc_handler_analyze]

// GCSEvent is the payload of a GCS event.
type GCSEvent struct {
	Bucket string `json:"bucket"`
	Name   string `json:"name"`
}

// ConvertImageToPNG convert PSD images uploaded to GCS to PNG.
func ConvertImageToPNG(ctx context.Context, e GCSEvent) error {
	outputBucket := os.Getenv("CONVERTED_BUCKET_NAME")
	if outputBucket == "" {
		return errors.New("CONVERTED_BUCKET_NAME must be set")
	}
	return convert(ctx, e.Bucket, outputBucket, e.Name)
}

// [END run_imageproc_handler_analyze]

// [START run_imageproc_handler_blur]

// blur blurs the image stored at gs://inputBucket/name and stores the result in
// gs://outputBucket/name.
func convert(ctx context.Context, inputBucket, outputBucket, name string) error {
	inputBlob := storageClient.Bucket(inputBucket).Object(name)
	r, err := inputBlob.NewReader(ctx)
	if err != nil {
		return fmt.Errorf("NewReader: %v", err)
	}

	convertType := os.Getenv("TYPE")
	if convertType == "" {
		convertType = "spreadshirt"
	}

	paths := strings.Split(name, "/");
	paths = append(paths, "");
	copy(paths[3:], paths[2:]);
	paths[2] = convertType;
	name = strings.Join(paths, "/"); 
	outputBlob := storageClient.Bucket(outputBucket).Object(convertType + "/" + strings.Replace(name, ".psd",".png", -1))
	w := outputBlob.NewWriter(ctx)
	defer w.Close()

	// Use - as input and output to use stdin and stdout.
//	cmd := exec.Command("convert", "wonder_to_flight.psd[0]", "png:-")
	cmd := exec.Command("convert", "-[0]", "png:-")
	if convertType == "amazon" {
		cmd = exec.Command("convert", "-[0]", "-background", "none", "-resize", "4500x5400", "-extent", "4500x5400", "png:-")
	} 
	cmd.Stdin = r
	cmd.Stdout = w
 
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("cmd.Run: %v", err)
	}

	log.Printf("Converted image uploaded to gs://%s/%s", outputBlob.BucketName(), outputBlob.ObjectName())

//	if err := storageClient.Bucket(inputBucket).Object(name).Delete(ctx); err != nil {
//		log.Printf("Source image could not deleted: %v", err)
//	} else {
//		log.Printf("Source image deleted: %s", name)
//	}
	return nil
}

// [END run_imageproc_handler_blur]
