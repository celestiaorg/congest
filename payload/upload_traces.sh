#!/bin/bash
export DEBIAN_FRONTEND=noninteractive

export GOOGLE_APPLICATION_CREDENTIALS="/root/payload/congest-remote-key-gbq.json"
sudo apt-get install apt-transport-https ca-certificates gnupg curl -y
curl https://packages.cloud.google.com/apt/doc/apt-key.gpg | sudo gpg --dearmor -o /usr/share/keyrings/cloud.google.gpg
echo "deb [signed-by=/usr/share/keyrings/cloud.google.gpg] https://packages.cloud.google.com/apt cloud-sdk main" | sudo tee -a /etc/apt/sources.list.d/google-cloud-sdk.list
# sudo apt-get update -y && sudo apt-get install google-cloud-cli -y

# ensure that the env vars are exported here
source /root/payload/vars.sh
echo "CHAIN_ID after sourcing vars.sh: $CHAIN_ID"

# Set environment variables
PROJECT_ID="numeric-mile-433416-e9"
DATASET_ID="traces"
CHAIN_ID="big-blonks-22"
LOCAL_DIR="/root/.celestia-app/data/traces"

# gcloud auth activate-service-account --key-file="/root/payload/congest-remote-key-gbq.json" -q
# gcloud config set project numeric-mile-433416-e9 -q

tmux kill-session -t app

# Get the hostname
hostname=$(hostname)

# Parse the first part of the hostname
nodeID=$(echo $hostname | awk -F'-' '{print $1 "-" $2}')

source_dir="/root/.celestia-app/data/traces"
logs_path="/root/logs.txt"

# clean the data by removing the last line
find $source_dir -type f -name "*.jsonl" -exec sed -i '$d' {} \;

#!/bin/bash



# # Iterate over all JSONL files in the directory
# for file in ${LOCAL_DIR}/*.jsonl; do
#   # Extract the filename without the extension to use as the table name
#   TABLE_ID=$(basename "${file}" .jsonl)

#   echo "Loading ${file} into BigQuery table ${TABLE_ID}..."

#   # Load each JSONL file directly into BigQuery
#   bq load --source_format=NEWLINE_DELIMITED_JSON --autodetect \
#     ${PROJECT_ID}:${DATASET_ID}.${TABLE_ID} \
#     ${file}
  
#   if [ $? -eq 0 ]; then
#     echo "Successfully loaded ${file} into table ${TABLE_ID}"
#   else
#     echo "Failed to load ${file} into table ${TABLE_ID}"
#   fi
# done

echo "All files loaded."

sudo snap install aws-cli --classic
# destination_file="/tmp/${CHAIN_ID}_${nodeID}_traces.tar.gz"


Compress the directory into a tar.gz file
# tar -czvf "$destination_file" -C "$source_dir" .

# Set the base S3 path
base_s3_path="s3://${S3_BUCKET_NAME}/${CHAIN_ID}/${nodeID}/"

# Upload the directory structure to S3
aws s3 cp "$source_dir" "$base_s3_path" --recursive --region $AWS_DEFAULT_REGION
aws s3 cp "$logs_path" "$base_s3_path" --region $AWS_DEFAULT_REGION


# # Upload the tar.gz file to S3
# aws s3 cp "$destination_file" "s3://${S3_BUCKET_NAME}/${CHAIN_ID}/${nodeID}/traces.tar.gz" --region $AWS_DEFAULT_REGION

# # Clean up the local tar.gz file
# rm "$destination_file"
