# Notes on FTP Protocol RFC

### Commands

When ftp command crosses the wire, the default port it uses is port 21. FTP also makes use of port 20 for the FTP Data channel by default.

Most ftp data sessions don't use port 20 for data transfer. Invidual control and data channels are used separat3eely from large file transfers.

The PORT command can be use to specify what port the server should sent data to.

When client is on active mode FTP, this does not work well if the client is behind a firewall or NAT is used ont he clients network. Passive mode ftp solves this problem, by letting the client send a `PASV` command to the server. The server returns a high port normal (normally greater than 1024) that the client used use for data connection

In summary, the PORT command is used in FTP to communicate the TCP port number to use for the data transfer channel. In active mode FTP, the client uses the PORT command to tell the server which high-numbered port the client will use for the data channel, and the server opens a connection to that port. In passive mode, the PASV command is sent by the client, and the server responds with the high-numbered port on which it will accept the data connection


If passive mode is turned on (default), ftp will send a PASV command for all data connections instead of a PORT command. The PASV command requests that the remote server open a port for the data connection and return the address of that port. The remote server listens on that port and the client connects to it. When using the more traditional PORT command, the client listens on a port and sends that address to the remote server, who connects back to it. Passive mode is useful when using ftp through a gateway router or host that controls the directionality of traffic. (Note that though FTP servers are required to support the PASV command by RFC 1123, some do not.)

1. USER
   This command helps to identify if a user should have accees to file information. After this command it is common to have a `PASS` command for authenticating a user
   _params_
   1. NAME -> the identification of the user

2. PASSWORD
   This command is use for specifying a users password. It should be used immedieately after using the `USER` command. **This should be masked out** > the ftp client should handle this
   _params_
   1. PASS -> the password of the user described with the `USER` command

3. PORT
   This command allows the client to send a high port number that should be use for data transfer to the server
   _params_
   1. PORT -> A combination of the client ip (32bit) and port (16bit) separated by , in this format `h1,h2,h3,h4,p1,p2` 
   
      For exmaple `192,168,0,1,20,30`
      client network address -> 192,168,0,1
      port -> 20,30 which is in 16bits and can be translated to (P1 * 256) + P2 =>
      P1 = higher byte(most significant byte) = 20
      P2 - Lower byte (least significnat byte) = 30
      
      (20 * 256) + 30 = 5150
      


### Replies

FTP replies consists of a reply code (3 digits) and some text/
1. Reply Code
   The reply code is use to determine what state the ftp client should be placed in.
2. The text
   The text is meant for the human user using the ftp client
   
```bash
REPLY_CODE <SP> TEXT\r\n
```
   
The `\r\n` represents the telnet EOL code that needs to be used to signal the end of a response/request text



### References

https://userpages.umbc.edu/~dgorin1/451/OSI7/dcomm/ftp.htm#:~:text=The%20control%20connection%20is%20used,to%20actually%20send%20a%20file.

telnet (mainly for control connections process-process communication): https://datatracker.ietf.org/doc/html/rfc854
