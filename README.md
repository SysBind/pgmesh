# pgmesh - Automate common use cases for Postgres logical replication

## Upgrade Major Postgres Version or Switch Architecure / Deployment

### States
- Import basic schema
- Initiate PubSub
- Monitor replication state
- Enter Maintnence mode - if applicable
- Redirect traffic to new cluster
- Finalize PubSub
- Exit Maintnence - if applicable - only on new cluster
- Take down or make Read-Only old cluster
