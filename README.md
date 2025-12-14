# Direct Share

A tiny CLI tool for transferring files directly between two peers with a helpful progress bar, ETA, and file integrity verification.

## Installation

To install you can simply download the latest release from the "Releases" page, then copy it to whereever you want. 

Alternatively you can also build the tool yourself.

### Building

Simply run `go build` in the root of the project, and the executable will be built in the same directory.

## Usage

### 1. Receive a File (Listener)

Start the receiver on the machine that will accept the file. By default, it listens on port `:9000`.

```bash
# Listen on default port 9000
./direct-share listen

# Listen on a specific port
./direct-share listen -port :8080
```

### 2. Send a File (Sender)

Send a file to the listener's IP address and port.

```bash 
# Send a file to a specific IP and port
./direct-share send -addr 192.168.1.5:8080 -file ./large-video.mp4
```
