## devfeel/connx

#### Version 0.3
* New Feature: add Server.AddConnection, use to add new Connection
* New Feature: add Server.RemoveConnection, use to remove Connection with ConnIndex
* New Feature: add Server.GetConnectionCount, use to get connection count on current Server
* New Feature: add Server.GetConnectionMap, use to get connection map on current Server
* Add Debug log on Server.AddConnection\RemoveConnection\Start\Stop
* 2018-07-24 15:00

#### Version 0.2
* Add connx.SetHeadFlag used to set check head-flag on global mode
* 2018-06-19 15:00

#### Version 0.1
* init version
* feature:
* Server/Client Mode
* NewClient: New client which can send and read message with remoteAddr and OnConnHandle
* NewRequestOnlyClient: new client which can only send message with remoteAddr
* default Message protocol
* 2018-06-19 15:00