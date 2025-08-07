# clog version retrieved from CLOG_BIN
#                  _           _    _             _
#  _              ( )_        (_ ) (_ )          (_ )
# (_)  ___    ___ | ,_)   _ _  | |  | |      ___  | |    _      __
# | |/' _ `\/',__)| |   /'_` ) | |  | |    /'___) | |  /'_`\  /'_ `\
# | || ( ) |\__, \| |_ ( (_| | | |  | |   ( (___  | | ( (_) )( (_) |
# (_)(_) (_)(____/`\__)`\__,_)(___)(___)  `\____)(___)`\___/'`\__  |
#                                                            ( )_) |
#                                                             \___/'
#      curl https://mrmxf.com/get/clog | bash    # linux
#      curl https://mrmxf.com/get/clog | zsh     # mac
#
#     works in gitlab, github, gitpod, linux, mac
#
#   Command  ; Error     ; Info      ; File      ; Header      ; Success   ; Text      ; Url       ; Warning   ; eXit
cC="\e[34m";cE="\e[91m";cI="\e[93m";cF="\e[33m";cH="\e[36;1m";cS="\e[32m";cT="\e[30m";cU="\e[36m";cW="\e[35m";cX="\e[0m";

# set cOS: linux|alpine|mac|windows|gitpod, cPU: arm|amd, cErr: <some error>
fCompatibility() {
  local cErr=""
  # detect OS
  case "$(uname -s)" in
    Linux*|CYGWIN*|MINGW*|MSYS_NT*)     cOS="linux" ;;
    Darwin*)                            cOS="mac"   ;;
    win*)                               cOS="windows" ;;
    *)          cErr="Unsupported OS ($(uname -s))" ;;
  esac
  [ -n "${GITPOD_GIT_USER_NAME+x}" ] && cOS="gitpod"
  [ -n "$(uname -v|grep lpine)" ]    && cOS="alpine"
  # detect architecture
  case "$(uname -m)" in
      x86_64*|amd*)            cPU="amd"    ;;
    # i386*|i486*|i586*|i686*) cPU="i32"    ;;
      arm64|aarch64*)          cPU="arm"    ;;
      *) cErr="$cERR Unsupported ARCH ($(uname -m))" ;;
  esac

  [ -n "$cErr" ] && echo "$cErr" && exit 1
  export cOS
  export cPU
}
################################################################################

URL=https://mrmxf.com/get/CLOG_BIN
TMP=/tmp/clog
CPU="amd"
DST=/usr/local/bin
SUDO="sudo"

# figure out the combination of OS & CPU we are running on
fCompatibility

# disable sudo if we're in alpine or CICD
gotSudo="$(sudo >/dev/null 2>&1)"
[ $? -gt 1 ] && SUDO="" # $?==127 if sudo missing, $?==1 if it exists

# if we're on alpine replace gcc runtime with default musl runtime
[[ "$cOS" == "alpine" ]] && mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

case "$cOS" in
  linux|alpine) theOS="lnx" ;;
    mac)        theOS="mac" ;;
 gitpod)        theOS="lnx" ; DST=/workspace/go/bin ;;
      *) echo "clog doesn't run on $cOS" ; exit 1 ;;
esac

#echo $cOS $cPU $theOS $CPU $DST

# ensure the destination folder exists
$SUDO mkdir -p $DST

F="$URL/clog-$cPU-$theOS"
printf "%s ${cI}INF${cX} fetch $cW$cOS$cT clog ($cW$cPU$cT) ↓↓ $cF$F$cT →$cF$TMP\n" "$(date +"%Y-%m-%d %H:%M:%S")"
curl --location --silent "$F" --output "$TMP"

printf "%s ${cI}INF${cX} install →$cF$DST/clog\n" "$(date +"%Y-%m-%d %H:%M:%S")"
$SUDO mv "$TMP" "$DST/clog"
$SUDO chmod +x  "$DST/clog"
clog Log -I "$(clog --version)"
