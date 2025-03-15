## S3-clone (single Node)
This is a simple system design for a S3-like object store. I have used couple of practices used in S3. Instead of buckets, i relate each bulk-upload to a ```URI``` and each file to a ```Hash``` in that URI namespace.
### Services:
1. ```writer-service```: This consists of multiple servers that write to their respective storages.
2. ```reader-service```: This has multiple servers that read from their respective storages.
3. ```allocation-policy```: This is based on simple 4-bit ```Node-offset``` in the Hash, which decides to which storages a file(blob) is written into. This is NOT a separate service in this implementation, instead i merged it into client-cli.
4. ```db-config-service```: This is the config-db service that manages configuration of each file-hash and uri.
5. ```client(cli)```: Supports read/write/status to the system.
6. ```client(web-based)```: Suppors serving resources when requested to a specific url.
### How-to-test?
##### requirements:
1. Docker Daemon should be running.
2. (To run bash script in Windows) Use Git Bash or WSL.
##### Steps:
1. run ```./bnr.sh``` after cloning this repo.
2. open new terminal and follow basic-commands(cli) sections to read/write.
##### Basic Commands (CLI)
1. ```Write```: ```./main.exe write <path_of_dir_or_file_to_upload>```.
2. ```Read```: ```./main.exe read <URI_of_upload> <HASH_of_file> <BASE_DIR_TO_Write_the_read_to>```.
3. ```Status```: ```./main.exe status <URI_of_upload>```
