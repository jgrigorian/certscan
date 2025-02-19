# Certscan

`certscan` is a cli tool that let's you interact with the TLS secrets in a kubernetes cluster. By default,
it uses the current context in your kubeconfig file.

- [Usage](#Usage)
  - [Options](#options)
- [Examples](#Examples)

## Usage
```shell
$ certscan --help
NAME:
   certscan - Tool to interact with tls secrets in a kubernetes cluster

USAGE:
   certscan [global options] command [command options]

COMMANDS:
   list     Option for listing certificates
   show     Option for showing certificate details
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help
```

### Options
#### List
```shell
$ certscan list certificates --help
NAME:
   certscan list certificates - list certificates

USAGE:
   certscan list certificates [command options]

OPTIONS:
   --expiring, -e               List only expired or expiring certificates (default: false)
   --namespace value, -n value  Desired namespace (Example: default)
   --all-namespaces, -A         All namespaces (default: false)
   --help, -h                   show help
```

#### Show
```shell
$ certscan show certificate --help
NAME:
   certscan show certificate - show certificate

USAGE:
   certscan show certificate [command options]

OPTIONS:
   --secret value, -s value     Name of the desired secret. Example: worldofgoodbrandsio-ssl-secret
   --namespace value, -n value  Desired namespace (Example: default). NOTE: If no namespace is provided, default is used.
   --help, -h                   show help
```


## Examples
### List certificates in all namespaces
```shell
$ certscan list certificates -A
   Secret Name                                                         Namespace                           Expiration     Issuer            Days Remaining

   amazon-cloudwatch-observability-controller-manager-service-cert     amazon-cloudwatch                   2034-07-15                       3623
   jenkins-ssl-secret                                                  build                               2024-09-08     Let's Encrypt     26
   staging-ssl-secret                                                  default                             2024-09-18     DigiCert Inc      36
   aws-load-balancer-tls                                               kube-system                         2034-08-05                       3644
```

### List certificates in a specific namespace
```shell
certscan list certificates -n default
   Secret Name                        Namespace     Expiration     Issuer            Days Remaining

   staging-ssl-secret                 default       2024-09-18     DigiCert Inc      36
```

### List expiring certificates in a specific namespace
```shell
certscan list certificates -n default --expiring
   Secret Name                   Namespace     Expiration     Issuer           Days Remaining

   staging-ssl-secret            default       2024-09-18     DigiCert Inc     4
```
### Show certificate details
```shell
certscan show certificate -s staging-ssl-secret -n default
   staging-ssl-secret

   Namespace                          default
   Valid From                         2024-08-05
   Valid Until                        2024-11-03
   Issuer                             Let's Encrypt
   Days Remaining                     82
   Subject Alternate Name (DNS)       *.stg.acme.com, stg.acme.com
```