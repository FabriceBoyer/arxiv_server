Arxiv server
============

## Introduction

Minimal Arxiv metadata API server using local dumps

## Prerequisites

Use a recent Linux distribution

Install docker https://docs.docker.com/engine/install/

## Configuration

Copy .env.example and rename it to .env
Modify values as required to sore dump files in appropriate location

## Download data

Call ./download_dump.sh script to download latest dump

## Start server

Call ./start.sh to start the docker compose stack

Navigate to http://localhost:9097
