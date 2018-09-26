# etcd backblaze b2 snapshot

Downloads snapshot from b2 to initalize etcd instance

## Environment variables

- `EBS_B2_APPLICATION_ID`: b2 application id
- `EBS_B2_APPLICATION_KEY`: b2 application key
- `EBS_B2_BUCKET_ID`: b2 bucket id
- `EBS_B2_OBJECT`: b2 object to download (snapshot)
- `EBS_B2_OBJECT_ID`: b2 object ID to download, if defined, this takes precedence over `EBS_B2_OBJECT`
- `EBS_B2_DOWNLOAD_RETRY_INTERVAL`: time to retry download if retry exists (ms)
