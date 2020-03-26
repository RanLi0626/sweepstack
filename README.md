# 项目简介
1. 语言：Golang
2. 数据库：redis + mysql

# 项目功能流程
1. 按照conf.yaml中的配置初始化数据库。
2. 用户调用{host}/draw?username={username}进行抽奖。
3. 产生结果：win 或 failed

# database
## redis 
1. award_time[hash] value:award_name  score:winning_time
2. award_remain_num[zset]  field:award_name  value:award_remain_num

# 抽奖算法
1. 从剩余奖池中随机抽取此轮中奖的奖品
2. 判断本轮抽奖是否抽中
   - 平均发放奖品的时间间隔：(end_time - start_time) / prizes_total_amount  -> A
   - end_time - last_win_time 作为随机数种子 -> B
   - A * 已发放奖品数量 -> C
   - start_time + C + 由B得到的随机时间%A -> D
   - time.now() >= D -> 中奖

# 使用到的知识点
1. Go启动一个web服务器
2. Go连接redis,mysql及操作
3. time包的使用，时区的注意
4. 读取yaml配置文件
