# Bigtable Emulator In Docker

This repository contains everything to run the in memory Go implementation of Bigtable within a docker container
for testing applications which use Bigtable.

To run locally use the command:
```
docker run -it -p 8080:8080 spotify/bigtable-emulator:latest
```

You can specify tables and column families to be created at startup by using the `-cf` switch. The format is a comma separated list of `<instance>.<table>.<column family>`. Ex:
```
docker run -it -p 8080:8080 spotify/bigtable-emulator:latest -cf dev.records.data,dev.records.meta
```

## Using with the Google Cloud Bigtable client

To use with the [Cloud Bigtable Client](https://github.com/GoogleCloudPlatform/cloud-bigtable-client) you need to
set the Data Host and Table Admin Host to where the docker container is running, the Port to 8080, and to use the
Plain Text Negotiation since the in memory implementation does not support SSL encryption. It is also highly recommended to use null credentials, otherwise valid ones are required. An example of using it
with [BigtableOptions.java](https://github.com/GoogleCloudPlatform/cloud-bigtable-client/blob/master/bigtable-client-core/src/main/java/com/google/cloud/bigtable/config/BigtableOptions.java):

```java
BigtableOptions.Builder builder = new BigtableOptions.Builder();
// set required connection parameters
// ...
// ...
// Setup connection to docker container running locally
builder.setDataHost("localhost");
builder.setTableAdminHost("localhost");
builder.setPort(8080);
builder.setUsePlaintextNegotiation(true);
builder.setCredentialOptions(CredentialOptions.nullCredential());
BigtableSession session = new BigtableSession(builder.build());
```

## What Works And What Doesn't

Since this uses an in memory implementation of Bigtable, the Bigtable API is not 100% implemented (for example
certain row filter types are not supported). In order to get a full list of that you should look in the [source for the in memory implementation](https://github.com/GoogleCloudPlatform/gcloud-golang/blob/master/bigtable/bttest/inmem.go).
It is highly recommended to check the logs of the docker container if something is not working properly as there will be warning logs whenever an unsupported operation is attempted. An [issue](https://github.com/GoogleCloudPlatform/gcloud-golang/issues/261) has been opended to address this.

## Notes

Currently this uses the Go implementation because there is not an official Bigtable emulator provided in the google
cloud tools. Hopefully once the Bigtable emulator is included within the set of gcloud utlities, this will be changed to use
that implementation.


## Code of conduct
This project adheres to the [Open Code of Conduct][code-of-conduct]. By participating, you are expected to honor this code.

[code-of-conduct]: https://github.com/spotify/code-of-conduct/blob/master/code-of-conduct.md
