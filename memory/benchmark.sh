#!/usr/bin/env bash

go test -bench=. -benchmem -benchtime=3s -count=3
