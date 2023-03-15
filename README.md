# Nelly
Open source sdn with easy BPF capabilities, the ability to create bridges locally, remotely with a spine leaf architechture, and the combination of both, a management server for the spines and leaves, and a client/cli tool to manage it.
If you don't know about spine leaf architechture, you could read about it [here](https://community.fs.com/blog/leaf-spine-with-fs-com-switches.html).

## TODO
- [] Copy code for initial POC.
- [] Orginize logger.
- [] Create CLI.
- [] Make bridge with switches and mac learning.
- [] Create server for spine and leaf.
- [] Add functionality like latency, pl, bandwith, and ordering.
- [] Make filters more managed, with orders of filters and names. Some filters might have to be last and not displayed.
- [] Management server. This will also include heartbeating on the bridges.

## Stalker
---
A stalker is the component in charge of tunneling all packets going out of a network interface, and is in charge of peeling all incoming traffic heading to the interface.<br>
<br>

There will be two types of stalkers:

### Local Stalker
One will be a local stalker, where the interface being stalked is directly connected to the bridge host.

The local stalker will communicate with the bridge directly. The local stalker is used when nelly server is in single mode. Local stalkers will have the following configuration:
```
{
  interface: string
  id: UUID
  isRunning: bool
}
```


### Remote Stalker
The other will be a remote stalker, and will communicate with the bridge with a connected udp socket. The port will be the same number between the bridge and the stalker, and they will be unique across the entire system. This means that the number of stalkers is limited to 2^16 ports. The upside of limiting the ports is that migrating bridge to different spines will be extremly easy and without packet loss. If the port isn't unique between spines, when we migrate to a different spine we can't necessarily have the bridge communicate on the same port number, and then if we change the port number, we either have to deal with packet loss or latency.

The remote stalker will have the following configuration:
```
{
  interface: string
  spineAddress: string
  port: number
  id: UUID
  bridgeId: UUID // This might just be in the DB
  isRunning: bool
}
```

## Bridge
---
The bridge can work either like a switch or a hub. It will have the endpoint of every stalker in its configuration.

The configuration will look like this:
```
{
  id: UUID
  isSwitch: bool
  localAddress: string
  stalkers: [{
    isLocal: true
    id: UUID
    interface: string
  }, {
    isLocal: false
    id: UUID
    leafAddress: string
    port: number
  }]
}
```
<br>

## Nelly Server
---
The nelly server will be in charge of managing the bridges and stalkers. The nelly server has three modes: leaf, spine, and single. Single mode means that the server will have both stalkers and bridges connected to it, and it isn't part of a leaf spine architecture. It will be a grpc server that will only be exposed to the management server. <br>
<br>

### Routes 
Some routes will be exposed only in leaf mode, some only in spine mode, and the rest will be exposed in both modes.

<br>

#### Spine Mode Routes

- *POST* /changeBridgeMode/ <br>
  Change the mode of a bridge to hub if it is a switch, and to a switch if it is a hub.

- *POST* /createBridge/ <br>
  **Body**
  ```
  {
    bridgeId: UUID
  }
  ```
  Create a bridge with id.
  
- *POST* /deleteBridge/ <br>
  **Body**
  ```
  {
    bridgeId: UUID
  }
  ```
  Delete a bridge with specified id.

  - *GET* /getSwitchingTable/ <br>
  **Body**
  ```
  {
    bridgeId
  }
  ```
  Get the switching table of a bridge.
<br>
<br>

#### Leaf Mode Routes
- *POST* /createStalker/ <br>
  **Body**
  ```
  {
    stalkerId: UUID
    interface: string
    spineAddress: string
    port: number
  }
  ```
  Create a remote stalker.

  
- *POST* /killStalker/ <br>
  **Body**
  ```
  {
    stalkerId: UUID
  }
  ```
  Kill a stalker with the specified id. 

- *POST* /pauseStalker/ <br>
  **Body**
  ```
  {
    stalkerId: UUID
  }
  ```
  Pause stalker with the specified id.

- *POST* /unpauseStalker/ <br>
  **Body**
  ```
  {
    stalkerId: uuid
  }
  ```
  Unpause a stalker with the specified id.

#### Single Mode Routes

Single mode will have all the routes of the leaf and of the spine.

## Leaf
---
A leaf has to run 2 containers, the nelly management server, and the nelly server in leaf mode. 

## Nelly Management Server
The nelly management server, will be how you communicate with nelly. It wil be expose a rest server to the public network, and will be able to access systemwide state. It will also be in charge of all leaf state management, and request validation. Any future feature regarding users, accessing a db, or something similar will be added to here.


#### Routes

- *GET* /remote/ <br>
  **Body**
  ```
  {
    bridgeId?: UUID
  }
  ```
  **Returns** All remote stalkers. If bridge id is provided, returns only stalkers connected to bridgeId.

- *GET* /remote/{stalkerID}/ <br>
  **Returns** Returns the configuration for the remote stalker with stalker id.

- *POST* /remote/ <br>
  **Body**
  ```
  {
    spineAddress: string
    port: number
    bridgeId: UUID
  }
  ```
  Creates a stalker that sends all traffic to spineAddress:port(udp) and listens on port(udp). <br>
  This should be called after the stalker was created on the bridge, so that there will be no packet loss. <br>
  **Returns** the id of the stalker.

- *DELETE* /remote/ <br>
  **Body**
  ```
  {
    bridgeID: UUID
  }
  ```
  Deletes all stalkers that are connected to bridge with bridgeID.
  **Returns** a list containing all stalkers with id.

- *DELETE* /remote/{stalkerID}/ <br>
  Delete a stalker with id stalkerID.

- *PATCH* /remote/{stalkerID}/ <br>
  **Body**
  ```
  {
    remoteAddress: string
  }
  ```
  Make the stalker send all packets to remote address.

## Spine
---
The spine will house bridges. It can communicate with remote stalkers, or it can also host local stalkers and act just as a switch. Similar to the leaf, it will have a nelly management server on one container, and a nelly server in spine mode on another container.<br>
<br>

### Nelly server:
The nelly server will be in charge of managing the bridges and the local stalkers. It will be a grpc server that will only be exposed to the management server.


#### Routes:

- *GET* /bridge/{bridgeId} <br>
  **Returns** the bridge configuration.

- *GET* /bridge/{bridgeId}/switchingTable/ <br>
  **Returns** the switching table of bridge with id bridgeId. <br>
  **Except** 409 Conflict - if the bridge is set to hub.

- *GET* /bridge/{bridgeId}/local/ <br>
  **Returns** all local stalkers connected to the bridge.

- *GET* /bridge/{bridgeid}/local/{stalkerId}/ <br>
  **Returns** the configuration for a singular local stalker connected to a bridge.

- *GET* /bridge/{bridgeId}/remote <br>
  **Returns** all remote stalkers connected to the bridge.

- *GET* /bridge/{bridgeid}/remote/{stalkerId}/ <br>
  **Returns** the configuration for a singular remote stalker connected to a bridge.

- *POST* /bridge/ <br>
  **Body**:
  ```
  {
    type: 'CreateBridgeDTO' | 'MigrateBridgeDTO'
    params: CreateBridgeDTO | MigrateBridgeDTO
  }

  CreateBridgeDTO: {
    isSwitch: bool
  }

  MigrateBridgeDTO: {
    id: UUID
    ...CreateBridgeDTO
  }
  ```
  Creates a bridge. <br> 
  **Returns** the bridge id.

- *POST* /bridge/{bridgeId}/local/
  **Body**: 
  ```
  {
    interface: string
  }
  ```
  Creates a local stalker that listens on the interface. <br>
  **Returns** the stalker id.

- *POST* /bridge/{bridgeId}/remote/
  **Body**: 
  ```
  {
    type: 'CreateRemoteStakerDTO' | 'MigrateRemoteStalkerDTO'
    params: CreateRemoteStalkerDTO | MigrateRemoteStalkerDTO
  }

  CreateRemoteStalkerDTO: {
    leafAddress: string
    port: number
  }

  MigrateRemoteStalkerDTO: {
    UUID
    ...CreateRemoteSTalkerDTO
  }
  ```
  Create a remote stalker that listens to incoming transmition on the udp port, and will send packets to a stalker at the leafAdress that listens on the specified port. <br>
  **Returns** the stalker id. <br>
  **Except** 409 Conflict - if the port is already in use by another stalker.

- *DELETE* /bridge/{bridgeId}/ <br>
  Deletes a bridge with bridge id. <br>
  **Except** 409 Conflict - if the bridge has connected stalkers, it will return 409.

- *DELETE* /bridge/{bridgeId}/local/{stalkerId}/ <br>
  Deletes a local stalker with stalker id on bridge with bridge id.

- *DELETE* /bridge/{bridgeId}/remote/{stalkerId}/ <br>
  Deletes a remote stalker with stalker id on bridge with bridge id.

- *PATCH* /bridge/{bridgeId}/ <br>
  **Body**
  ```
  {
    mode: "switch" | "hub"
  }
  ```
  Change bridge mode to switch or to hub.
  