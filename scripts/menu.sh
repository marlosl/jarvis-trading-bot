#!/bin/bash

# Reset
Color_Off='\033[0m'       # Text Reset

# Regular Colors
Black='\033[0;30m'        # Black
Red='\033[0;31m'          # Red
Green='\033[0;32m'        # Green
Yellow='\033[0;33m'       # Yellow
Blue='\033[0;34m'         # Blue
Purple='\033[0;35m'       # Purple
Cyan='\033[0;36m'         # Cyan
White='\033[0;37m'        # White

# Bold
BBlack='\033[1;30m'       # Black
BRed='\033[1;31m'         # Red
BGreen='\033[1;32m'       # Green
BYellow='\033[1;33m'      # Yellow
BBlue='\033[1;34m'        # Blue
BPurple='\033[1;35m'      # Purple
BCyan='\033[1;36m'        # Cyan
BWhite='\033[1;37m'       # White

# Underline
UBlack='\033[4;30m'       # Black
URed='\033[4;31m'         # Red
UGreen='\033[4;32m'       # Green
UYellow='\033[4;33m'      # Yellow
UBlue='\033[4;34m'        # Blue
UPurple='\033[4;35m'      # Purple
UCyan='\033[4;36m'        # Cyan
UWhite='\033[4;37m'       # White

# Background
On_Black='\033[40m'       # Black
On_Red='\033[41m'         # Red
On_Green='\033[42m'       # Green
On_Yellow='\033[43m'      # Yellow
On_Blue='\033[44m'        # Blue
On_Purple='\033[45m'      # Purple
On_Cyan='\033[46m'        # Cyan
On_White='\033[47m'       # White

# High Intensity
IBlack='\033[0;90m'       # Black
IRed='\033[0;91m'         # Red
IGreen='\033[0;92m'       # Green
IYellow='\033[0;93m'      # Yellow
IBlue='\033[0;94m'        # Blue
IPurple='\033[0;95m'      # Purple
ICyan='\033[0;96m'        # Cyan
IWhite='\033[0;97m'       # White

# Bold High Intensity
BIBlack='\033[1;90m'      # Black
BIRed='\033[1;91m'        # Red
BIGreen='\033[1;92m'      # Green
BIYellow='\033[1;93m'     # Yellow
BIBlue='\033[1;94m'       # Blue
BIPurple='\033[1;95m'     # Purple
BICyan='\033[1;96m'       # Cyan
BIWhite='\033[1;97m'      # White

# High Intensity backgrounds
On_IBlack='\033[0;100m'   # Black
On_IRed='\033[0;101m'     # Red
On_IGreen='\033[0;102m'   # Green
On_IYellow='\033[0;103m'  # Yellow
On_IBlue='\033[0;104m'    # Blue
On_IPurple='\033[0;105m'  # Purple
On_ICyan='\033[0;106m'    # Cyan
On_IWhite='\033[0;107m'   # White

SERVICES_NAME[0]=""
SERVICES_NAME[1]="Build Bot"
SERVICES_NAME[2]="Postgres Server"
SERVICES_NAME[3]="PG Admin"
SERVICES_NAME[4]="Bot Server"
SERVICES_NAME[5]="Mock WebSocket Server"
SERVICES_NAME[6]="Mock Sender"
SERVICES_NAME[7]="Elastic ELK"

SERVICES_FILE[0]=""
SERVICES_FILE[1]="build-bot.yaml"
SERVICES_FILE[2]="db-postgres.yaml"
SERVICES_FILE[3]="db-pgadmin.yaml"
SERVICES_FILE[4]="bot-server.yaml"
SERVICES_FILE[5]="mock-websocket-server.yaml"
SERVICES_FILE[6]="mock-sender.yaml"
SERVICES_FILE[7]="elk.yaml"

SERVICES_PS[0]=""
SERVICES_PS[1]="build-bot"
SERVICES_PS[2]="postgres"
SERVICES_PS[3]="pgadmin"
SERVICES_PS[4]="bot-trading-server"
SERVICES_PS[5]="mock-websocket-server"
SERVICES_PS[6]="mock-sender"
SERVICES_PS[7]="kibana"

SERVICES_STATUS[0]=0
SERVICES_STATUS[1]=0
SERVICES_STATUS[2]=0
SERVICES_STATUS[3]=0
SERVICES_STATUS[4]=0
SERVICES_STATUS[5]=0
SERVICES_STATUS[6]=0
SERVICES_STATUS[7]=0

DOCKER_VOLUMES[0]="jarvis-postgres"
DOCKER_VOLUMES[1]="go-bot-packages"

DOCKER_NETWORK="jarvis-network"

MAX_SERVICES=7
MAX_DOCKER_VOLUMES=1

SCRIPT_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"
DOCKER_BASE_DIRECTORY="$SCRIPT_DIR/docker"

function createDockerNetwork() {
    ret=`docker network ls | grep $DOCKER_NETWORK`
    if [[ -z $ret ]]; then
        docker network create $DOCKER_NETWORK
    fi 
}

function createDockerVolumes() {
    for i in $(seq 0 $MAX_DOCKER_VOLUMES); do
        volume=${DOCKER_VOLUMES[$i]}
        ret=`docker volume ls | grep $volume`
        if [[ -z $ret ]]; then
            docker volume create $volume
        fi     
    done
}

function copyConfigFilesForBuild() {
    sudo rm -rf $SCRIPT_DIR/docker/build
    mkdir $SCRIPT_DIR/docker/build    
    cp $SCRIPT_DIR/../.env $SCRIPT_DIR/docker/build/.env
    cp $SCRIPT_DIR/../.parameters $SCRIPT_DIR/docker/build/.parameters
    cp $SCRIPT_DIR/docker/Dockerfile.bot $SCRIPT_DIR/docker/build/Dockerfile.bot
    cp $SCRIPT_DIR/docker/filebeat-8.0.0-linux-arm64.tar.gz $SCRIPT_DIR/docker/build/filebeat-8.0.0-linux-arm64.tar.gz
    cp $SCRIPT_DIR/docker/filebeat.yml $SCRIPT_DIR/docker/build/filebeat.yml
    cp $SCRIPT_DIR/docker/entrypoint.sh $SCRIPT_DIR/docker/build/entrypoint.sh
}

function buildBot() {
    docker-compose -f $DOCKER_BASE_DIRECTORY/${SERVICES_FILE[1]} down
    docker-compose -f $DOCKER_BASE_DIRECTORY/${SERVICES_FILE[1]} up --build
    docker-compose -f $DOCKER_BASE_DIRECTORY/${SERVICES_FILE[1]} down
}

function updateServiceStatus() {
    if [[ ! -z $1 ]]; then
        file=$DOCKER_BASE_DIRECTORY/${SERVICES_FILE[$1]}

        if [[ -f $file ]]; then
            SERVICES_STATUS[$1]=0
            ret=$(docker-compose -f "$file" --env-file "$SCRIPT_DIR"/../.env top | grep "${SERVICES_PS[$1]}")
            t=$(echo $ret | sed -e 's/\r//g')
            if [[ ! -z $t ]]; then
                SERVICES_STATUS[$1]=1
            fi
        fi
    fi
}

function updateAllServiceStatus() {
    echo "Validating Services status..."
    for i in $(seq 2 $MAX_SERVICES); do
        updateServiceStatus $i
    done
}

function printMenu() {
    clear
    title="Services status"
    echo ""
    echo -e " $On_IPurple$BIWhite$title$Color_Off"
    echo ""
    
    for i in $(seq 1 $MAX_SERVICES); do
        color=$BIRed
        status="not running"
        if [[ "${SERVICES_STATUS[$i]}" == "1" ]]; then
            color=$BIGreen
            status="running"
        fi
        
        space=" "
        if [[ $i -gt 9 ]];then
            space=""
        fi
        
        if [[ "$i" == "1" ]]; then
            color=$BIYellow
            echo -e "\t $space$IWhite 1)$color ${SERVICES_NAME[1]}"        
        else
            echo -e "\t $space$IWhite $i)$color ${SERVICES_NAME[$i]} is $status $Color_Off"
        fi
    done
    echo ""
    echo -n "Type Q to quit, R to refresh services status, L to Log or the number [1-$MAX_SERVICES] of service to change its status: "
    read option
    case $option in
        q|Q) exit ;;
        r|R) main ;;
        l|L) logServicesMenu ;;
        [1-9]) validateChangeServiceStatus $option ;;
        [1-9][0-9]) validateChangeServiceStatus $option ;;
        *) multipleChangeServiceStatus $option ;;
    esac
}

function logServicesMenu() {
    echo ""
    echo -n "Type the number [2-$MAX_SERVICES] of service for logging or Q to quit: "
    read option
    case $option in
        q|Q) main ;;
        [1-9]) logServices $option ;;
        *) logServices ;;
    esac    
}

function multipleChangeServiceStatus() {
    test=`echo $* | grep -E "[0-9]+([\s\,\:\;\-]{1})*"`
    params=$*
    space=" "
    dash="-"
    params=${params//$space/$dash}

    if [[ -n $test ]]; then
        for i in $(echo $params | tr "," "\n"); do
            if [[ ${#params} -eq ${#i} ]]; then
                break
            fi
            doValidateChangeServiceStatus $i
        done

        for i in $(echo $params | tr ";" "\n"); do
            if [[ ${#params} -eq ${#i} ]]; then
                break
            fi
            doValidateChangeServiceStatus $i
        done

        for i in $(echo $params | tr "-" "\n"); do
            if [[ ${#params} -eq ${#i} ]]; then
                break
            fi
            doValidateChangeServiceStatus $i
        done
    else
        echo "Invalid option"        
    fi
    main
}

function validateChangeServiceStatus() {
    createDockerNetwork
    createDockerVolumes
    copyConfigFilesForBuild
    if [[ "$1" == "1" ]]; then
        buildBot
    else
        doValidateChangeServiceStatus $*
    fi
    main
}

function doValidateChangeServiceStatus() {
    if [[ $1 =~ ^([1-9][0-9]?|$MAX_SERVICES)$ ]]; then
        updateServiceStatus $1

        action="START"
        if [[ "${SERVICES_STATUS[$1]}" == "1" ]]; then
            action="STOP"
        fi
        changeServiceStatus $1 $action
    else
        echo "Invalid service!"
    fi
}

function logServices() {
    file=$DOCKER_BASE_DIRECTORY/${SERVICES_FILE[$1]}

    if [[ -f $file ]]; then
        docker-compose -f $file --env-file $SCRIPT_DIR/../.env logs --no-log-prefix -f
    fi
}

function changeServiceStatus() {
    file=$DOCKER_BASE_DIRECTORY/${SERVICES_FILE[$1]}

    if [[ -f $file ]]; then
        if [[ "$2" == "START" ]]; then
            docker-compose -f $file --env-file $SCRIPT_DIR/../.env down
            ret=`docker-compose -f $file --env-file $SCRIPT_DIR/../.env up -d --build`
            SERVICES_STATUS[$1]=1
        else
            ret=`docker-compose -f $file --env-file $SCRIPT_DIR/../.env down`
            SERVICES_STATUS[$1]=0
        fi
    fi
}

function main() {
    updateAllServiceStatus
    printMenu
}

main
