#!/usr/bin/env bash

# check TARGET_BRANCH existance
# if [ "$TARGET_BRANCH" = "" ]; then
# 	echo "you must set TARGET_BRANCH"
# 	exit 0
# fi
# 
# # checkout branch
# echo "-------- fetch source ---------"
HERE=/home/isucon/isubata
# git fetch origin "$TARGET_BRANCH":"$TARGET_BRANCH"
# git co "$TARGET_BRANCH"


# compile
echo "-------- compile source ---------"
cd webapp/go
make

# restart services
echo "-------- restart services ---------"
sudo systemctl restart nginx.service
sudo systemctl restart isubata.ruby
# sudo systemctl restart isubata.golang

# sleep 2
echo "-------- 2sec sleeping ---------"
sleep 2

# run benchmark
echo "-------- run benchmark ---------"
cd "$HERE"/bench
./bin/bench -remotes=127.0.0.1 -output result.json; cat result.json | jq

# back project root
cd "$HERE"
git co master
