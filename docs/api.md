# API Docs for WG MGR

## Submitting new public key and psk

__POST /add_peer__

*Payload*
```javascript
{
  "public_key": "The public key of the peer connecting",
  "psk":        "PSK used between the peers",
  "name":       "Name of the peer",
  "ip":         "CIDR notation ip address",
  "network_id": "The network to associate this peer with",
}
```

The IP should be in CIDR notation i.e.: 10.10.10.10/24

__POST /add_network__

*Payload*
```javascript
{
  "private_key": "The private key of the WG interface",
  "name":        "Name of the interface and network",
  "port":        "Port to listen on, this should be an int",
  "ip":          "IP Address the interface will have, CIDR notation", 
}
```

The IP should be in CIDR notation i.e.: 10.10.10.1/24

## Example

*Requirments*
- HTTPie

```bash
http POST 127.0.0.1:8080/add_peer PublicKey=oDs4M1XlPYjEKiPVYnisLHDicBA1vEjr5921TX+b31g= Name=asd Psk=AjC7NUjxJNn9/AQDCGTuPGfWMhBzdCJNszdoxuuybAI= IP=10.10.10.10/24
```
