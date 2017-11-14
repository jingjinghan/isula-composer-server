#!/usr/bin/env bash

DB_NAME=testdb
DB_USER=testuser
DB_PASS=testpass

echo "USE mysql;\nCREATE DATABASE IF NOT EXISTS ${DB_NAME} DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci;\n" | mysql -u root
echo "USE mysql;\nCREATE USER ${DB_USER}@'%' IDENTIFIED BY '${DB_PASS}';\nFLUSH PRIVILEGES;\n" | mysql -u root 
echo "USE mysql;\nGRANT ALL PRIVILEGES ON ${DB_NAME}.* TO ${DB_USER}@'%' IDENTIFIED BY '${DB_PASS}';\n" | mysql -u root 
