#!/bin/bash

WRK_SCRIPT="genegate_user_table_data.lua"
WRK_URL="http://217.16.20.177/auth/register"
# WRK_URL="http://localhost:8080/auth/register"

wrk -t100 -c150 -d150s -s "$WRK_SCRIPT" --latency "$WRK_URL" 
