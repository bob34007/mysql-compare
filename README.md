# mysql-compare
read the files persisted by the simulation tool, compare the execution time 
and results of the production and simulation environments, and output the 
comparison statistics

#比对工具
读取仿真工具持久化的文件，比对生产环境和仿真环境的执行时间和执行结果，并输出比对统计结果

# demo

./mysql-compare  text compare -d "/home/store/" -t 20 -B 2000 -P7700 --log-output="/home/log/mysql-compare-26.log" &