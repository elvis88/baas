#!/bin/bash

# This is just a little script that can be downloaded from the internet to
# install ft. It just does platform detection, downloads the installer
# and runs it.

set -u

# constant variable
BASS_ROOT="http://localhost:8080/api/v1"
APP_NAME="agent"
BINARY_NAME="agent"
BINARY_CONFIG="client.yaml"

# derived variable
BINARY_DIR="${HOME}/.baas/${APP_NAME}"
BINARY_CMD="${BINARY_DIR}/$BINARY_NAME"

say() {
    printf "${BINARY_CMD}: %s\n" "$1"
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
            # echo "Warning: Not forcing TLS v1.2, this is potentially less secure"
            curl --silent --show-error --fail --location "$1" --output "$2" $3
        else
            # echo "curl --proto '=https' --tlsv1.2 --silent --show-error --fail --location $1 --output $2"
            curl --proto '=https' --tlsv1.2 --silent --show-error --fail --location "$1" --output "$2" $3
        fi
    elif [ "$_dld" = wget ]; then
        if [[ $1 =~ https ]] || ! check_help_for wget --https-only --secure-protocol; then
            # echo "Warning: Not forcing TLS v1.2, this is potentially less secure"
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
agent 1.0.0 (2019-12-24)
The installer for agent

USAGE:
    agent.sh [FLAGS] COMMAND

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
    ensure ps aux | grep "${BINARY_CMD} -id ${AgentID}" | grep -v grep > /dev/null
    if [ $? -eq 0 ] 
    then
        err "already runing....."
    fi

    _ostype="$(uname -s)"
    _file=${BINARY_NAME}
    _url="${BASS_ROOT}/agent/${_ostype}/${_file}"
    
    ensure mkdir -p "${BINARY_DIR}"
    if [ ! -f ${BINARY_DIR}/$_file ]; then
        ensure downloader $_url ${BINARY_DIR}/$_file "--header Authorization:${BASS_Authorization}"
        ensure chmod +x ${BINARY_DIR}/$_file
    fi

    # TODO
    _file="${BINARY_CONFIG}"
    _url="${BASS_ROOT}/agent/${_ostype}/${_file}"
    ensure downloader $_url ${BINARY_DIR}/$_file "--header Authorization:${BASS_Authorization}"

    local _timestamp="$(date "+%Y%m%d%H%M%S")"
    ensure  cd ${BINARY_DIR}; ${BINARY_CMD} -id ${AgentID} > ${BINARY_DIR}/${_timestamp}_${AgentID}.log 2>&1 &
    say "started"
}

stop() {
    ensure ps aux | grep "${BINARY_CMD} -id ${AgentID}" | grep -v grep | awk '{print "kill -9 " $2}' | sh
    say "stoped"
}

restart() {
    stop ${AgentID}
    sleep 1s
    start ${AgentID}
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
