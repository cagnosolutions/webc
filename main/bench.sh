#!/bin/bash

HOST="http://192.168.0.81:8080"

PATH[0]="/index"
PATH[1]="/home"
PATH[2]="/404"
PATH[3]="/user"
PATH[4]="/user/add"
PATH[5]="/user/"
PATH[6]="/"
PATH[7]="/login/"
PATH[8]="/logout/"
PATH[9]="/protected/"

SLUG[0]="foo"
SLUG[1]="bar"
SLUG[2]="baz"
SLUG[3]="nei-inc"
SLUG[4]="cagno-solutions-llc"
SLUG[5]="the-blue-lotus"
SLUG[6]="black-light"
SLUG[7]="0"
SLUG[8]="238"
SLUG[9]="awesome"
SLUG[10]="douche-bag"
SLUG[11]="coolnessinabucket"

I="0"
COUNT="5000"

while [ $I -lt $COUNT ]; do
	PATH_ID=$((RANDOM%10))
	HOSTNAME=${HOST}${PATH[$PATH_ID]}
	if [ $PATH_ID -gt 1 ]; then
		HOSTNAME+=${SLUG[$((RANDOM%12))]}
	fi
	#echo "curl -X GET $HOSTNAME" ## just logging output
	/usr/bin/curl -so /dev/null $HOSTNAME
	I=$[$I+1]
done
