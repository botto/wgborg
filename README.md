# Wireguard Manager

This is a tool to manage a wireguard server.  
It's not complete and is an initial concept.  
The project is going to be considered unstable as long as the main wiregaurd project is.  
Once wireguard is stable then this project will try to get to a V1 as soon as possible.  

## Goals
- Add and remove peers without restarting the server
- Do not rely on local config file (store all config in a DB)
- Let a peer submit a public key and receive back the server connection details
- Webui to manage the interfaces on your server
- Multi backend support

## Long term
- Manage wireguard on remote nodes
- Managment distributed across multiple servers
- Mesh (unless offical support comes out later)
- Visually manamge routes between peers
- Create groups of users that can communicate with each other


## To get started
Grab a copy of the code
```
git clone git@github.com/botto/wg_mgr.git
```

Set up deps (this includes some dummy data)
```
docker-compose up -d
```

Run the server
```
make debug
```

This will run the bin under sudo and sets up an interface with peers.

After you kill the server you will have clean up:
```
sudo ip link del TestNet
```

## Todo for an initial working version
- Clean up after ourselves
- Switch to bbolt db (want less deps)
- Set up remaing api endpoints


## Disclaimer
This project isn't officaly supported or part of the WireGuard project by Jason A. Donenfeld  
