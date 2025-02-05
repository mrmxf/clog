#  clog> hello
# short> Display hello-world
# extra> and give some extra help

eval "$(clog Inc)"
fInfo "   _     _         _   _            _  _  _                _       _"
fInfo "  (_)   (_)       | | | |          (_)(_)(_)              | |     | |"
fInfo "   _______  _____ | | | |   ___     _  _  _   ___    ____ | |   __| |"
fInfo "  |  ___  || ___ || | | |  / _ \   | || || | / _ \  / ___)| |  / _  |"
fInfo "  | |   | || ____|| | | | | |_| |  | || || || |_| || |    | | ( (_| |"
fInfo "  |_|   |_||_____) \_) \_) \___/    \_____/  \___/ |_|     \_) \____|"
echo
fInfo "${cC}clog$cT scripts$cW v$(clog -v)"
echo
fInfo "${cC}clog$cT searches$cF \$(pwd)/clogrc$cT for files with a$cF .sh$cT extension"
fInfo "to be visible to clog, the script must start:"
echo
fInfo "    #$cC  clog> hello"
fInfo "    #$cC short> Display hello-world"
fInfo "    #$cC extra> and give some extra help"

