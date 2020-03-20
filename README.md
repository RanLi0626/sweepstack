# database:
## redis 
1. award_time[zset] value:award_name  score:winning_time
2. award_remain_num[hash]  field:award_name  value:award_remain_num

# 抽奖算法：
1. 从剩余奖池中随机抽取此次抽奖的奖品
2. 判断本次抽奖是否中奖
   - 平均发放奖品的时间间隔：(end_time - start_time) / prizes_total_amount  -> A
   - end_time - last_win_time 作为随机数种子 -> B
   - A * 已发放奖品数量 -> C
   - start_time + C + 由B得到的随机时间 -> D
   - time.now() >= D -> 中奖
