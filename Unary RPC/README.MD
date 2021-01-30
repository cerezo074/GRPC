This example shows you how to build a Unary RPC server that works like Rest(1 Request -> 1 Response), the idea is to send a text from any language and this web service can tell you what language it belongs to, this project is structured into 2 folders: 

1. **Demo** folder shows you a simple app in which you can interact with the [Language Detector library](https://github.com/abadojack/whatlanggo), you just need to run *main.go* golang file. 

2. **App** folder named contains a server/client RPC application, just run the *server.go* file first and after the *client.go* file.

If you want to build by yourself the protofile located inside **rpc/proto** folder, inside rpc folder run ***make build*** command and you'll see the output inside **languagepb** folder, but first of all be sure you have protocol buffers compiler and GRPC plugin already setup, see this [link](https://github.com/cerezo074/GRPC/tree/main/Code%20Generation) if you want to know how to do it. 

***Note: Remember that this example uses Go Modules so you don't need to worry about externals dependencies those dependencies are downloaded automatically and you don't need to execute any go file inside your go path 😬.***