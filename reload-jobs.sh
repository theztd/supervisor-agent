#!/bin/bash -l

echo "Restarting postgresql database..."
pushd ./devel
docker-compose restart postgresql
popd
echo ""

sleep 5

echo "Reloading supervisor jobs..."
curl -kL 'http://localhost:9001/index.html?action=restartall'
echo ""

sleep 5
echo "Waiting for connections..."
