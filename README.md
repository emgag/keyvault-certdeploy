# keyvault-certdeploy

[![Build Status](https://travis-ci.org/emgag/keyvault-certdeploy.svg?branch=master)](https://travis-ci.org/emgag/keyvault-certdeploy)
[![Go Report Card](https://goreportcard.com/badge/github.com/emgag/keyvault-certdeploy)](https://goreportcard.com/report/github.com/emgag/keyvault-certdeploy)

**WORK IN PROGRESS**: more or less feature complete, but not used in production yet.

**keyvault-certdeploy** is a helper tool used to facilitate X.509 certificate deployment to Linux VMs with [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault/). Unlike the built-in method via VM secrets, it does support RSA and ECDSA certificates, local deployment and update hooks. It can be used to push certificates to Key Vault from a Let's Encrypt deployment hook and to refresh VM certificates on boot or to periodically poll for updates via cronjob. 

## Requirements

* [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault) vault and an account with write permissions to vault secrets (for uploading new certs) and a read-only secret permissions for fetching certs. 
* Go >=1.10 for building.  

## Limitations

To support importing ECC keys and ECDSA certificates, it does not use Key Vault's internal mechanism for key and certificate management and stores key and certificate in a single PEM-encoded secret. Thus, it cannot be used with HSM-backed keys or Key Vault's automatic certificate renewal.

## Authentication

keyvault-certdeploy uses the [automatic environment authentication from the Azure Golang SDK](https://docs.microsoft.com/en-us/go/azure/azure-sdk-go-authorization#use-environment-based-authentication) to authenticate itself to Key Vault. Credentials can be either provided through environment variables or, if running in a Azure VM, through [Managed Service Identity (MSI)](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/overview). No environment variables need to be set when using MSI.

Required access policies to vault:
* To push certificates: Secret Get & Set
* To fetch (dump & sync) certificates: Secret Get


## Configuration

See [keyvault-certdeploy.yml.dist](keyvault-certdeploy.yml.dist) for a sample configuration. 



## Usage

```
X.509-Certificate deployment helper for Azure Key Vault

Usage:
  keyvault-certdeploy [command]

Available Commands:
  dump        Dump certificate and key from vault to current directory or dir, if supplied
  help        Help about any command
  push        Push a certificate to the vault
  sync        Sync configured certificates from vault to system

Flags:
  -c, --config string   Config file (default locations are $HOME/.config/keyvault-certdeploy.yml, /etc/keyvault-certdeploy.yml, $PWD/keyvault-certdeploy.yml)
  -h, --help            help for keyvault-certdeploy
  -q, --quiet           Be quiet
  -v, --verbose         Be more verbose
```

## Build

On Linux:

```
$ mkdir keyvault-certdeploy && cd keyvault-certdeploy
$ export GOPATH=$PWD
$ go get -d github.com/emgag/keyvault-certdeploy
$ cd src/github.com/emgag/keyvault-certdeploy
$ make install
```

will download the source and build binary called _keyvault-certdeploy_ in $GOPATH/bin.

## Bundled libraries

* [Cobra](https://github.com/spf13/cobra) command line processing 
* [go-homedir](https://github.com/mitchellh/go-homedir) home directory expansion 
* [go-playground/log](https://github.com/go-playground/log) enhanced logging
* [Microsoft Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go)
* [Viper](https://github.com/spf13/viper) configuration library

## License

keyvault-certdeploy is licensed under the [MIT License](http://opensource.org/licenses/MIT).