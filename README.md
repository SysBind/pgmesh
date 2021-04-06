# pgmesh - Automate common use cases for Postgres logical replication


## Build
```make pgmesh```

## Test
```make e2e```

Currently contains one test - online upgrade/migrate Moodle instance

## Upgrade Major Postgres Version or Switch Architecure / Deployment

Can be carried out for each database of the cluster seperatly

Ref: https://www.2ndquadrant.com/en/blog/upgrading-to-postgresql-11-with-logical-replication/


[Moodle's e2e Test](e2e/moodle/run.sh) Outlines the upgrade procedure:
(For a single database)

- ```pgmesh pubsub```
  - Copy basic schema
  - Copy primary keys
  - Initiate Logical Replication (Publish/Subscribe)
  - Copy all other constraints 
     (not copied before to speedup initial data copy)  
- Enter Maintnence mode - i.e: don't accept new meaningful data
- ```pgmesh pubsub --detach```
  Detaches previously established logical replication PubSub
- ```pgmesh copyseq```
  Copy all sequence values (not being replicated in logical replication)
  Optionally add some slack (--slack=N, added to each copied value)
  To avoid sequence number conflicts, N depends on site business
  (in moodle, only the log table is expected to recieve updated
   in maintnence mode, so for 1000 users site --slack=1000 should
   be more than enough, if things carried out in automatic/fast manner)
- Redirect traffic to new cluster
  (in moodle, change config.php)
- Exit Maintnence - on new database 
  (Moodle: change value in mdl_config)


The test runs few docker containers:
- postgres10 - to upgrade from
- postgres13 - to upgrade to
- moodle - the application
- pgmesh - to establish / teardown the logical replication
           and copy sequence values

_NOTE:_ If you test it on your postgres installation,
         Dont use the same cluster as source and destination,
		 It will get stuck (see: https://stackoverflow.com/questions/59492957/can-i-have-a-logical-replication-between-2-databases-in-the-same-postgres-server)
