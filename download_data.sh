#!/bin/bash

# Common utils
wget_cmd="wget -q --show-progress --limit-rate=10M"

source $(dirname -- "$0")/.env
out_dir=$DUMP_PATH
mkdir $out_dir

# Arxiv Metadata
$wget_cmd https://storage.googleapis.com/arxiv-dataset/metadata-v5/arxiv-metadata-oai.json -O $out_dir/arxiv-metadata-oai.json
# zip compressed version : https://www.kaggle.com/datasets/Cornell-University/arxiv/download # ?datasetVersionNumber=103
