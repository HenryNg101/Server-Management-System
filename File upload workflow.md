1. User import servers by making a request: POST /jobs/import-server with csv file

In **jobs service**:

2. Unique Job ID is generated using UUID, and returned to the user in the response for further inspection that user can do
3. File is uploaded to MinIO with an object key (In this case, use "jobs/<job-id>/input.csv")
4. A new job with the above ID is created inside Postgres table for import jobs
5. A message for the import job is published to Kafka

In **files-import-consumer worker**:

6. Files-import-consumer worker is a consumer, which is a long-running application, will consume messages, skip poison messages (i.e messages that are invalid JSON that can't be validated), and process. After processing and all side effects are succeeded (DB updates, MinIO upload), it will commit that one message manually (This is final, so from next step onwards, its about the processing of the message)
7. Processing the message: Check if the job status is done or not. If it's done -> Processed before already but couldn't commit yet, maybe because of system being down -> No process again, just commit -> Ensure idempotency. But otherwise, update the status of the job as "Processing"
8. Download the file from MinIO
9. Import servers: It process the CSV line by line until EOF. For each row, if invalid or can't be parsed, thats a failed row. If yes, add to batch. If the batch is large enough, or there's no more line (reached EOF) to process, the batch is inserted in batch into DB. Then, the amount of success/faild rows will be added to the result
10. Batch upsert: It will try upsert in batch. If there's any error (i.e at least one row has problem, maybe invalid data to fail validation or smt, but definitely not duplicate, because it's...upsert), it will fallback and insert one by one to record failed rows
11. Write to output: After every rows are processed and upserted into DB, the final result of failed/success rows counts are updated to the DB for the job system, and the specific failed rows are written into a JSON file, uploaded to MinIO, and finally, in the worker (Yeah, the consumer for this), Kafka message is committed

Job status check (Within the **jobs service**): GET /job/:id

1. If job doesn't exist, it will be notified to users -> Ends. If it does, continue getting more info
2. It will then try getting a temporary presigned URL from MinIO and outputting it to allow users be able to download the file through that URL. The system has two addresses: One public for presign URL, and one is private for internal usage of deleting/uploading files (In the future, this could be changed to use smt like a reverse proxy to help dealing with file download)
3. It will also try doing the same thing above to get users the output file, but if that doesn't exist yet, because the processing of the input csv file hasn't been done yet, this will be empty and ommitted
4. Output it to the users

Jobs and files cleanup: 

**cron-scheduler**, which is a light cron job is run every 24 hrs to check all the jobs and related files on MinIO. If they are created more than 7 days ago, delete them