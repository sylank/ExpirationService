#!/bin/bash

CONFIG_DIR="./config"

VERSION=$(cat infrastructure/VERSION)
APP_NAME=$(cat infrastructure/APP_NAME)

echo "Download credentials..."
#rm $CONFIG_DIR/calendars.json
#rm $CONFIG_DIR/credentials.json
#rm $CONFIG_DIR/secrets.json
#rm $CONFIG_DIR/token.json

#aws s3 cp s3://artifactory.levendulabalatonmaria.info/credentials/$APP_NAME/calendars.json $CONFIG_DIR/calendars.json
#aws s3 cp s3://artifactory.levendulabalatonmaria.info/credentials/$APP_NAME/credentials.json $CONFIG_DIR/credentials.json
#aws s3 cp s3://artifactory.levendulabalatonmaria.info/credentials/$APP_NAME/secrets.json $CONFIG_DIR/secrets.json
#aws s3 cp s3://artifactory.levendulabalatonmaria.info/credentials/$APP_NAME/token.json $CONFIG_DIR/token.json

ls -la $CONFIG_DIR