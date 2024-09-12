# Telnet (Transmission Control Protocol)

Telnet provides bi-directional communcation between two temrinal devices or two process (distributed computing)

## Ideas of Telnet
1. NVT -> Network virtual terminal
   This is an immaginery terminal device that mimics what a terminal is like.
2. negotiated options
   NVT provides basic networking features, that might not meet the needs of the participants in a network communcation or process-to-process communcation, negotiated princple gives the pariticpants the ability to agree on extra features. For example if a party wants to have the posiblity of displaying certain characters that is not supported on NVT, the both parties can agree on this, but if one disagrees, the system sticks to the basic NVT features
3. symmetric view of temrinals and processess
   This represents a way in which parties need to coordinate acknowlegements of request to prevent acknowledgemnet loops. 
   a. A party must not send a request for the sake of announcing what mode it is in
   b. if a party receives a request to enter a mode it is already in, the request should not be acknowleged. This prevents ack loops

