package conf

import (
	"encoding/base64"
)

// Config
const (
	// Info
	AppName       = "Samuno"  // specify program name
	AppVersion    = "220228"  // specify version (year month day)
	RcloneVersion = "v1.57.0" // specify rclone version

	// General settings
	ValidTo         = "31 Dec 22"                        // specify the program expiration date like "31 Dec 22". No expiry is set if empty
	IsDev           = false                              // enable extra debugging info if set to true
	EncryptKey      = "qa3rvd6333bjfn39gfebj9u7wvvkzohr" // specify random 32 characters string - needed to encrypt key file
	WebhookInterval = "1m"                               // specify the webhook poling interval. Examples are (without quotes): '2h','5m','14s'

	// Rclone  settings
	/*
		Samuno heavily relies on (wraped) rclone cli utility that works as rc server (https://rclone.org/rc/)
		To enhance security fronend communicates with the server using following credentials and settings.
	*/
	RcCreateConfig     = true                     // enable rclone config should be automatically created upon program start
	RcAutodeleteConfig = true                     // enable if rclone config should be always deleted on program start/shutdown
	RcUsername         = "someUsername"           // specify username for rclone rcd server authentication
	RcPassword         = "somePassword"           // specify password for rclone rcd server authentication
	RcHost             = "http://localhost:5579/" // specify hostname for rclone rcd

	// N8N settings
	/*
		N8N acts as API bridge between Samuno and Airtable. To make it work it is necessary to
		import N8N workflow (see /user-access folder) to actual N8N instanse and setup following
		parameters on both ends (here and in the imported workflow).
	*/
	N8nUsername = "someUsername"                                           // specify username for n8n webhook authentication
	N8nPassword = "somePassword"                                           // specify password for n8n webhook authentication
	N8nURL      = "https://n8n.example.com/webhook/some-prod-endpoint"     // specify webhook production url
	N8nDevURL   = "https://n8n.example.com/webhook-test/some-dev-endpoint" // specify webhook development url

	// S3 Object Storage settings
	/*
		Settings to connect to S3 that are being used by rclone to moun/unmount buckets.
		In order to get Read&Write and Read-only capabilities S3 administrator need to create
		two different pairs of secret and access keys (using ACL or similar): one pair for full access (ReadWrite) and
		another is for read only. That would allows us to toggle rw/r behavior in Airtable.
	*/
	S3Name     = "providerName"                // specify name of the provider
	S3Region   = "someRegion"                  // specify region e.g. 'eu-central-1'
	S3Endpoint = "s3.eu-central-1.example.com" // specify endpoint URL

	S3ReadWriteAccessKey = "ACCESS KEY GOES HERE" // specify access key with read & write permissions
	S3ReadWriteSecretKey = "SECRET KEY GOES HERE" // specify secret key with read & write permissions

	S3ReadOnlyAccessKey = "ACCESS KEY GOES HERE" // specify access key with read only permissions
	S3ReadOnlySecretKey = "SECRET KEY GOES HERE" // specify secret key with read only permissions
)

var (
	// Encoded username:password for rcd server and n8n webhook
	RcAuthEncoded  = base64.StdEncoding.EncodeToString([]byte(RcUsername + ":" + RcPassword))
	N8nAuthEncoded = base64.StdEncoding.EncodeToString([]byte(N8nUsername + ":" + N8nPassword))

	/*
		Controls if we want to scope S3 bucket selection by matching part of their names.
		For example we have 3 buckets: my-example1, my-example-2, my-example-3
		If we specify in AllowedKeywords: []string{"example-1", "example-2"}
		that means that we can mount buckets containing 'example-1' and 'example-2' in their names.
	*/
	AllowedKeywords = []string{} // keywords that should be in bucket's name
)
