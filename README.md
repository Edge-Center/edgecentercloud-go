EdgeCenter cloud API client
====================================

Command line client to EdgeCenter cloud API.

Building the binary locally
------------------------------------
To build the binary locally, run the following command:
```bash
make build
```
This will create a binary called `ec_client` in the `bin` directory.

Getting started
------------------------------------
### Downloading and using the environment file from Vault
To download the environment file from Vault, first, install Vault and jq by running:
```bash
make install_vault
make install_jq
```
Next, set your Vault token to terminal for `install_vault` command

Then, download the environment file by running:
```bash
make download_env_file
```
This will download the environment file and save it as .env :
* **EC_CLOUD_USERNAME** - username
* **EC_CLOUD_PASSWORD** - user's password
* **EC_CLOUD_PROJECT** - project id
* **EC_CLOUD_REGION** - region id
* **EC_CLOUD_AUTH_URL** - authentication url, you could use the same as in example above
* **EC_CLOUD_API_URL** - api url, you could use the same as in example above
* **EC_CLOUD_CLIENT_TYPE** - client type, you could use the same as in example above

### Running the client:
After setting the env, use `-h` key to retrieve all available commands:
```bash
./bin/ec_client -h

   NAME:
   ec_client - EdgeCloud API client

   Environment variables example:

   EC_CLOUD_AUTH_URL=
   EC_CLOUD_API_URL=
   EC_CLOUD_API_VERSION=v1
   EC_CLOUD_USERNAME=
   EC_CLOUD_PASSWORD=
   EC_CLOUD_REGION=
   EC_CLOUD_PROJECT=

USAGE:
   ./bin/ec_client [global options] command [command options] [arguments...]

VERSION:
   v0.3.00

COMMANDS:
   network        EdgeCloud networks API
   task           EdgeCloud tasks API
   keypair        EdgeCloud keypairs V2 API
   volume         EdgeCloud volumes API
   subnet         EdgeCloud subnets API
   flavor         EdgeCloud flavors API
   loadbalancer   EdgeCloud loadbalancers API
   instance       EdgeCloud instances API
   heat           EdgeCloud Heat API
   securitygroup  EdgeCloud security groups API
   floatingip     EdgeCloud floating ips API
   port           EdgeCloud ports API
   snapshot       EdgeCloud snapshots API
   image          EdgeCloud images API
   region         EdgeCloud regions API
   project        EdgeCloud projects API
   keystone       EdgeCloud keystones API
   quota          EdgeCloud quotas API
   limit          EdgeCloud limits API
   cluster        EdgeCloud k8s cluster commands
   pool           EdgeCloud K8s pool commands
   l7policy       EdgeCloud l7policy API
   router         EdgeCloud router API
   fixed_ip       EdgeCloud reserved fixed ip API
   help, h        Shows a list of commands or help for one command

```

Running tests locally
------------------------------------
To run the tests locally using the following command:
```bash
make run_local_tests
```
This command will run the tests and display the output, excluding lines with 'no test files'.

Running linters locally
------------------------------------
To run linters locally, you will need to run a series of checks, including go vet, go fmt, gofumpt, and golangci-lint. 
First, run the following commands to check the code for errors and format it:
```bash
make checks
```
Next, run the linters using the following command:
```bash
make linters
```
This will run the golangci-lint tool and display any issues found in the code.

Note: If you don't have golangci-lint installed, the make linters command will automatically download and install it for you.
