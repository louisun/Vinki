#!/bin/bash

qtc -dir ./templates wiki_template.qtpl

go build -o vinki main.go

./vinki
