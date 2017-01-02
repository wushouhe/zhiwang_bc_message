/*
Navicat MySQL Data Transfer

Source Server         : localhost_3306
Source Server Version : 50155
Source Host           : localhost:3306
Source Database       : ethereum

Target Server Type    : MYSQL
Target Server Version : 50155
File Encoding         : 65001

Date: 2017-01-02 16:18:18
*/

SET FOREIGN_KEY_CHECKS=0;
-- ----------------------------
-- Table structure for `accounts`
-- ----------------------------
DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts` (
  `address` varchar(11) NOT NULL DEFAULT '',
  `balance` varchar(11) DEFAULT NULL,
  PRIMARY KEY (`address`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of accounts
-- ----------------------------

-- ----------------------------
-- Table structure for `blocks`
-- ----------------------------
DROP TABLE IF EXISTS `blocks`;
CREATE TABLE `blocks` (
  `hash` varchar(255) NOT NULL,
  `parenthash` varchar(255) DEFAULT NULL,
  `nonce` varchar(255) DEFAULT NULL,
  `number` varchar(255) DEFAULT NULL,
  `extraData` varchar(255) DEFAULT NULL,
  `gasLimit` varchar(255) DEFAULT NULL,
  `gasUsed` varchar(255) DEFAULT NULL,
  `miner` varchar(255) DEFAULT NULL,
  `mixHash` varchar(255) DEFAULT NULL,
  `receiptsRoot` varchar(255) DEFAULT NULL,
  `stateRoot` varchar(255) DEFAULT NULL,
  `sha3Uncles` varchar(255) DEFAULT NULL,
  `logsBloom` varchar(255) DEFAULT NULL,
  `size` varchar(255) DEFAULT NULL,
  `difficulty` varchar(255) DEFAULT NULL,
  `totalDifficulty` varchar(255) DEFAULT NULL,
  `timestamp` varchar(255) DEFAULT NULL,
  `transactionsRoot` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;


-- ----------------------------
-- Table structure for `events`
-- ----------------------------
DROP TABLE IF EXISTS `events`;
CREATE TABLE `events` (
  `id` varchar(11) DEFAULT NULL,
  `block_id` varchar(11) DEFAULT NULL,
  `block_number` varchar(11) DEFAULT NULL,
  `log_index` varchar(11) DEFAULT NULL,
  `transaction_id` varchar(11) DEFAULT NULL,
  `contract_id` varchar(11) DEFAULT NULL,
  `name` varchar(11) DEFAULT NULL,
  `event_data` varchar(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of events
-- ----------------------------

-- ----------------------------
-- Table structure for `peers`
-- ----------------------------
DROP TABLE IF EXISTS `peers`;
CREATE TABLE `peers` (
  `id` varchar(11) DEFAULT NULL,
  `status` varchar(11) DEFAULT NULL,
  `node_url` varchar(11) DEFAULT NULL,
  `node_name` varchar(11) DEFAULT NULL,
  `node_ip` varchar(11) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of peers
-- ----------------------------

-- ----------------------------
-- Table structure for `transactions`
-- ----------------------------
DROP TABLE IF EXISTS `transactions`;
CREATE TABLE `transactions` (
  `hash` varchar(255) NOT NULL DEFAULT '',
  `blockHash` varchar(255) DEFAULT NULL,
  `blockNumber` varchar(255) DEFAULT NULL,
  `tx_from` varchar(255) DEFAULT NULL,
  `tx_to` varchar(255) DEFAULT NULL,
  `value` varchar(255) DEFAULT NULL,
  `input` varchar(255) DEFAULT NULL,
  `nonce` varchar(255) DEFAULT NULL,
  `transactionIndex` varchar(255) DEFAULT NULL,
  `gas` varchar(255) DEFAULT NULL,
  `gasPrice` varchar(255) DEFAULT NULL,
  `v` varchar(255) DEFAULT NULL,
  `r` varchar(255) DEFAULT NULL,
  `s` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`hash`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of transactions
-- ----------------------------
