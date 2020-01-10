docker exec -it cli.org1 bash -c "peer channel create -c testchannel -f ./channels/testchannel.tx -o orderer.fabart.example.com:7050"

docker exec -it cli.org1 bash -c "peer channel join -b testchannel.block"
docker exec -it cli.org2 bash -c "peer channel join -b testchannel.block"

docker exec cli.org1 bash -c 'peer chaincode install -p fabart -n fabart -v 0'
docker exec cli.org2 bash -c 'peer chaincode install -p fabart -n fabart -v 0'

docker exec cli.org1 bash -c "peer chaincode instantiate -C testchannel -n fabart -v 0 -c '{\"Args\":[]}'"
docker exec cli.org2 bash -c "peer chaincode instantiate -C testchannel -n fabart -v 0 -c '{\"Args\":[]}'"


docker exec cli.org1 bash -c "peer chaincode invoke -C testchannel -n fabart -c '{\"Args\":[\"register\", \"1\", \"2000\", \"bar\", \"foo\"]}'"
sleep 3

docker exec cli.org2 bash -c "peer chaincode query -C testchannel -n fabart -c '{\"Args\":[\"queryDetails\", \"1\"]}'"
sleep 1

docker exec cli.org1 bash -c "peer chaincode invoke -C testchannel -n fabart -c '{\"Args\":[\"transfer\", \"1\", \"caaaaa\"]}'"
sleep 1
