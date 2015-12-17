#!/bin/bash

INV=$(pwd)/../involucro

set -e
$INV -e "inv.task('test').using('busybox').run('/bin/echo', VAR['k'])" -s k=asd test 2>&1| grep "stdout: asd"
$INV -e "inv.task('test').using('busybox').run('/bin/echo', VAR['k'])" --set k=asd test 2>&1| grep "stdout: asd"
$INV -e "inv.task('test').using('busybox').run('/bin/echo', VAR['k'])" --set k=asd=6 test 2>&1| grep "stdout: asd=6"
