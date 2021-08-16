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

## Repo structure and design considerations
The core logic of this package is subdivided into modules available inside the /pkg directory.

### Storage module
The Storage module exposes an interface for adding, removing and retrieving raw content from the underlying storage such as the local file system. This module is mainly concerned with how raw data is persisted and it does not deal with any metadata with the exception of the content path. The current implementation used here bypasses the actual storage of content and should only be seen as a noop version.

### Indexer module
The Indexer module is responsible for managing the provider storage state. It’s the layer between the raw storage and IO operations and provides functionalities for efficiently fetching data and metadata as well as managing the provider storage capacity.
The main functionalities of the Indexer interface are:
- the retrieval of content given its CID from the underlying storage layer
- the storage of blocks in the underlying storage layer

Its `MemIndexer` implementation uses an in memory key-value map for associating the CID of individual blocks of content to metadata specific to the content itself.
This metadata includes:
- The block CID
- Block Address (path to the actual content)
- Hits count (the number of times the block has been queried)
- Block size
- Block creation timestamp
- Last block retrieval timestamp

The internal usage of a map has been chosen for providing fast access to content given its identifier. Using CIDs as keys ensure that only unique blocks of content will be stored at any given time.

**Storage Size Management**
The `MemIndexer` storage size is naively measured in numbers of blocks being stored. While this option provides a simple mechanism to test against, it should be replaced by using the actual file size of the various blocks of content being committed to storage.
As new blocks of content are submitted by clients, the `MemIndexer` will check the current storage size against its maximum capacity and when the two values are equal the block with the lowest number of hits will be removed from both the key-value map and the underlying data storage.
The `Evictor` interface used by the `MemIndexer` struct allows for more sophisticated algorithms to be used for removing least frequently accessed items.

### Provider
The Provider module mediates access with the content Indexer and it is the outer access layer that is being made available through an HTTP server. This layer could be potentially used for providing access to a UI, integrating monitoring or auditing middlewares and more.

## Further considerations
