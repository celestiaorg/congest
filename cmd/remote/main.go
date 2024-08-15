package main

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/pkg/trace"
)

// NewDownloadCmd creates and returns a new download Cobra command.
func NewDownloadCmd() *cobra.Command {
	var (
		dst       string
		prefix    string
		region    string
		accessKey string
		secretKey string
		bucket    string
	)

	// Set defaults from environment variables if they exist
	if v := os.Getenv("AWS_REGION"); v != "" {
		region = v
	}
	if v := os.Getenv("AWS_ACCESS_KEY_ID"); v != "" {
		accessKey = v
	}
	if v := os.Getenv("AWS_SECRET_ACCESS_KEY"); v != "" {
		secretKey = v
	}
	if v := os.Getenv("S3_BUCKET_NAME"); v != "" {
		bucket = v
	}

	downloadCmd := &cobra.Command{
		Use:   "download",
		Short: "Download files from S3",
		Run: func(cmd *cobra.Command, args []string) {
			cfg := trace.S3Config{
				Region:     region,
				AccessKey:  accessKey,
				SecretKey:  secretKey,
				BucketName: bucket,
			}

			err := trace.S3Download(dst, prefix, cfg)
			if err != nil {
				log.Fatalf("Failed to download files: %v", err)
			}

			fmt.Println("Download completed successfully")
		},
	}

	// Flags for the download command with environment variable defaults
	downloadCmd.Flags().StringVar(&dst, "dst", "", "Destination directory (required)")
	downloadCmd.Flags().StringVar(&prefix, "prefix", "", "S3 prefix to filter files (required)")
	downloadCmd.Flags().StringVar(&region, "region", region, "AWS region (default from AWS_REGION environment variable)")
	downloadCmd.Flags().StringVar(&accessKey, "access-key", accessKey, "AWS Access Key ID (default from AWS_ACCESS_KEY_ID environment variable)")
	downloadCmd.Flags().StringVar(&secretKey, "secret-key", secretKey, "AWS Secret Access Key (default from AWS_SECRET_ACCESS_KEY environment variable)")
	downloadCmd.Flags().StringVar(&bucket, "bucket", bucket, "S3 bucket name (default from S3_BUCKET_NAME environment variable)")

	// Marking dst and prefix as required, others have environment variable defaults
	downloadCmd.MarkFlagRequired("dst")
	downloadCmd.MarkFlagRequired("prefix")

	return downloadCmd
}

func main() {
	rootCmd := &cobra.Command{Use: "s3cli"}
	rootCmd.AddCommand(NewDownloadCmd())
	rootCmd.Execute()
}
