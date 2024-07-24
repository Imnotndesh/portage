# Portage
A simple file sending utility written in golang

## Usage
- On Linux : ```./portage [OPTIONS] ```
- On Windows : ```./portage.exe [OPTIONS]```

## Valid Options
| Flag | Usage                        | Description                                                 |
|------|------------------------------|-------------------------------------------------------------|
| -s   | -s [REMOTE_IP] [FILES...]    | Sends passed FILES to the server running on the REMOTE_IP   |
| -r   | -r                           | Starts the server at the default port                       |
- The flags above are not case-sensitive

## In the works
- [x] Support for multiple files
- [ ] Progress tracker
- [ ] Notification when transfer complete on both sender and receiver
- [ ] Speed meter
- [ ] Concurrency for multiple file transfer connections