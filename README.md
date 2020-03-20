#database:
> redis 
award_time[zset] value:award_name  score:winning_time
award_remain_num[hash]  field:award_name  value:award_remain_num

#抽奖算法：
1. 从剩余奖池中随机抽取此次抽奖的奖品
2. 判断本次抽奖是否中奖
   1）平均发放奖品的时间间隔：(end_time - start_time) / prizes_total_amount  -> A
   2) end_time - last_win_time 作为随机数种子 -> B
   3）A * 已发放奖品数量 -> C
   4) start_time + C + 由B得到的随机时间 -> D
   5）time.now() >= D -> 中奖
