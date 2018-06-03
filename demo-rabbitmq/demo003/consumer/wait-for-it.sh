#!/usr/bin/env bash

echo "等待5秒"

echo "等待consulManger 把配置文件加载完成"
sleep 5s


echo "执行的文件名称是:"$1

$1