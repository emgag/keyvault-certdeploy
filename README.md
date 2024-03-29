# keyvault-certdeploy

![build](https://github.com/emgag/keyvault-certdeploy/workflows/build/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/emgag/keyvault-certdeploy)](https://goreportcard.com/report/github.com/emgag/keyvault-certdeploy)
[![Docker Pulls](https://img.shields.io/docker/pulls/emgag/keyvault-certdeploy.svg)](https://hub.docker.com/r/emgag/keyvault-certdeploy)


**keyvault-certdeploy** is a tool used to facilitate X.509 certificate deployment to Linux systems with [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault/), supporting RSA and ECDSA certificates, local deployment and update hooks. Originally built to work around shortcomings in the Azure certificiate provisioning via VM secrets, it's possible to provision systems outside of Azure as well, provided they have access to the corresponding vault. It can be used to push certificates to Key Vault from a Let's Encrypt deployment hook and to refresh VM certificates on boot or to periodically poll for updates via cronjob.

## Requirements

* [Azure Key Vault](https://azure.microsoft.com/en-us/services/key-vault) vault and an account with write permissions to vault secrets (for uploading new certs) and a read-only secret permissions for fetching certs. 
* Go >=1.19 and [goreleaser](https://goreleaser.com/) for building.  

## Limitations

To support importing ECC keys and ECDSA certificates, it does not use Key Vault's internal mechanism for key and certificate management and stores key and certificate in a single PEM-encoded secret. Thus, it cannot be used with HSM-backed keys or Key Vault's automatic certificate renewal.

## Authentication

keyvault-certdeploy uses the [automatic environment authentication from the Azure Golang SDK](https://docs.microsoft.com/en-us/go/azure/azure-sdk-go-authorization#use-environment-based-authentication) to authenticate itself to Key Vault. Credentials can be either provided through environment variables or, if running in a Azure VM, through [Managed Service Identity (MSI)](https://docs.microsoft.com/en-us/azure/active-directory/managed-service-identity/overview). No environment variables need to be set when using MSI.

Required access policies to vault:
* To push certificates: Secret Get & Set
* To fetch (dump & sync) certificates: Secret Get
* To list certificates: Secret Get & List
* To delete certificates (delete & prune): Secret Delete, Get & List 

## Install 

Either download a [compiled release](https://github.com/emgag/keyvault-certdeploy/releases), use the [emgag/keyvault-certdeploy](https://hub.docker.com/r/emgag/keyvault-certdeploy) docker image or see build instructions in this README.

Also checkout example [ansible role](https://github.com/emgag/keyvault-certdeploy/tree/master/examples/ansible) for deploying keyvault-certdeploy and individual certificates or the [systemd unit files](https://github.com/emgag/keyvault-certdeploy/tree/master/examples/systemd) for running keyvault-certdeploy as an onboot- and/or periodic service.

## Configuration

See [keyvault-certdeploy.yml.dist](keyvault-certdeploy.yml.dist) for a sample configuration. 

## Usage

```
X.509-Certificate deployment helper for Azure Key Vault

Usage:
  keyvault-certdeploy [command]

Available Commands:
  delete      Deletes certificate from vault
  dump        Dump certificate and key from vault to current directory or dir, if supplied
  help        Help about any command
  list        List certificates in vault
  prune       Remove expired certificates from vault
  push        Push a certificate to the vault
  sync        Sync configured certificates from vault to system

Flags:
  -c, --config string   Config file (default locations are $HOME/.config/keyvault-certdeploy.yml, /etc/keyvault-certdeploy/keyvault-certdeploy.yml, $PWD/keyvault-certdeploy.yml)
  -h, --help            help for keyvault-certdeploy
  -q, --quiet           Be quiet
  -v, --verbose         Be more verbose
      --version         version for keyvault-certdeploy

Use "keyvault-certdeploy [command] --help" for more information about a command.
```

### delete

```
Deletes certificate from vault

Usage:
  keyvault-certdeploy delete <subject> <keyalgo> [flags]

Flags:
  -h, --help   help for delete
  -y, --yes    Don't confirm before deleting
```

This command removes a single certificate from vault. 

### dump

```
Dump certificate and key from vault to current directory or dir, if supplied

Usage:
  keyvault-certdeploy dump <subject> <keyalgo> [flags]

Flags:
      --cert string               Name of the leaf certificate file (default "cert.pem")
      --chain string              Name of the certificate chain file (default "chain.pem")
  -d, --dir string                Directory to save files to (default ".")
      --fullchain string          Name of the full certificate chain file (default "fullchain.pem")
      --fullchainprivkey string   Name of the full certificate chain + private key file (default "fullchain.privkey.pem")
  -h, --help                      help for dump
      --key string                Name of the private key file (default "privkey.pem")
```

The files generated by this command follow the same convention as common Let's Encrypt utilities like [certbot](https://github.com/certbot/certbot) or [dehydrated](https://github.com/lukas2511/dehydrated):

* `privkey.pem`  : the private key for the certificate.
* `fullchain.pem`: the certificate file containing all CA certificates as well as the leaf certificates. That's usually the one to use in TLS configurations of the services (e.g. in nginx's `ssl_certificate`)
* `chain.pem`    : the CA certificates only, e.g. for OCSP stapling.
* `cert.pem`     : just the leaf certificate. 

See [certbot's documentation about this](https://certbot.eff.org/docs/using.html#where-are-my-certificates) for more info.

Additionally, following files will be generated as well:
* `fullchain.privkey.pem` : the concatenation of fullchain and privkey.

### list

```
List certificates in vault

Usage:
  keyvault-certdeploy list [flags]
```

This command will list all certificates in configured vault.

### prune

```
Remove expired certificates from vault

Usage:
  keyvault-certdeploy prune [flags]

Flags:
  -d, --days int   Delete certificates after this many days of being expired (default 7)
  -h, --help       help for prune
  -n, --noop       Just list expired certificates, don't actually remove the certs
  -y, --yes        Don't confirm before pruning
```

Automatically remove all certificates which expired _-d_ days ago from vault.  

### push

```
Push a certificate to the vault

Usage:
  keyvault-certdeploy push <privkey.pem> <fullchain.pem> [flags]
```

To push a certificate to the vault, the push command requires a key- and a certificates file. It's designed to work with the PEM files generated by Let's Encrypt utilities (privkey.pem & fullchain.pem).

### sync

```
Sync configured certificates from vault to system

Usage:
  keyvault-certdeploy sync [flags]

Flags:
  -f, --force     Force update even if version on disk matches the one in vault
  -h, --help      help for sync
      --nohooks   Disable running hooks after cert update
```

This command will fetch all certificates configured in the `certs` list in the config file and run hooks if a certificate was updated. Hooks will only run once per sync run, after all certificates are processed. Duplicate hooks are ignored.

## Vault format

Certificates are pushed to the vault as an unencrypted, single PEM-formated file containing the ECC or RSA private key, the chain certificates and the leaf certificate. The secret is named SubjectCN-PublicKeyAlgo (all lowercase, e.g. example.org-rsa), its content type is set to _application/x-pem-file_ and following tags are defined:

* fingerprint: SHA256 fingerprint (hex-encoded string)
* keyalgo: Either _ECDSA_ or _RSA_
* notafter: UNIX timestamp of the certificate expire date
* subjectcn: The certificate's subject common name


## Build

### On Linux

```
$ git clone github.com/emgag/keyvault-certdeploy
$ cd keyvault-certdeploy
$ make snapshot
```

will download the source and build binary called _keyvault-certdeploy_ in ./dist directory and a locally tagged docker image.

## License

keyvault-certdeploy is licensed under the [MIT License](http://opensource.org/licenses/MIT).
