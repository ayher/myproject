# EOS

## 测试网
浏览器：https://jungle.bloks.io/account/karlmarxbest#keys

水龙头：http://monitor.jungletestnet.io/#home

rpc：https://jungle3.cryptolions.io/v1/chain/get_info

## 账户
```bigquery
账户

karlmarxbest

Public Key: EOS7mR7jwMgc68Bf1azYhG9qxkNxWifVN1aMkCMsYVzH33xa9mYWi
Private key: 5JJkAHR6CBPbJsvKbZbx8VRdrhaVCoGZcsHrpjfWWdGh5fwqd7J

Public Key: EOS6XsuncnNJGB6J9fJxZSwJQP2hsWcjcEg2imybkSK5WVUVCeuFs
Private key: 5JETSXKfw8bX9FbXhczqwuRKffveSQsiBtxkcNNTMLTPdUZbGwd

executed transaction: a191086854a4a923c07bd9cbf340d7ae7e6cf6bcd56fc48ea55918f78c7d58f8 336 bytes 480 us warn 2021-05-16T01:59:41.597 cleos main.cpp:506 print_result ] warning: transaction executed locally, but may not be confirmed by the network yet

active2：
Public Key: EOS8Nr2xPFZnmNSmDhdztsqaP3wKKXsJbPSSs6cQxpmuHrbnkYgKW
Private key: 5Htkk3DNYz41Cnmrv2rSyJJUSpaaYdRJMzQHjiqFR3oWGNnQyyT

# eosio <= eosio::newaccount {"creator":"junglefaucet","name":"karlmarxbest","owner":{"threshold":1,"keys":[{"key":"EOS7mR7jwMgc6... # eosio <= eosio::buyrambytes {"payer":"junglefaucet","receiver":"karlmarxbest","bytes":4096} # eosio <= eosio::delegatebw {"from":"junglefaucet","receiver":"karlmarxbest","stake_net_quantity":"1.0000 EOS","stake_cpu_quanti... # eosio.token <= eosio.token::transfer {"from":"junglefaucet","to":"eosio.ram","quantity":"0.0796 EOS","memo":"buy ram"} # eosio.token <= eosio.token::transfer {"from":"junglefaucet","to":"eosio.ramfee","quantity":"0.0005 EOS","memo":"ram fee"} # eosio.token <= eosio.token::transfer {"from":"eosio.ramfee","to":"eosio.rex","quantity":"0.0005 EOS","memo":"transfer from eosio.ramfee t... # junglefaucet <= eosio.token::transfer {"from":"junglefaucet","to":"eosio.ram","quantity":"0.0796 EOS","memo":"buy ram"} # eosio.ram <= eosio.token::transfer {"from":"junglefaucet","to":"eosio.ram","quantity":"0.0796 EOS","memo":"buy ram"} # junglefaucet <= eosio.token::transfer {"from":"junglefaucet","to":"eosio.ramfee","quantity":"0.0005 EOS","memo":"ram fee"} # eosio.ramfee <= eosio.token::transfer {"from":"junglefaucet","to":"eosio.ramfee","quantity":"0.0005 EOS","memo":"ram fee"} # eosio.ramfee <= eosio.token::transfer {"from":"eosio.ramfee","to":"eosio.rex","quantity":"0.0005 EOS","memo":"transfer from eosio.ramfee t... # eosio.rex <= eosio.token::transfer {"from":"eosio.ramfee","to":"eosio.rex","quantity":"0.0005 EOS","memo":"transfer from eosio.ramfee t... # eosio.token <= eosio.token::transfer {"from":"junglefaucet","to":"eosio.stake","quantity":"2.0000 EOS","memo":"stake bandwidth"} # junglefaucet <= eosio.token::transfer {"from":"junglefaucet","to":"eosio.stake","quantity":"2.0000 EOS","memo":"stake bandwidth"} # eosio.stake <= eosio.token::transfer {"from":"junglefaucet","to":"eosio.stake","quantity":"2.0000 EOS","memo":"stake bandwidth"}


karlmarxtest

Public Key: EOS7Kg71LfoU4Hjr5KWQHt4Npsz3kr6KtnC53frEzMY1Zb52t4Hn4
Private key: 5K2BDTSGTTi6pu6uNySGMyHVCzABujtixD2bwsoRy3rgCdBPhui

Public Key: EOS88D7r9hUJn9Pouz3yUnWXV8ZRpT6R2Mch1ynzWEvPbWZR8LzKd
Private key: 5JQmtV3KTfzjPRDsuWGdShSMCBLUbryUdFWZXp3KY4U96ZF4nRj

```

添加 cpu ：https://www.jianshu.com/p/c64d33d457b6
##

## eosio
```bigquery
[root@iZwz9fbvkcoj7mzb5cqa3wZ bin]# cleos wallet create --to-console
Creating wallet: default
Save password to use in the future to unlock this wallet.
Without password imported keys will not be retrievable.
"PW5K91HEdM1L5pASYZeysUN5NkDPbfvupZBTXhVzEVQzfxuJpPT7a"

cleos wallet create -n panyi --to-console
Creating wallet: panyi
Save password to use in the future to unlock this wallet.
Without password imported keys will not be retrievable.
"PW5JMeCKe5Qo11tAavm2V5JKMQqamPrRt5koiatyKBHgabD5GACEM"

修改私钥：https://www.jianshu.com/p/c6bd914dc2d0
cleos -u https://jungle3.cryptolions.io set account permission karlmarxbest active '{"threshold":1,"keys":[{"key":"EOS8Nr2xPFZnmNSmDhdztsqaP3wKKXsJbPSSs6cQxpmuHrbnkYgKW","weight":1}]}'owner
```

