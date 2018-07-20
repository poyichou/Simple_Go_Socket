# Simple golang socket  
## Server:  
* Accepts up to 10 clients and handles each of them with other goroutine  
* Stores every client information (including closed clients) with map  
	* client address  
	* connection time  
	* online or offline  
* Receive message sent from every client  
	* If the message contains "list":  
		Send client history to the client  
		format:  
		```bash
		======== Client history ========  
		Client: [address] | Connection time: [time] | [online|offline][ <- you|]   
		...  
		======== Online Client: [number] ========  
		```
	* Otherwise:  
		Sent "Received: " + [received message] to the client  
## Client:  
* Connect to server  
* Deal with two event with select:  
	* Local input:  
		* if it is "exit", close the connection and return  
		* otherwise, send it to the server  
	* Message from the server:  
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
