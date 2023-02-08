# Omega Drive client
## What is Omega Drive
Omega Drive is my pet project that I started while working as a CTO in a small animation studio.
It is used to provide colleagues with distributed access to various paths of the S3 object storage and map them as local drives. 
It supports read-only access or full access depending on the server settings.
Omega Drive client acts as GUI wrapper on top of [Rclone](https://github.com/rclone/rclone) - cli application to managing and mounting various type of network storages.

## What components does OD consist of
Omega Drive has three main components:
1. A client using [Fyne](https://github.com/fyne-io/fyne) as GUI
2. [RCD binary](https://github.com/keshon/omega-drive-client-rcd) - a compiled version of [Rclone](https://github.com/rclone/rclone) with a small wrapper on top.
3. [Server environment](https://github.com/keshon/omega-drive-server) - mix of N8n.io workflow and Airtable.

## How to assemble the client
To build the client part, you need to perform several steps:
1. Rename `/conf/conf.go.BLANK` to `/conf/conf.go`
2. Populate the `conf.go` file with the appropriate information (mostly N8n endpoint paths).
3. Compile the RCD client binary and place it in the root of the client folder.