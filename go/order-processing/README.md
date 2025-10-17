## Request for Proposal (kinda)

### Order Processing System Implementation in Go

This program should be able to process a big amount of orders and make their corresponding purchase orders to the providers. It is very important to ensure that all the concurrent processes are secure and they should work well, no matter the amount of workers that are set to do the jobs. 

The basic process should be the following: 
1. **Searching orders to be processed:**
    1. A worker should be polling the DB to query new orders ready to be processed. 
2. **Pending and processing orders queues:**
    1. Once the query is done, they should be filtered and inserted into a queue of pending orders. Also, it should check if the orders are not in the processing orders queue too.
    2. Then another worker should take the pending orders from the queue and mark them as being processed in the db. Once that is done, they should be moved from the pending to the processing orders queue. 
3. **Order processing:**
    1. From the processing queue a new worker should take each one and create their corresponsing purchase orders with the providers. 
    2. Once is done, the order should be marked as completed in the DB and only then should be pop out from the processing orders queue. 

This software should guarantee some important points: 
- After querying pending orders there has to be a validation to prevent the inclusion of duplicated orders. 
- It is crusial that many workers do not process one specific order at the same time. 