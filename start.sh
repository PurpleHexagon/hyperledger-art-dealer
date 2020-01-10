rm -rf crypto-config/**
rm -rf orderer/genesis.block
rm -rf channels/*

docker-compose down
docker rmi --force $(docker images -q dev-peer*)

export FABRIC_CFG_PATH=./
cryptogen generate --config crypto-config.yaml

configtxgen -profile FabartOrdererGenesis -outputBlock ./orderer/genesis.block

configtxgen -profile testchannel -outputCreateChannelTx ./channels/testchannel.tx -channelID testchannel

configtxgen -profile testchannel -outputAnchorPeersUpdate ./channels/org1archor.tx -channelID testchannel -asOrg Org1MSP
configtxgen -profile testchannel -outputAnchorPeersUpdate ./channels/org2archor.tx -channelID testchannel -asOrg Org2MSP

docker-compose up
