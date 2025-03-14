#!/bin/bash
## It Should do all the setup automatically 
set -e 
echo "Step 1: Creating Docker Images..."
echo "Building: s3-reader"
cd reader-server && docker build -t s3-reader . && cd ..
echo "Building: s3-writer"
cd writer-server && docker build -t s3-writer . && cd ..
echo "Building: config-db-service"
cd db-config-service && docker build -t config-db-service . && cd ..
echo "Step 2: Building Client..."
cd client
if [ -f "Makefile" ]; then
    make
else
    echo "No Makefile found in client directory!"
    exit 1
fi
cp main.exe ../
cd ..
# cd proxy-client && docker build -t proxy-client . && cd ..
echo "Step 3: Running Containers..."
docker-compose up -d
echo "Client Built and Ready."
cd ./proxy-client && npm install && npm start
echo "DONOT Close this terminal, open a new one to work in now..."
