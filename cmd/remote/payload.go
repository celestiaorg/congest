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

const (
	chainIDFlag    = "chainID"
	payloadFlag    = "payload"
	validatorsFlag = "validators"
)

// PayloadCmd is the Cobra command for creating the payload for the experiment.
func PayloadCmd() *cobra.Command {
	payloadCmd := &cobra.Command{
		Use:   "payload",
		Short: "Create the payload for the experiment",
		RunE: func(cmd *cobra.Command, args []string) error {
			chainID, ppath, vpath := readFlags(cmd)
			ips, err := ReadValidatorsFromFile(vpath)
			if err != nil {
				return err
			}
			err = createPayload(ips, chainID, ppath)
			if err != nil {
				log.Fatalf("Failed to create payload: %v", err)
			}

			return nil
		},
	}

	// Flags for the payload command
	payloadCmd.Flags().StringP(chainIDFlag, "c", "test", "Chain ID (required)")
	payloadCmd.MarkFlagRequired(chainIDFlag)
	payloadCmd.Flags().StringP(payloadFlag, "p", "./payload", "Path to the payload directory (required)")
	payloadCmd.MarkFlagRequired(payloadFlag)
	payloadCmd.Flags().StringP(validatorsFlag, "v", "./payload/validators.json", "Path to the validators file (required)")
	payloadCmd.MarkFlagRequired(validatorsFlag)
	return payloadCmd
}

func readFlags(cmd *cobra.Command) (chainID, ppath, vpath string) {
	chainID, err := cmd.Flags().GetString(chainIDFlag)
	if err != nil {
		log.Fatalf("Failed to read chainID flag: %v", err)
	}
	ppath, err = cmd.Flags().GetString(payloadFlag)
	if err != nil {
		log.Fatalf("Failed to read payload flag: %v", err)
	}
	vpath, err = cmd.Flags().GetString(validatorsFlag)
	if err != nil {
		log.Fatalf("Failed to read validators flag: %v", err)
	}
	return chainID, ppath, vpath
}

// createPayload takes ips created by pulumi the path to the payload directory
// to create the payload required for the experiment.
func createPayload(ips map[string]NodeInfo, chainID, ppath string) error {
	n, err := NewNetwork(chainID)
	if err != nil {
		return err
	}

	for _, info := range ips {
		n.SyncAddValidator(
			info.Name,
			info.Region,
			info.IP,
			ppath,
		)
	}

	err = n.InitNodes(ppath)
	if err != nil {
		return err
	}

	err = n.SaveAddressBook(ppath, n.Peers())
	if err != nil {
		return err
	}

	return nil

}

func main() {
	rootCmd := &cobra.Command{Use: "s3cli"}
	rootCmd.AddCommand(NewDownloadCmd())
	rootCmd.AddCommand(PayloadCmd())
	rootCmd.Execute()
}
