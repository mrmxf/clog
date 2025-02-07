#   _            _                             _
#  | |_    ___  | |  _ __   ___   __ _   ___  | |  __ _   _ _    __ _
#  | ' \  / -_) | | | '_ \ |___| / _` | / _ \ | | / _` | | ' \  / _` |
#  |_||_| \___| |_| | .__/       \__, | \___/ |_| \__,_| |_||_| \__, |
#                   |_|          |___/                          |___/
# =============================================================================
#   _             _   _      _
#  | |__   _  _  (_) | |  __| |
#  | '_ \ | || | | | | | / _` |
#  |_.__/  \_,_| |_| |_| \__,_|
# -----------------------------------------------------------------------------
# golang build helpers to be included during build
if [ -z "$1" ];then
  # be quiet if $1 provided
  _m="Using"

  _m="$_m$cC golang$cT" && _a="$(which go 2>/dev/null)"
  [ -z "$_a" ] && _m="$_m$cE missing$cT"
  [ -n "$_a" ] && _m="$_m $(go version 2>/dev/null|grep -oE '[0-9]+\.[0-9]+\.[0-9]+')"

  _m="$_m,$cC docker$cT" && _a="$(which docker 2>/dev/null)"
  [ -z "$_a" ] && _m="$_m$cE missing$cT"
  [ -n "$_a" ] && _m="$_m $(docker system info 2>/dev/null|grep Name|cut -c 8-)"

  _m="$_m,$cC ko$cT" && _a="$(which ko 2>/dev/null)"
  [ -z "$_a" ] && _m="$_m$cE missing$cT"
  [ -n "$_a" ] && _m="$_m $(ko version 2>/dev/null)"

  _m="$_m,$cC aws-cli$cT" && _a="$(which aws 2>/dev/null)"
  [ -z "$_a" ] && _m="$_m$cE missing$cT"
  [ -n "$_a" ] && _m="$_m $(aws --version|grep -oE '[0-9]+\.[0-9]+\.[0-9]+'|head -1 2>/dev/null)"

  fInfo "$_m"
fi
# -----------------------------------------------------------------------------

# You can get the semver path with the following command:
#      go tool objdump -S tmp/opentsg-amd-lnx | grep /semver.SemVerInfo
# e.g. github.com/mrmxf/opentsg-node/semver.SemVerInfo

fGoBuild(){
  local gofile=$(printf "%-30s" "$1")
  local goos="$2"
  local goarch="$3"
  local commitHash="$4"
  local buildSuffix="$5"
  local buildAppName="$6"
  local buildAppTitle="$7"
  local linkerDataSemverPath="$8"
  
  #tidy params...
  # get ISO date
  printf -v buildDate '%(%Y-%m-%d)T' -1
  # remove spaces from App Title
  buildAppTitle=$(echo "$buildAppTitle" | tr ' ' '_')

    # set colors for printing logs
  local cos="$cLnx"
  [[ "$goos" == "darwin" ]] && cos=$cMac
  [[ "$goos" == "windows" ]] && cos=$cWin
  local car="$cArm"
  [[ "$goarch" == "amd64" ]] && car=$cAmd

  # determine build local OS & bCPU
  bCPU="${cAmd}amd$cX" && case $(uname -m) in arm*) bCPU="${cArm}arm$cX";; esac
  case "$(uname -s)" in
    Linux*)  bOSV=${cLnx}Linux$cX;;
    Darwin*)  bOSV=${cMac}Mac$cX;;
          *)  bOSV="untested:$(uname -s)";;
  esac

  # pretty format platform strings
  pad=$((17-${#goos}-${#goarch}))
  spaces="$(echo '----------'|head -c $pad)"
  printf -v tPlatform "for %s %s" "$cos$goos$cX/$car$goarch$cX" "$spaces"
  printf -v bPlatform "(built on %14s)" "$bOSV/$bCPU"

  # create linker data info:
  ldi="$commitHash|$buildDate|$buildSuffix|$buildAppName|$buildAppTitle"

  #create linker data string
  lds="-X $linkerDataSemverPath='$ldi'"

  # prepare build message
  buildMsg="$cos$gofile$cX $tPlatform $bPlatform"
  printf "${cI}>>   INFO$cT $buildMsg\r"
  
  GOOS="$goos" GOARCH="$goarch" go build -ldflags "$lds" -o $gofile
  if [ $? -gt 0 ]; then
    fError "$buildMsg ...  build failed"
    fError "Linker data string was:$cC -ldflags \"$lds\""
  else
    size="$(du --apparent-size --block-size=M $gofile)"
    fInfo "$buildMsg ... $size"
  fi
}

# -----------------------------------------------------------------------------

fDoHeading(){
  local dockerfile="$1"
  local loadOrPush="$2"
  local         os="$3"
  local       arch="$4"
  #check that arch exists in the MAKE string
  [[ "${MAKE#*"$arch"}" == "$MAKE" ]] && echo "$arch not in ($MAKE) - skipping" && return 0
  #do the build
  local platform="$os/$arch"
  local cArch="$cAmd"
  echo $arch | grep "arm" >/dev/null && cArch="$cArm"
  local cOs="$cLnx"
  echo $os | grep "darwin" >/dev/null && cOs="$cMac"
  echo $os | grep "indows" >/dev/null && cOs="$cWin"
  local t1=""   && [ -n "$5" ] && t1="$cX${cT}tag1=$cS$5 "
  local t2=""   && [ -n "$6" ] && t2="$cX${cT}tag2=$cC$6 "
  local t3=""   && [ -n "$7" ] && t3="$cX${cT}tag3=$cI$7 "
  local xtra="" && [ -n "$8" ] && t4="$cX${cT}tag4=$cW$8 "
  fInfo "Build $cOs$PROJECT$cT for $cArch$platform$cX$cT from=$cW$dockerfile$cT"
  fInfo "      tags: $t1 $t2 $t3 $xtra"
}

fMakeTags(){
  local slug=$1
  local arch="$(echo $2|grep -oE '.+[^0-9]{1,2}')"
  T1="$DOCKER_NS/$bBASE-$slug-${arch}:latest"
  T2="$DOCKER_NS/$bBASE-$slug-${arch}:$vCODE"
  #T3="$bBASE-$SLUG-${arch}:$vCODE"
}
# -----------------------------------------------------------------------------

fDockerBuild(){
  local dockerfile="-f $1"
  local LoadOrPush="--$2"
  local         os="$3"
  local       arch="$4"
  #check that arch exists in the MAKE string
  [[ "${MAKE#*"$arch"}" == "$MAKE" ]] && return 0
  local t1="" && [ -n "$5" ] && t1="--tag $5"
  local t2="" && [ -n "$6" ] && t2="--tag $6"
  local t3="" && [ -n "$7" ] && t3="--tag $7"
  docker buildx build $dockerfile . $LoadOrPush --platform "$os/$arch" $t1 $t2 $t3
  if [ $? -gt 0 ]; then
    printf "${cS}FAIL$cF $FF$cT build failed$cX\n"
    exit 1
  fi
}

# -----------------------------------------------------------------------------
fDockerLogin () {
  fInfo "Using Docker $cF$(docker system info 2>/dev/null| grep Name)"
  # ensure we have logged into docker (docker doesn't store state so this is idempotent)
  echo "$DOCKER_PAT" | docker login -u "$DOCKER_USR" --password-stdin
}

# -----------------------------------------------------------------------------
fShouldMake(){
  local test="$1"
  # return err if $1 is not in the MAKE env var
  [[ "${MAKE#*"$test"}" == "$MAKE" ]] && return 1
  return 0
}