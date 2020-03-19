#database:
redis 
award_time[zset] value:award_name  score:winning_time
award_remain_num[hash]  field:award_name  value:award_remain_num
