#!/bin/bash

source "$(dirname $0)/../../scripts/base.sh"
info "[$0]"
(
    set -e
    cd $ROOT_DIR

    PROMOTE="--no-promote"
    if [ "$1" == "--promote" ] ; then
        PROMOTE="--promote"
    fi

    GIT_VERSION=$(git rev-parse --short HEAD)-$(git describe --all --dirty --broken)
    echo GIT_VERSION $GIT_VERSION
    echo ""


    info "Collect source files"
    STAGE_DIR="tmp/deployments/gae"
    rm -rf $STAGE_DIR
    mkdir -p $STAGE_DIR
    cp -r deployments/gae/. $STAGE_DIR/
    cp -r internal $STAGE_DIR/
    cp -r web/static $STAGE_DIR/
    cp -r web/template $STAGE_DIR/
    cp -r go.mod $STAGE_DIR/
    cp -r go.sum $STAGE_DIR/
    echo $GIT_VERSION > $STAGE_DIR/version.txt
    echo "done"
    echo ""

    info "Deploy with $PROMOTE"
    (
        cd $STAGE_DIR
        echo "gcloud --project=$GC_PROJECT app deploy --quiet $PROMOTE"
        gcloud --project=$GC_PROJECT app deploy --quiet $PROMOTE
    )
    echo "done"
    echo ""

    info "Done, go to https://console.cloud.google.com/ to try, promote the deployed version and delete old one."
)
if [ $? == 0 ] ; then
    info "[$0 OK]"
    exit 0
else
    error "[$0 FAILED]"
    exit 1
fi