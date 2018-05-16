# keyvault-certdeploy

**WORK IN PROGRESS**

**keyvault-certdeploy** is a helper tool used to facilitate X.509 certificate deployment to Linux VMs with [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault/). Unlike the built-in method via VM secrets, it does support RSA and ECDSA certificates, local deployment and update hooks. It can be used to push certificates to Key Vault from a Let's Encrypt deployment hook and to refresh VM certificates on boot or to periodically poll for updates via cronjob. 

## Requirements

* [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault) vault and an account with write permissions to vault secrets (for uploading new certs) and a read-only secret permissions for fetching certs.  

## Limitations

To support importing ECC keys and ECDSA certificates, it does not use Key Vault's internal mechanism for key and certificate management and stores key and certificate in a single PEM-encoded secret. Thus, it cannot be used with HSM-backed keys or Key Vault's automatic certificate renewal.

## Configuration

TBD

## Usage

### Push certificate to vault

TBD

### Dump certificate from vault

TBD

### Fetch/Refresh certificates from vault

TBD

## Build

On Linux:

```
$ mkdir keyvault-certdeploy && cd keyvault-certdeploy
$ export GOPATH=$PWD
$ go get -d github.com/emgag/keyvault-certdeploy
$ cd src/github.com/emgag/keyvault-certdeploy
$ make install
```

will download the source and builds binary called _keyvault-certdeploy_ in $GOPATH/bin.

## Bundled libraries

* [Cobra](https://github.com/spf13/cobra) command line processing 
* [go-homedir](https://github.com/mitchellh/go-homedir) home directory expansion 
* [go-playground/log](https://github.com/go-playground/log)  
* [Microsoft Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go)
* [Viper](https://github.com/spf13/viper) configuration library

## License

keyvault-certdeploy is licensed under the [MIT License](http://opensource.org/licenses/MIT).