# Content Provider
This repository contains a simple Content Provider that can be used for storing blocks of content. These blocks can be retrieved using their CID, which is returned by the provider itself when a new block is added to it.
A basic http server and a CLI are available for running the provider, storing and requesting content.

## Dependencies
Go(lang).
The repository has been tested using v 1.15 but older versions should work too.

## Install
Run
```
make install
```
to compile the code into an executable called `provider` and install it on the Go install path.
The executable exposes the following interface via the command line
```
> provider

Usage:
  provider [command]

Available Commands:
  client      The client interface for interacting with a provider
  help        Help about any command
  server      Runs a content provider behind an HTTP server

Flags:
  -h, --help   help for provider

Use "provider [command] --help" for more information about a command.
```
## Usage
Once the executable has been compiled the easiest way to interact with it is to
1. Start a new local server to use for exposing the provider functionalities.

```
> provider server

server running at localhost:9999
```
While this command will run using default values, it is possible to modify some of its behaviours via
provided flags. Use the `-h` flag for help on how to configure the http port and the provider maximum capacity.

2. Once the server is running, the same executable can be used to interact with it.
In a separate command line panel run
```
> provider client add --item=hello
```
To add a new object (in this case the string `hello`). The provider will return the CID generated for the provided item.
The CID can be used to retrieve some item information such as its path, hits count etc
```
> provider client get --cid=<the-cid>
```
Issuing these commands will print out the result of the operation formatted as JSON.

## Repo structure and design considerations
The core logic of this package is subdivided into modules available inside the /pkg directory.

### Storage module
The Storage module exposes the `Storage` interface for adding, removing and retrieving raw content from the underlying storage. This module is mainly concerned with how and where raw data is persisted. Concrete implementations are not directly responsible for handling the eviction of content in case of overflow. Instead they will simply return an `ErrNoStorageAvailable` error type when this happens further up the chain.

The `BlockStore` implementation used by the content provider is based on a merkel DAG representation of the available blocks of content backed by an in-memory datastore.
This was my attempt at bridging together content and its representation as a tree of linked nodes whose identifiers can be used to verify the actual raw content of each node.
Even if this resulted in adding another level of indirection, a benefit of this implementation is that the underlying block storage can be replaced by any type implementing the ipfs `Blockstore` interface (such as a wrapper layer around a database, the file system etc), while delegating and abstracting away nodes management to the `DAGService`.

The `SimpleStore` type also included in the /storage folder was the first concrete implementation of the `Storage` interface. This version does not provide a tree-like layout of the available data but simply stores content in a key value map.

The choice of using the `BlockStore` implementation by default felt a more appropriate attempt at providing a solution based on content addressing rather than location addressing.

### Indexer module
The Indexer module is responsible for managing the provider storage state. It’s the layer between the storage module and IO operations and provides functionalities for efficiently fetching data and metadata as well as managing content eviction in case there is no available storage space.
The main functionalities of the `Indexer` interface are:
- the retrieval of content given its CID from the underlying storage layer
- the storage of blocks in the underlying storage layer

Its `MemIndexer` implementation uses an in memory key-value map for associating the CID of individual blocks of content to metadata specific to the content itself.
This metadata includes:
- Block Address
- Hits count (the number of times the block has been queried)
- Block size
- Block creation timestamp
- Last block retrieval timestamp

The internal usage of a map has been chosen for providing fast access to content given its identifier. Using CIDs as keys ensure that only unique blocks of content will be stored at any given time.

<<<<<<< HEAD
**Storage Size Management**
The `MemIndexer` storage size is naively measured in numbers of blocks being stored. While this option provides a simple mechanism to test against, it should be replaced by using the actual file size of the various blocks of content being committed to storage.
As new blocks of content are submitted by clients, the `MemIndexer` will check the current storage size against its maximum capacity and when the two values are equal the block with the lowest number of hits will be removed from both the key-value map and the underlying data storage.
The `Evictor` interface used by the `MemIndexer` struct allows for more sophisticated algorithms to be used for removing least frequently accessed items.

### Provider
The Provider module mediates access with the content Indexer and it is the outer access layer that is being made available through an HTTP server. This layer could be potentially used for providing access to a UI, integrating monitoring or auditing middlewares and more.

## Further considerations
While exploring the IPLD and IPFS ecosystem I considered the idea of representing the content available to the indexer using linked data, where the root node would be represented by the hash of the content of its children (EG files).
The advantage of this representation is the verifiability of the Indexer content integrity as well as the possibilty of traversing and searching the resulting linked tree given any node CID.
Due to my lack of experience - and probably limited understanding - in working with the existing libraries my attempts fell short of a working implementation.
=======
### Provider module
The Provider module mediates access to the content Indexer. It is the outer layer which is made available to external clients through an HTTP server. This layer could potentially be used for  integrating monitoring or auditing middlewares and more.

## Storage Size Management
The storage size is naively measured in numbers of blocks being stored. While this option provides a simple mechanism to test against, it should be replaced by using the actual storage size of the various blocks of content being committed.
As new blocks of content are submitted by clients, the `MemIndexer` will check the storage layer current size against its maximum capacity and when the two values are equal the block with the lowest number of hits will be removed from both the key-value map and the underlying data storage.
The `Evictor` interface used by the `MemIndexer` type allows for more sophisticated algorithms to be used for removing least frequently accessed items.

## Further Considerations
While the proposed solution is the result of applying newly absorbed concepts around content addressing to storage and retrieval, it should be seen as a working prototype rather than a ready-to-go solution. Few topics which haven't been explored here can be incorporated by expanding on plugging into the existing types and interfaces. These include
- the actual storage of files (and how to calculate the total storage size)
- how to ensure consistency between the indexer and the storage types (IE how to keep them in sync)
- how to load content from different root folders
- how to load content already available in storage when the application starts
- improve thread-safety
>>>>>>> updated design and further considerations sections
