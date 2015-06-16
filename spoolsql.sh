#!/bin/sh
 
#define the file location
SQL_DIR=~/zhaoyf/
 
#spool the result to
#spool options
#set colsep ' ' #delimiter
#set echo off #sql command show
#set feedback off #sqlline count
#set heading off 
#set pagesize 0 #devide page
#set termout off #terminal result
#set trimout on #every line have last space stdout
#set trimspool on #every line have last space
#set verify  on
#SET NEWPAGE NONE
#SET SPACE 0
#SET LINESIZE 2500
#set num 18
#set long 300
#set markup html on
#set entmap ON
#set spool on
#set preformat off
#in test.sql select col1||','||col2||','||col3||','||col4||'..' from tablename;
sqlplus acc01_gwtst/1003@XXXGW 1>sql.log 2>&1 <<EOF
set echo off;
set linesize 1024
set term on verify off feedback off pagesize 40000 heading on
set trimspool on
----------------
spool $SQL_DIR/test.txt;
@$SQL_DIR/test.sql
 
spool off;
exit;
EOF
