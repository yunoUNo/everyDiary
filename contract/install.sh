#!/bin/bash
set -ev 

# 1. 설치
docker exec cli peer chaincode install -n everyDiary -v 1.0 -p github.com/everyDiary

# 2. 인스톨
docker exec cli peer chaincode instantiate -n everyDiary -v 1.0 -C mychannel -c '{"Args":[]}' -P 'OR ("Org1MSP.member", "Org2MSP.member")'
sleep 5

# 3. Initial Diary 2 user
docker exec cli peer chaincode invoke -n everyDiary -C mychannel -c '{"Args":["set","yuno_1996-11-28","initial Diary"]}'
docker exec cli peer chaincode invoke -n everyDiary -C mychannel -c '{"Args":["set","arm_2000-08-08","initial Diary"]}'

# 4. query
docker exec cli peer chaincode query -n everyDiary -C mychannel -c '{"Args":["get","yuno_1996-11-28"]}'
docker exec cli peer chaincode query -n everyDiary -C mychannel -c '{"Args":["get","arm_2000-08-08"]}'
sleep 3
docker exec cli peer chaincode query -n everyDiary -C mychannel -c '{"Args":["history","yuno"]}'
sleep 3
docker exec cli peer chaincode query -n everyDiary -C mychannel -c '{"Args":["history","arm"]}'

# 5. 전체 검색
docker exec cli peer chaincode query -n everyDiary -C mychannel -c '{"Args":["checkUser"]}'

