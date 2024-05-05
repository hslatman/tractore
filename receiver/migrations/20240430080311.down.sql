-- reverse: create index "outbox_topic_id" to table: "outbox"
DROP INDEX "outbox_topic_id";
-- reverse: create "outbox" table
DROP TABLE "outbox";
-- reverse: create "incoming_mails" table
DROP TABLE "incoming_mails";
