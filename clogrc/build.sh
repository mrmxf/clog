#  clog> build     #  build & inject metadata into clog
#  _                _        _
# ( )            _ (_ )     ( )
# | |_    _   _ (_) | |    _| |
# | '_`\ ( ) ( )| | | |  /'_` |
# | |_) )| (_) || | | | ( (_| |
# (_,__/'`\___/'(_)(___)`\__,_)
# ------------------------------------------------------------------------------
# load build config and script helpers
eval "$(clog Source project config)"    # configure project - local config
eval "$(clog Inc)"                      # shell embedded help (sh, zsh & bash)
help="core/sh/help-golang.sh"           # build embedded help
eval "$(clog Cat $help)"                # golang build helpers

clog Log -I "Build$cC $PROJECT $cT using clog's$cF $help"

clog Check pre-build && [ $? -gt 0 ] && exit 1
clog Check tools     && [ $? -gt 0 ] && exit 1
clog Check build     && [ $? -gt 0 ] && exit 1
# ------------------------------------------------------------------------------
go test
[ $? -gt 0 ] && exit 1
# ------------------------------------------------------------------------------

# ensure tmp dir exists
mkdir -p tmp

branch="$(clog git branch)"
hash="$(clog git hash head)"                                # use the head hash as the build hash
suffix="" && [[ "$branch" != "main" ]] && suffix="$branch"  # use the branch name as the suffix
app=clog                                                    # command you type to run the build
title="Command Line Of Go"                                  # title of the software
linkerPath="github.com/mrmxf/clog/semver.SemVerInfo"        # go tool objdump -S clog|grep /semver.SemVerInfo

fGoBuild tmp/$app-amd-lnx     linux   amd64 $hash "$suffix" $app "$title" "$linkerPath"
fGoBuild tmp/$app-amd-win.exe windows amd64 $hash "$suffix" $app "$title" "$linkerPath"
fGoBuild tmp/$app-amd-mac     darwin  amd64 $hash "$suffix" $app "$title" "$linkerPath"
fGoBuild tmp/$app-arm-lnx     linux   arm64 $hash "$suffix" $app "$title" "$linkerPath"
fGoBuild tmp/$app-arm-win.exe windows arm64 $hash "$suffix" $app "$title" "$linkerPath"
fGoBuild tmp/$app-arm-mac     darwin  arm64 $hash "$suffix" $app "$title" "$linkerPath"

clog Log -I "${cT}All built to the$cF tmp/$cT folder\n"