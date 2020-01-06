#!/bin/bash

# This is just a little script that can be downloaded from the internet to
# install ft. It just does platform detection, downloads the installer
# and runs it.

set -u

# constant variable
BASS_ROOT="http://localhost:8080/api/v1"
RELEASE_ROOT="https://github.com/fractalplatform"
RELEASE_NAME="fractal"
RELEASE_VERSION="1.0.0"
RELEASE_EXT="tar.gz"
BINARY_NAME="ft" 
RELEASE_CONFIG="genesis.json"
DEPLOY_CONFIG="config.yaml"

# template variable
DEPLOY_USER="testuser" # node user
DEPLOY_NAME="ft"  # node name
APP_NAME="ft"

# derived variable
BINARY_RELEASE_CONFIG="${RELEASE_CONFIG}"
BINARY_DEPLOY_CONFIG="${DEPLOY_CONFIG}"
RELEASE_CONFIG_ROOT="${BASS_ROOT}/file/${APP_NAME}/application"
DEPLOY_CONFIG_ROOT="${BASS_ROOT}/file/${DEPLOY_NAME}/deployment"
DEPLOY_DIR="${HOME}/.baas/${BINARY_NAME}/${DEPLOY_NAME}"
BINARY_ARG="-g ${DEPLOY_DIR}/${BINARY_RELEASE_CONFIG} -c ${DEPLOY_DIR}/${BINARY_DEPLOY_CONFIG}"
DEPLOY_CMD="${DEPLOY_DIR}/$BINARY_NAME ${BINARY_ARG}"

say() {
    printf "${DEPLOY_USER} ${DEPLOY_NAME}: %s\n" "$1"
}

err() {
    say "$1" >&2
    exit 1
}

need_cmd() {
    if ! check_cmd "$1"; then
        err "need '$1' (command not found)"
    fi
}

check_cmd() {
    command -v "$1" > /dev/null 2>&1
}

# Run a command that should never fail. If the command fails execution
# will immediately terminate with an error showing the failing
# command.
ensure() {
    if ! "$@"; then err "command failed: $*"; fi
}

# This is just for indicating that commands' results are being
# intentionally ignored. Usually, because it's being executed
# as part of error handling.
ignore() {
    "$@"
}

# This wraps curl or wget. Try curl first, if not installed,
# use wget instead.
downloader() {
    local _dld
    if check_cmd curl; then
        _dld=curl
    elif check_cmd wget; then
        _dld=wget
    else
        _dld='curl or wget' # to be used in error message of need_cmd
    fi

    if [ "$1" = --check ]; then
        need_cmd "$_dld"
    elif [ "$_dld" = curl ]; then
        if ! [[ $1 =~ https ]] || ! check_help_for curl --proto --tlsv1.2; then
            echo "Warning: Not forcing TLS v1.2, this is potentially less secure"
            curl --silent --show-error --fail --location "$1" --output "$2" $3
        else
            echo "curl --proto '=https' --tlsv1.2 --silent --show-error --fail --location $1 --output $2"
            curl --proto '=https' --tlsv1.2 --silent --show-error --fail --location "$1" --output "$2" $3
        fi
    elif [ "$_dld" = wget ]; then
        if [[ $1 =~ https ]] || ! check_help_for wget --https-only --secure-protocol; then
            echo "Warning: Not forcing TLS v1.2, this is potentially less secure"
            wget "$1" -O "$2" $3
        else
            wget --https-only --secure-protocol=TLSv1_2 "$1" -O "$2" $3
        fi
    else
        err "Unknown downloader"   # should not reach here
    fi
}

check_help_for() {
    local _cmd
    local _arg
    local _ok
    _cmd="$1"
    _ok="y"
    shift

    for _arg in "$@"; do
        if ! "$_cmd" --help | grep -q -- "$_arg"; then
            _ok="n"
        fi
    done

    test "$_ok" = "y"
}

usage() {
    cat 1>&2 <<EOF
ft 1.0.0 (2019-12-24)
The installer for ft

USAGE:
    ft [FLAGS] COMMAND

FLAGS:
    -h, --help              Prints help information
    -v, --version           Prints version information

COMMAND:
    start                install & run
    stop                 kill 
    restart              kill & run
EOF
}

start() {
    ensure ps aux | grep "${DEPLOY_CMD}" | grep -v grep > /dev/null
    if [ $? -eq 0 ] 
    then
        err "already runing....."
    fi

    local _ostype _cputype _arch _file _url
    _ostype="$(uname -s)"
    _cputype="$(uname -m)"
    _arch="${_ostype}_${_cputype}"
    _file="${RELEASE_NAME}_${RELEASE_VERSION}_${_arch}.${RELEASE_EXT}"
    _url="${RELEASE_ROOT}/${RELEASE_NAME}/releases/download/v${RELEASE_VERSION}/${_file}"
    
    ensure mkdir -p "${DEPLOY_DIR}"
    if [ ! -f ${DEPLOY_DIR}/$_file ]; then
        ensure downloader $_url ${DEPLOY_DIR}/$_file
        ensure tar -xzf ${DEPLOY_DIR}/$_file -C ${DEPLOY_DIR}
    fi

    # TODO
    _file="${BINARY_RELEASE_CONFIG}"
    _url="${RELEASE_CONFIG_ROOT}/${_file}"
    ensure downloader $_url ${DEPLOY_DIR}/$_file "--header \"Authorization:${BASS_Authorization}\""

    _file="${BINARY_DEPLOY_CONFIG}"
    _url="${DEPLOY_CONFIG_ROOT}/${_file}"
    ensure downloader $_url ${DEPLOY_DIR}/$_file "--header \"Authorization:${BASS_Authorization}\""

    local _timestamp="$(date "+%Y%m%d%H%M%S")"
    ensure  ${DEPLOY_CMD} > ${DEPLOY_DIR}/${_timestamp}_${DEPLOY_NAME}.log 2>&1 &
    say "started"
}

stop() {
    ensure ps aux | grep "${DEPLOY_CMD}" | grep -v grep | awk '{print "kill -9 " $2}' | sh
    say "stoped"
}

restart() {
    stop
    sleep 1s
    start
}

main() {
    downloader --check
    need_cmd uname
    need_cmd mkdir
    need_cmd rm
    need_cmd chmod
    need_cmd tar
    need_cmd awk

    for arg in "$@"; do
        case "$arg" in
            -h|--help)
                usage
                exit 0
                ;;
            start)
                start
                exit 0
                ;;
            stop)
                stop
                exit 0
                ;;
            restart)
                restart
                exit 0
                ;;
            *)
                err "Unknown Arg: $arg";;
        esac
    done
}

main "$@" || exit 1
