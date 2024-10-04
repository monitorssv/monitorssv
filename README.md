# MonitorSSV
An information display and monitoring alarm platform for the SSV network. Monitor the execution layer contract status, monitor the consensus layer validator status, and provide valuable information display and alarms for cluster owners.

## Introduction
For validator clusters running on the SSV network, cluster owners usually pay attention to operator fee changes, cluster available balance, operating runway and other data, which are related to the health and stability of the validator cluster and affect the user's validator's income.

MonitorSSV monitors all cluster data of the SSV network and displays it publicly, and allows cluster owners to configure monitoring strategies, regularly sends cluster status data that users are concerned about, and issues alarms when operator rates change and cluster balances reach the liquidation threshold.

MonitorSSV allows cluster owners to configure monitoring policies,
When the operator rate changes and the cluster runway reaches the alarm threshold, an alarm is triggered.
Regularly send cluster status data, including the number of validators and the current available balance of the cluster.
Supports multiple alarm methods, telegram, discord, etc.

## Functions
Display of SSV Network information
- total number of operators
- total number of validators
- network fee
- cluster size limit
- Liquidation threshold
- cluster minimum collateral

Display of all cluster information
- clusterOwner
- operatorIds
- operatorsFee
- validatorCount
- cluster status
- feeRecipientAddress
- clusterBalance
- runwayDay

Display of operator information

Display of validator information

Search cluster information through clusterOwner/clusterID

Search operator information through Owner/ID

Monitoring consensus layer validator performance

Owner can configure monitoring service
- Support discord/telegram alarm
- Allow setting cluster liquidation runway alarm threshold
- Report when operator fee change
- Report when network fee change
- The validator proposed a block.
- The validator missed a block
- The validator balance decreased or even slashed.

## License
MIT
