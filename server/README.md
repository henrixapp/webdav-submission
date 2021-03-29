# webdav-submissions server implementation

based on Minio and database

## Handling of user id

The user id is based through the string "userID"

## Quota

There is a limit of 15 MB per file upload. And per assignment an file limit can be configured!

## ADMIN-API

description: located in [submissions.yml](api/submissions.yml) 

### Planned Integration of (Admin)-API in MaMpf

Please configure a proxy that sets `X-Forwarded-User` header to the id of the user.
The [frontend](../frontend) can be used for it