#!/bin/bash

# Prepare git folder
mkdir -p ${HOME} && chmod 0700 ${HOME}
if [ ! -w ${HOME} ]; then echo "${HOME} is not writable"; exit 1; fi

# Prepare custom folder
mkdir -p ${KMUP_CUSTOM} && chmod 0700 ${KMUP_CUSTOM}

# Prepare temp folder
mkdir -p ${KMUP_TEMP} && chmod 0700 ${KMUP_TEMP}
if [ ! -w ${KMUP_TEMP} ]; then echo "${KMUP_TEMP} is not writable"; exit 1; fi

#Prepare config file
if [ ! -f ${KMUP_APP_INI} ]; then

    #Prepare config file folder
    KMUP_APP_INI_DIR=$(dirname ${KMUP_APP_INI})
    mkdir -p ${KMUP_APP_INI_DIR} && chmod 0700 ${KMUP_APP_INI_DIR}
    if [ ! -w ${KMUP_APP_INI_DIR} ]; then echo "${KMUP_APP_INI_DIR} is not writable"; exit 1; fi

    # Set INSTALL_LOCK to true only if SECRET_KEY is not empty and
    # INSTALL_LOCK is empty
    if [ -n "$SECRET_KEY" ] && [ -z "$INSTALL_LOCK" ]; then
        INSTALL_LOCK=true
    fi

    # Substitute the environment variables in the template
    APP_NAME=${APP_NAME:-"Kmup: Git of Kumose"} \
    RUN_MODE=${RUN_MODE:-"prod"} \
    RUN_USER=${USER:-"git"} \
    SSH_DOMAIN=${SSH_DOMAIN:-"localhost"} \
    HTTP_PORT=${HTTP_PORT:-"3326"} \
    ROOT_URL=${ROOT_URL:-""} \
    DISABLE_SSH=${DISABLE_SSH:-"false"} \
    SSH_PORT=${SSH_PORT:-"2222"} \
    SSH_LISTEN_PORT=${SSH_LISTEN_PORT:-$SSH_PORT} \
    DB_TYPE=${DB_TYPE:-"sqlite3"} \
    DB_HOST=${DB_HOST:-"localhost:3306"} \
    DB_NAME=${DB_NAME:-"kmup"} \
    DB_USER=${DB_USER:-"root"} \
    DB_PASSWD=${DB_PASSWD:-""} \
    INSTALL_LOCK=${INSTALL_LOCK:-"false"} \
    DISABLE_REGISTRATION=${DISABLE_REGISTRATION:-"false"} \
    REQUIRE_SIGNIN_VIEW=${REQUIRE_SIGNIN_VIEW:-"false"} \
    SECRET_KEY=${SECRET_KEY:-""} \
    envsubst < /etc/templates/app.ini > ${KMUP_APP_INI}
fi

# Replace app.ini settings with env variables in the form KMUP__SECTION_NAME__KEY_NAME
environment-to-ini --config ${KMUP_APP_INI}
