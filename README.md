# Samuno - mounting S3 drives in Windows with user access
## About
**Samuno** - is a simple software that allows you to mount S3 buckets (or directories) in Windows as network shares via access key that is bound to a specific person.

The aim is to manage and control access to network drives for employees in a company. 

The software relies on Rclone as backend (+ WinFsp library) and N8N + Airtable services for authenications.

<kbd>![# Demo](https://raw.githubusercontent.com/keshon/assets/main/samuno/demo.gif)</kbd>

## How is it working
0. Create two pair of access & secret keys for S3 provider: one for full access, another for read only.
1. Add new user inside Airtable and create a unique key for him.
2. Create new path inside Airtable using `butcket-name/subpath` pattern.
3. Assing paths to the user.
4. Provide user with the key and link to compiled binary of Samuno.
5. User downloads the software (and WinFcp library), run it and under the Settings fill-in the given key.
6. User press Connect button - the key is being validated and assigned buckets are being mounted.

## Folder content
- **3rd-party-tools** - binary tools are needed to generate icon and syso files.
- **rcd-server** - sources for server part of the application.
- **user-access** - example for [Airtable](http://airtable.com "Airtable") tables structure and json workflow for [N8N](https://n8n.io/ "N8N").
- **win-client** - sources for front-end part of the application. GUI is based on modified version of [NanoGUI](https://github.com/shibukawa/nanogui-go "NanoGUI") framework.