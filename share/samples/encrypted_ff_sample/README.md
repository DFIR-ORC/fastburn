# Sample of encrypted FastFind result

## File format

Encrypted archives are PKCS7(CMS) DER encoded files.

But the encrypted archived is not pure 7zip, it is an ORC specific version to allow for a streamed creation.
It needs to be "unstreamed" to get the inner standard 7zip file.

A decryption/unstreaming tool is available at: https://github.com/DFIR-ORC/orc-decrypt



## Manual decryption

If it is smaller than 2Gbyte, the archive can be manually decrypted using openssl CLI.

Exemple:
```sh

openssl cms --decrypt 
    -in ORC_WorkStation_W11-22000-51_FastFind.7z.p7b 
    -out ORC_WorkStation_W11-22000-51_FastFind.7zs
    -inkey mastecontoso.com.key

```

# Files MD5 hashes

To properly test samples processing, the following hashes can be checked

```
2a4e97c2eb6ed77096b76ae02cf37446  ORC_WorkStation_W11-22000-51_FastFind.7z.p7b
5cfd91a98c9f648d6f74fbbe743855e0  ORC_WorkStation_W11-22000-51_FastFind.7zs
13c648b48ef995e1f31586def4b63eee  ORC_WorkStation_W11-22000-51_FastFind.7z
```


