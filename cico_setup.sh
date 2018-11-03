#!/bin/bash

# Output command before executing
set -x

# Exit on error
set -e

# Source environment variables of the jenkins slave
# that might interest this worker.
function load_jenkins_vars() {
  if [ -e "jenkins-env" ]; then
    cat jenkins-env \
      | grep -E "(DEVSHIFT_TAG_LEN|REGISTRY|REGISTRY_ORG|QUAY_USERNAME|QUAY_PASSWORD|JENKINS_URL|GIT_BRANCH|GIT_COMMIT|BUILD_NUMBER|ghprbSourceBranch|ghprbActualCommit|BUILD_URL|ghprbPullId)=" \
      | sed 's/^/export /g' \
      > ~/.jenkins-env
    source ~/.jenkins-env
  fi
}

function install_deps() {
  # We need to disable selinux for now, XXX
  /usr/sbin/setenforce 0 || :

  # Get all the deps in
  yum -y install \
    docker \
    make \
    git \
    curl

  service docker start

  echo 'CICO: Dependencies installed'
}

function cleanup_env {
  EXIT_CODE=$?
  echo "CICO: Cleanup environment: Tear down test environment"
  make integration-test-env-tear-down
  echo "CICO: Exiting with $EXIT_CODE"
}

function prepare() {
  # Let's test
  make docker-start
  make docker-deps
  make docker-generate
  make docker-build
  echo 'CICO: Preparation complete'
}

function run_tests_without_coverage() {
  make docker-test
  echo "CICO: ran tests without coverage"
}

function tag_push() {
  local tag
  local registry 

  tag=$1

  docker tag build-tool-detector-deploy $tag
  docker push $tag
}

function deploy() {

  if [ -z "$REGISTRY" ]
    then
        echo "\$var is empty, assigning default"
        REGISTRY="quay.io"
    else
        echo "\$Using custom registry $REGISTRY"
  fi

  if [ -z "$REGISTRY_ORG" ]
    then
        echo "\$var is empty, assigning default"
        REGISTRY_ORG="openshiftio"
    else
        echo "\$Using custom registry $REGISTRY"
  fi

  if [ -n "${QUAY_USERNAME}" -a -n "${QUAY_PASSWORD}" ]; then
    docker login -u ${QUAY_USERNAME} -p ${QUAY_PASSWORD} ${REGISTRY}
  else
    echo "Could not login, missing credentials for the registry"
  fi

  # Build build-tool-detector-depoy
  make docker-image-deploy

  TAG=$(echo $GIT_COMMIT | cut -c1-${DEVSHIFT_TAG_LEN})

  if [ "$TARGET" = "rhel" ]; then
    tag_push ${REGISTRY}/${REGISTRY_ORG}/rhel-fabric8-services-build-tool-detector:$TAG
    tag_push ${REGISTRY}/${REGISTRY_ORG}/rhel-fabric8-services-build-tool-detector:latest
  else
    tag_push ${REGISTRY}/${REGISTRY_ORG}/fabric8-services-build-tool-detector:$TAG
    tag_push ${REGISTRY}/${REGISTRY_ORG}/fabric8-services-build-tool-detector:latest
  fi

  echo 'CICO: Image pushed, ready to update deployed app'
}

function cico_setup() {
  load_jenkins_vars;
  install_deps;
  prepare;
}