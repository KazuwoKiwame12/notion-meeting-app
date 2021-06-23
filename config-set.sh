#!/bin/bash
heroku_app_name=""
cat ./golang/.env.heroku | while read line; do
    if test "`echo ${line} | cut -c 1`" = "HEROKU_APP_NAME";then
        heroku_app_name=`echo ${line} | cut -d '=' -f 2`
        break
    fi
done

echo ${heroku_app_name}

cat ./golang/.env | while read line; do
    if test "`echo ${line} | cut -c 1`" != "#";then
        key=`echo ${line} | cut -d '=' -f 1`
        value=`echo ${line} | cut -d '=' -f 2`
        heroku config:set ${key}=${value} --app ${heroku_app_name}
    fi
done