#!/usr/bin/env bash

# --- begin runfiles.bash initialization v3 ---
# Copy-pasted from the Bazel Bash runfiles library v3.
set -uo pipefail; set +e; f=bazel_tools/tools/bash/runfiles/runfiles.bash
# shellcheck source=/dev/null.
source "${RUNFILES_DIR:-/dev/null}/$f" 2>/dev/null || \
source "$(grep -sm1 "^$f " "${RUNFILES_MANIFEST_FILE:-/dev/null}" | cut -f2- -d' ')" 2>/dev/null || \
source "$0.runfiles/$f" 2>/dev/null || \
source "$(grep -sm1 "^$f " "$0.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null || \
source "$(grep -sm1 "^$f " "$0.exe.runfiles_manifest" | cut -f2- -d' ')" 2>/dev/null || \
{ echo>&2 "ERROR: cannot find $f; this script must be run with 'bazel run'"; exit 1; }; f=; set -e
# --- end runfiles.bash initialization v3 ---

docker buildx build - --load --tag=pachyderm/ci-runner:"$(cat "$(rlocation _main/etc/ci-image/version)")" < "$(rlocation _main/etc/ci-image/Dockerfile)" --tag=pachyderm/ci-runner:latest
exec docker run -it -v "$BUILD_WORKSPACE_DIRECTORY":/home/circleci/project:ro \
     pachyderm/ci-runner:latest
