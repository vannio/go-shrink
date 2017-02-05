# Shrink

A WIP URL-shortener in Golang. Like bit.ly, except using my own domain :tada:

Short URL structure:
`vann.io/s/xxxxxxxx` where `xxxxxxxx` is _obviously_ a unique token

Table currently looks like:

| id | token | url | created_at | last_accessed | access_count |
| --- | --- | --- | --- | --- | --- |
| key | string | string | timestamp | timestamp | integer |
| 1 | 55b207d6 | http://www.google.com | 2017-02-05 12:25:30.044385 | 2017-02-05 10:34:05.186949 | 5 |

#### TODO
- Create some kind of self-executing `sql` script to generate the table
- Implement POST-redirect
- Make it look less ugly
- Maybe think of a better name than _Shrink_
- Get it working on AWS
- Hope no-one abuses it
