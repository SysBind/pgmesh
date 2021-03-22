# pgmesh - Automate common use cases for Postgres logical replication

## Upgrade Major Postgres Version or Switch Architecure / Deployment

Can be carried out for each database of the cluster seperatly

Ref: https://www.2ndquadrant.com/en/blog/upgrading-to-postgresql-11-with-logical-replication/

### Plan - Repeated for each database
- Import basic schema
- Initiate PubSub
- Monitor replication state
- Enter Maintnence mode - if applicable
- Fixup Sequences
- Redirect traffic to new cluster
- Finalize PubSub
- Exit Maintnence - if applicable - only on new database
- Take down or make Read-Only old database

