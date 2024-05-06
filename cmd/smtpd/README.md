# smtpd 

A local SMTP relay server

## Compilation

```console
go build -o smtpd cmd/smtpd/smtpd.go
```

## Usage

```console 
./smtpd --token <mail-api-token> --environment <env>
```