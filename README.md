# Shrink

A WIP URL-shortener in Golang :tada:

Short URL structure:
`domain.tld/s/xxxxxxxx` where `xxxxxxxx` is a unique token

Table currently looks like:

| id | token | url | created_at | last_accessed | access_count |
| --- | --- | --- | --- | --- | --- |
| key | string | string | timestamp | timestamp | integer |
| 1 | 55b207d6 | http://www.google.com | 2017-02-05 12:25:30.044385 | 2017-02-05 10:34:05.186949 | 5 |
