# pkimanager [![Build Status](https://travis-ci.org/thomas-maurice/pkimanager.svg?branch=master)](https://travis-ci.org/thomas-maurice/pkimanager)
A tool to manage your own Private Key Infrastructure

## Introduction
This tool aims at making easier managing your very own PKI, for small Infrastructures
that don't require trusted certificates that are pretty expensive. This is adapted to
programs such as OpenVPN for instance, if you know how a pain in the ass easy-rsa is,
you might want to give it a try.

## Building
Just enter in a prompt:
```
make
```

Alternatively, you can built it using:
```
go get
go build
```

## Features
This program enables you to manage one PKI. The configuration is done in the pkimanager.yml file
which syntax is easy.
```yaml
# Directory in which will be installed the CA
ca_root: /some/path/where/to/install/your/ca
```

The following directories will be created:
* `certificates` for the signed certificates
* `revoked` for the revoked certificates
* `keys` for the certificate private keys

The following files will be created:
* `ca.crt` for the root certificate
* `ca.key` for the root certificate's key
* eventually `ca.crl` for the certificate revokation list

## Creating your CA
You can do it via the command `pkimanager ce create <commonName>`. The following options are available:
* `-c`: Country code
* `-l`: Location
* `-k`: Private key size (defaults to 4096)
* `-v`: Validity in years (default 1)
* `-f`: Force, will overwrite any existing file

And some more, slightly mode exotic are available too, please check the `ca create --help` command for
more information.

Example:
```
$ pkimanager ca create maurice.fr -k 2048 -c FR -l Lille -v 100
INFO[0000] Generating the CA maurice.fr with a 2048 bits key and a 100 years validity
INFO[0000] Generation complete, certificate written to CertificationAutority/ca.crt, and key to CertificationAutority/ca.key
```

## Create a certificate
You can create a certificate with:
```
$ pkimanager certificate create vpn.maurice.fr
INFO[0000] Creating certificates directory              
INFO[0000] Creating keys directory                      
INFO[0000] Generating the certificate vpn.maurice.fr with a 4096 bits key and a 1 years validity
INFO[0002] Generation complete, certificate written to CertificationAutority/certificates/vpn.maurice.fr.crt, and key to CertificationAutority/keys/vpn.maurice.fr.key
```

Basically the same options are available, plus support for IP addresses and alternate
DNS management, check `pkimanager certificate create --help` for more information.

## Revoking a certificate
```
$ pkimanager certificate revoke vpn.maurice.fr
INFO[0000] Certificate revoked, run the crl regen command to make it effective
```

## Regenerate the CRL
```
$ pkimanager crl regen
INFO[0000] Adding certificate vpn.maurice.fr, serial 105791853041754404084114357416071945739 to revokation list
INFO[0000] Regeneration complete, CRL written to CertificationAutority/ca.crl
```

And you are done :)

## License
```
           DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
                   Version 2, December 2004

Copyright (C) 2016 Thomas Maurice <thomas@maurice.fr>

Everyone is permitted to copy and distribute verbatim or modified
copies of this license document, and changing it is allowed as long
as the name is changed.

           DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE
  TERMS AND CONDITIONS FOR COPYING, DISTRIBUTION AND MODIFICATION

 0. You just DO WHAT THE FUCK YOU WANT TO.
```
