# About

This tool is a part of Tönnjes developer challenge.
The purpose of this CLI tool is to display relevant information regarding X.509 certificates located in the IDeTRUST 
credential repository at https://idetrust.com.

## Build

To build the project all you have to do is run the build script via the following command:

```shell
$ ./build.sh
```

The script will build the application for: linux, windows_x32, and windows_x64.

## Usage

You can run the "certinfo" executable with no arguements:

```shell
$ ./certinfo
```

Providing no arguement will make the program use the following default values:

 - host: "https://idetrust.com"

 - daid: "QC-DEMO"

 - cid: 3

 - skip-root: false (this flag can be used to skip the validation for the root certificate)

and we end up with this URL: https://idetrust.com/daid/QC-DEMO/cid/3

It is possible to change any of the previously listed values when running the application.

```shell
$ ./certinfo --cid 2 --daid "lorem ipsum" --host https://example.com --skip-root
```

This example will target the following URL: https://example.com/daid/lorem%20ipsum/cid/2
Note that this is just an arbitrary example and it will throw an error.

```shell
$ ./certinfo --cid 1 --skip-root
```

The target URL here is: https://idetrust.com/daid/QC-DEMO/cid/1
We use the skip root flag because the root certificate is self-signed and it is not present in the system trust store.

The output should be a valid JSON object containing relevant information for each certificate.

```json
[
  {
    "not_before": "2026-01-08T07:34:23Z",
    "not_after": "2027-01-08T07:34:23Z",
    "issuer": "CN=QC DigSig Demo QC-DEMO https://www.idetrust.io 2026,O=QC DigSig Demo Inc.,L=Delmenhorst,C=DE",
    "subject": "CN=https://idetrust.com/daid/QC%20DEMO/cid/1,O=IDeTRUST GmbH,C=DE",
    "is_ca": false,
    "serial_number": "265611437842169415274341616385990281489",
    "fingerprint_sha256": "31f227f04ddad3497cfd6492d930eed57f1b2d48e53b67892b2eb709ef4ba5db"
  },
  {
    "not_before": "2025-11-13T00:00:00Z",
    "not_after": "2046-12-31T23:59:59Z",
    "issuer": "CN=DigSig CA IDeTRUST 2026,O=IDeTRUST GmbH,ST=Lower Saxony,C=DE",
    "subject": "CN=QC DigSig Demo QC-DEMO https://www.idetrust.io 2026,O=QC DigSig Demo Inc.,L=Delmenhorst,C=DE",
    "is_ca": true,
    "serial_number": "284119739675548703326845937359568515688255311725",
    "fingerprint_sha256": "3f3901fabbf442a637b05dbf3eea3bdd8bffeadbb563d2d0fcef6a560fe7ca3b"
  },
  {
    "not_before": "2025-10-07T00:00:00Z",
    "not_after": "2046-12-31T23:59:59Z",
    "issuer": "CN=DigSig CA IDeTRUST 2026,O=IDeTRUST GmbH,ST=Lower Saxony,C=DE",
    "subject": "CN=DigSig CA IDeTRUST 2026,O=IDeTRUST GmbH,ST=Lower Saxony,C=DE",
    "is_ca": true,
    "serial_number": "412711715293269178734742633292716723066828337464",
    "fingerprint_sha256": "2c2abeaebe54bb32b31954f0c8e29c4b078a6ada2da1ec9248dc9775d57dcf32"
  }
]
```
