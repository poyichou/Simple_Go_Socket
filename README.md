# Simple golang socket  
## Server:  
* Accepts up to 10 clients and handles each of them with other goroutine  
* Stores every client information (including closed clients) with map  
	* client address  
	* connection time  
	* online or offline  
* Receive message sent from every client  
	* if the message contains "list":  
		send client history to the client  
		format:  
		======== Client history ========  
		Client: [address] | Connection time: [time] | [online|offline][ <- you|]   
		...  
		======== Online Client: [number] ========  
	* otherwise:  
		sent "Received: " + [received message] to the client  
## Client:  
* Connect to server  
* Deal with two event with select:  
	* local input:  
		* if it is "exit", close the connection and return  
		* otherwise, send it to the server  
	* message from the server:  
		simply print out to screen  

## To run:  
Run server first before running any client  
Server:  
```bash
cd go_socket_server
go run go_socket_server.go
```
Client:  
```bash
cd go_socket_client
go run go_socket_client.go
```
