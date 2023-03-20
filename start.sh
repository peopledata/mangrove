#! /bin/bash
/app/server &
wait # 等待回调执行完，主进程再退出
