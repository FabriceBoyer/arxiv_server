#!/bin/bash

# Common utils
wget_cmd="wget -q --show-progress --limit-rate=1000M"

source $(dirname -- "$0")/.env
out_dir=$DUMP_PATH
mkdir -p $out_dir

# Arxiv Metadata
# old_version : $wget_cmd https://storage.googleapis.com/arxiv-dataset/metadata-v5/arxiv-metadata-oai.json -O $out_dir/arxiv-metadata-oai.json
# up to date zip compressed version (requires login) : https://www.kaggle.com/datasets/Cornell-University/arxiv/download # ?datasetVersionNumber=103
# unzip it in $DUMP_PATH
