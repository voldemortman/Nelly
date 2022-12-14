# Nelly
Open source sdn with easy BPF capabilities, the ability to create bridges locally, remotely with a spine leaf architechture, and the combination of both, a management server for the spines and leaves, and a client/cli tool to manage it

## TODO
- [] Copy code for initial POC
- [] Orginize logger
- [] Create CLI
- [] Make bridge with switches and mac learning
- [] Create server for spine and leaf
- [] Add functionality like latency, pl, bandwith, and ordering
- [] Make filters more managed, with orders of filters and names. Some filters might have to be last and not displayed.
- [] Management server. This will also include heartbeating on the bridges.

## Stalker
A stalker is the component in charge of tunneling all packets going out of a network interface, and is in charge of peeling all incoming traffic heading to the interface.

There will be two types of stalkers. 

### Local Stalker
One will be a local stalker, where the interface being stalked is directly connected to the bridge host.

The local stalker will communicate with the bridge directly, and will have the following configuration:
```
{
  interface: string
  id: UUID
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
}
```

## Bridge
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

## Leaf
A leaf will host remote stalkers. 

### Routes:

- *GET* /remote/ <br>
  **Body**
  ```
  {
    bridgeId?: UUID
  }
  ```
  **Returns** All remote stalkers. If bridge id is provided, returns only stalkers connected to bridgeId.

- *GET* /remote/{stalkerID}/ <br>
  **Returns** Returns the configuration for the remote stalker with stalker id

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
  **Returns** the id of the stalker

- *DELETE* /remote/ <br>
  **Body**
  ```
  {
    bridgeID: UUID
  }
  ```
  Deletes all stalkers that are connected to bridge with bridgeID
  **Returns** a list containing all stalkers with id

- *DELETE* /remote/{stalkerID}/ <br>
  Delete a stalker with id stalkerID

- *PATCH* /remote/{stalkerID}/ <br>
  **Body**
  ```
  {
    remoteAddress: string
  }
  ```
  Make the stalker send all packets to remote address.

## Spine
The spine will house bridges. It can communicate with remote stalkers, or it can also host local stalkers and act just as a switch.

### Routes:

- *GET* /bridge/{bridgeId} <br>
  **Returns** the bridge configuration.

- *GET* /bridge/{bridgeId}/switchingTable <br>
  **Returns** the switching table of bridge with id bridgeId. <br>
  **Except** 409 Conflict - if the bridge is set to hub.

- *GET* /bridge/{bridgeId}/local <br>
  **Returns** all local stalkers connected to the bridge.

- *GET* /bridge/{bridgeid}/local/{stalkerId} <br>
  **Returns** the configuration for a singular local stalker connected to a bridge.

- *GET* /bridge/{bridgeId}/remote <br>
  **Returns** all remote stalkers connected to the bridge.

- *GET* /bridge/{bridgeid}/remote/{stalkerId} <br>
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

- *POST* /bridge/{bridgeId}/local
  **Body**: 
  ```
  {
    interface: string
  }
  ```
  Creates a local stalker that listens on the interface. <br>
  **Returns** the stalker id.

- *POST* /bridge/{bridgeId}/remote
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
  **Returns** the stalker id <br>
  **Except** 409 Conflict - if the port is already in use by another stalker.

- *DELETE* /bridge/{bridgeId} <br>
  Deletes a bridge with bridge id. <br>
  **Except** 409 Conflict - if the bridge has connected stalkers, it will return 409.

- *DELETE* /bridge/{bridgeId}/local/{stalkerId} <br>
  Deletes a local stalker with stalker id on bridge with bridge id.

- *DELETE* /bridge/{bridgeId}/remote/{stalkerId} <br>
  Deletes a remote stalker with stalker id on bridge with bridge id.
  