CREATE SCHEMA `sweepstake` DEFAULT CHARACTER SET utf8 ;

CREATE TABLE `sweepstake`.`winner_record` (
  `username` VARCHAR(255) NULL,
  `award` VARCHAR(45) NULL,
  `time` TIMESTAMP NULL);