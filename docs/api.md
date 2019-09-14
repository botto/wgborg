# API Docs for WG MGR

## Submitting new public key and psk

__POST /add_peer__

*Payload*
```javascript
{
  "PublicKey": "",
  "Psk": "",
  "Name": "",
  "IP": ""
}
```

The IP should be in CIDR format i.e.: 10.10.10.10/24

## Example

*Requirments*
- HTTPie

```bash
http POST 127.0.0.1:8080/add_peer PublicKey=oDs4M1XlPYjEKiPVYnisLHDicBA1vEjr5921TX+b31g= Name=asd Psk=AjC7NUjxJNn9/AQDCGTuPGfWMhBzdCJNszdoxuuybAI= IP=10.10.10.10/24
```
