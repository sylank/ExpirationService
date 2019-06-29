#!/bin/bash

VERSION=$(cat infrastructure/VERSION)
APP_NAME=$(cat infrastructure/APP_NAME)

ls -la

mv main $APP_NAME-$VERSION

zip -r $APP_NAME-$VERSION.zip $APP_NAME-$VERSION config

ls -la

#Delete if exists
aws s3 rm s3://artifactory.levendulabalatonmaria.info/$APP_NAME/$VERSION --recursive

#Copy the new content
aws s3 cp $APP_NAME-$VERSION.zip s3://artifactory.levendulabalatonmaria.info/$APP_NAME/$VERSION/$APP_NAME-$VERSION.zip

#Generate hash
hash=$(cat $APP_NAME-$VERSION | base64 | sha256sum)
echo $hash > $APP_NAME-$VERSION.hash

aws s3 cp $APP_NAME-$VERSION.hash s3://artifactory.levendulabalatonmaria.info/$APP_NAME/$VERSION/$APP_NAME-$VERSION.hash