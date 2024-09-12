# Notes on FTP Protocol RFC

### Commands

1. USER
   This command helps to identify if a user should have accees to file information. After this command it is common to have a `PASS` command for authenticating a user
   _params_
   1. NAME -> the identification of the user

2. PASSWORD
   This command is use for specifying a users password. It should be used immedieately after using the `USER` command. **This should be masked out** > the ftp client should handle this
   _params_
   1. PASS -> the password of the user described with the `USER` command




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
