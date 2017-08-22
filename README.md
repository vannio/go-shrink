# Shrink

A WIP URL-shortener in Go :tada:

Short URL structure:
`domain.tld/xxxxxxxx` where `xxxxxxxx` is a unique slug

Table currently looks like:

| id | slug | url | created_at | last_accessed | access_count |
| --- | --- | --- | --- | --- | --- |
| key | string | string | timestamp | timestamp | integer |
| 1 | 55b207d6 | http://www.google.com | 2017-02-05 12:25:30.044385 | 2017-02-05 10:34:05.186949 | 5 |
