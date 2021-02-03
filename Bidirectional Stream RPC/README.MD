This example shows you how to build a Client Stream RPC service which works by creating a stream from the client to the server, the idea is to send an video in chunks and the server needs to rebuild this file. I've tested sending files greater than 50MB of size. The project is structured into 4 folders mainly: 

1. **rpc** folder contains **proto** folder which is where our protofiles are located and the **videopb** folder where we will have our compiled protofiles.

2. **api** folder contains a server/client rpc application, just run the *server/server.go* file first and after the *client/client.go* file.

3. **assets/server** folder contains videos that were received from the client.

4. **assets/client** folder contains videos to upload to the sever.

If you want to build by yourself the protofile(located inside **rpc/proto** folder), run ***make build*** command and you'll see the output inside **videopb** folder, but first of all be sure you have protocol buffers compiler and grpc plugin already setup, see this [link](https://github.com/cerezo074/GRPC/tree/main/Code%20Generation) if you want to know how to do it. 

***Note: Remember that this example uses Go Modules so you don't need to worry about externals dependencies those dependencies are downloaded automatically and you don't need to execute any go file inside your go path 😬.***