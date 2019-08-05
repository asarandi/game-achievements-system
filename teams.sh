#!/usr/bin/env bash

for i in `seq 1 3`; do curl -s -X POST http://0.0.0.0:4242/teams/42/members/"$i" | jq; done
for i in `seq 11 13`; do curl -s -X POST http://0.0.0.0:4242/teams/84/members/"$i" | jq; done
curl -s -X POST http://0.0.0.0:4242/games | jq
curl -s -X POST http://0.0.0.0:4242/games/1/teams/42 | jq
curl -s -X POST http://0.0.0.0:4242/games/1/teams/84 | jq

for i in `seq 21 24`; do curl -s -X POST http://0.0.0.0:4242/teams/33/members/"$i" | jq; done
for i in `seq 31 34`; do curl -s -X POST http://0.0.0.0:4242/teams/66/members/"$i" | jq; done
curl -s -X POST http://0.0.0.0:4242/games | jq
curl -s -X POST http://0.0.0.0:4242/games/2/teams/33 | jq
curl -s -X POST http://0.0.0.0:4242/games/2/teams/66 | jq


for i in `seq 41 45`; do curl -s -X POST http://0.0.0.0:4242/teams/55/members/"$i" | jq; done
for i in `seq 51 55`; do curl -s -X POST http://0.0.0.0:4242/teams/77/members/"$i" | jq; done
curl -s -X POST http://0.0.0.0:4242/games | jq
curl -s -X POST http://0.0.0.0:4242/games/3/teams/55 | jq
curl -s -X POST http://0.0.0.0:4242/games/3/teams/77 | jq
