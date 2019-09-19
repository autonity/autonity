#!/bin/sh

TESTDIR=test_race_files

if [ ! -d $TESTDIR ]; then
  mkdir $TESTDIR
fi

echo "Running TestTendermintSuccess..."
GORACE="history_size=7" go test ./... -run=TestTendermintSuccess                -v -race -timeout=60m >  $TESTDIR/success-race.txt

echo "Running TestTendermintOneMalicious..."
GORACE="history_size=7" go test ./... -run=TestTendermintOneMalicious           -v -race -timeout=60m >  $TESTDIR/one-malicious-race.txt

echo "Running TestTendermintSlowConnections..."
GORACE="history_size=7" go test ./... -run=TestTendermintSlowConnections        -v -race -timeout=60m >  $TESTDIR/slow-connection-race.txt

echo "Running TestTendermintLongRun..."
GORACE="history_size=7" go test ./... -run=TestTendermintLongRun                -v -race -timeout=60m >  $TESTDIR/long-run-race.txt

echo "Running TestTendermintStopUpToFNodes..."
GORACE="history_size=7" go test ./... -run=TestTendermintStopUpToFNodes         -v -race -timeout=60m >  $TESTDIR/stop-up-to-f-nodes-race.txt

echo "Running TestTendermintStartStopSingleNode..."
GORACE="history_size=7" go test ./... -run=TestTendermintStartStopSingleNode    -v -race -timeout=60m >  $TESTDIR/start-stop-single-node-race.txt

echo "Running TestTendermintStartStopFNodes..."
GORACE="history_size=7" go test ./... -run=TestTendermintStartStopFNodes        -v -race -timeout=60m >  $TESTDIR/start-stop-f-nodes-race.txt

echo "Running TestTendermintStartStopFPlusOneNodes..."
GORACE="history_size=7" go test ./... -run=TestTendermintStartStopFPlusOneNodes -v -race -timeout=60m >  $TESTDIR/start-stop-f-plus-one-nodes-race.txt

echo "Running TestTendermintStartStopFPlusTwoNodes..."
GORACE="history_size=7" go test ./... -run=TestTendermintStartStopFPlusTwoNodes -v -race -timeout=60m >  $TESTDIR/start-stop-f-plus-two-nodes-race.txt

echo "Running TestTendermintStartStopAllNodes..."
GORACE="history_size=7" go test ./... -run=TestTendermintStartStopAllNodes      -v -race -timeout=60m >  $TESTDIR/start-stop-all-nodes-race.txt