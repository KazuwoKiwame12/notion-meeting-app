#!/bin/bash
getHerokuAppName() {
    cat ./golang/.env.heroku | while read line || [ -n "${line}" ]; do
    if test "`echo ${line} | cut -d '=' -f 1`" = "HEROKU_APP_NAME";then
        echo "`echo ${line} | cut -d '=' -f 2`"
        break
    fi
done
}
heroku_app_name=`getHerokuAppName`

cat ./golang/.env | while read line; do
    if test "`echo ${line} | cut -c 1`" != "#";then
        key=`echo ${line} | cut -d '=' -f 1`
        value=`echo ${line} | cut -d '=' -f 2`
        heroku config:set ${key}=${value} --app ${heroku_app_name}
    fi
done