#   _            _            _
#  | |_    ___  | |  _ __    | |_    _  _   __ _   ___
#  | ' \  / -_) | | | '_ \   | ' \  | || | / _` | / _ \
#  |_||_| \___| |_| | .__/   |_||_|  \_,_| \__, | \___/
#                   |_|                    |___/
#   _             _   _      _
#  | |__   _  _  (_) | |  __| |
#  | '_ \ | || | | | | | / _` |
#  |_.__/  \_,_| |_| |_| \__,_|
# -----------------------------------------------------------------------------
# build a hugo site `into public/`
fHugoBuild() {
  local opts
  opts="$1"
  # set the default opts for a production build if none given
  [ -z "$opts" ] && opts="--gc --minify --forceSyncStatic --ignoreCache --logLevel info"

  # tidy up before starting
    clog Log -I "purging old builds:  $ ${cC}rm ${cW}-rf$cF public/*$cX"
  rm -rf public/*
  [ $? -gt 0 ] && fWarn "purging failed - continuing anyway"

  clog Log -I "building hugo site with opts$cC $opts"
  hugo build $opts
  [ $? -gt 0 ] && fError "hugo build failed" && exit 1
  fOk   "static website (${cW}$(clog git tag ref)$cT) in$cF public/$cX"
}

#      _              _
#   __| |  ___   __  | |__  ___   _ _
#  / _` | / _ \ / _| | / / / -_) | '_|
#  \__,_| \___/ \__| |_\_\ \___| |_|
# -----------------------------------------------------------------------------
# build a dockerfile with Hugo artifacts
fHugoDocker() {
  local opts dockerfile t1 t2 t3 t4
  opts="$1"
  dockerfile="$2"
  platform="$3"
  t1="$4"
  t2="$5"
  t3="$6"
  t4="$7"

  # override empty parameters with production defaults
  [ -z "$opts"] && opts="--no-cache --push --progress auto"
  if [ -z "$dockerfile" ]; then
    dockerfile="."
  else
    dockerfile="--file=$dockerfile"
  fi
  [ -z "$platform" ] && platform="linux/amd64"
  [ -n "$t1" ] && t1="--tag=$t1"
  [ -n "$t2" ] && t2="--tag=$t2"
  [ -n "$t3" ] && t3="--tag=$t3"
  [ -n "$t4" ] && t4="--tag=$t4"
  clog Log -I "building hugo docker image"
  clog Log -I "${cC}docker buildx build$cW $opts$cS \"--platform=$platform\"$cT $t1 $t2 $t3 $t4 \"$dockerfile\" ."
              docker buildx build    $opts      "--platform=$platform" \
              $t1   $t2 $t3 $t4 \
              "$dockerfile" .
}