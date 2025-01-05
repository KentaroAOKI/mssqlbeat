# Mssqlbeat

Mssqlbeat is a custom Beat for collecting data from Microsoft SQL Server and sending it to the Elastic Stack. 

Ensure that this folder is at the following location:
`${GOPATH}/src/github.com/KentaroAOKI/mssqlbeat`

## Getting Started with Mssqlbeat

### Requirements

* [Golang](https://golang.org/dl/) 1.7

### Init Project
To get running with Mssqlbeat and also install the
dependencies, run the following command:

```
make setup
```

It will create a clean git history for each major step. Note that you can always rewrite the history if you wish before pushing your changes.

To push Mssqlbeat in the git repository, run the following commands:

```
git remote set-url origin https://github.com/KentaroAOKI/mssqlbeat
git push origin master
```

For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).

### Build

To build the binary for Mssqlbeat run the command below. This will generate a binary
in the same directory with the name mssqlbeat.

```
make
```


### Run

To run Mssqlbeat with debugging output enabled, run:

```
./mssqlbeat -c mssqlbeat.yml -e -d "*"
```


### Test

To test Mssqlbeat, run the following command:

```
make testsuite
```

alternatively:
```
make unit-tests
make system-tests
make integration-tests
make coverage-report
```

The test coverage is reported in the folder `./build/coverage/`

### Update

Each beat has a template for the mapping in elasticsearch and a documentation for the fields
which is automatically generated based on `fields.yml` by running the following command.

```
make update
```


### Cleanup

To clean  Mssqlbeat source code, run the following command:

```
make fmt
```

To clean up the build directory and generated artifacts, run:

```
make clean
```


### Clone

To clone Mssqlbeat from the git repository, run the following commands:

```
mkdir -p ${GOPATH}/src/github.com/KentaroAOKI/mssqlbeat
git clone https://github.com/KentaroAOKI/mssqlbeat ${GOPATH}/src/github.com/KentaroAOKI/mssqlbeat
```


For further development, check out the [beat developer guide](https://www.elastic.co/guide/en/beats/libbeat/current/new-beat.html).


## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.

## Configuration Settings
### period

    Defines the interval at which the SQL query is executed.
    Example: 600s

### threads

    Description: Specifies the number of threads used for data collection.
    Example: 2

### inputs

    Description: Defines the input settings for data collection. Multiple inputs can be configured.

### field

    Specifies the field name for the input.
    Example: "sql1"

### enabled

    Indicates whether the input is enabled.
    Example: true

### mssqlserver_host

    Specifies the hostname of the SQL Server to connect to.
    Example: "servername.database.windows.net"

### mssqlserver_port

    Specifies the port number of the SQL Server to connect to.
    Example: "1433"

### mssqlserver_userid

    Specifies the user ID for connecting to the SQL Server.
    Example: "userid"

### mssqlserver_password

    Specifies the password for connecting to the SQL Server.
    Example: "password"

### mssqlserver_database

    Specifies the name of the database to connect to.
    Example: "db"

### mssqlserver_tlsmin

    Description: Specifies the minimum TLS version to use.
    Example: "1.2"

### sql_query

    Specifies the SQL query to execute. The @LastTime placeholder will be replaced with the latest date.
    Example: "SELECT * FROM TimeCount WHERE CurrentTime > @LastTime"

### sql_time_column

    Specifies the column name used for time-based filtering.
    Example: "CurrentTime"

### sql_time_initialize_with_current_time

    Indicates whether to initialize with the current time.
    Example: false

### field_prefix

    Specifies the prefix for field names when sending data to the Elastic Stack. This prefix is added to the column names from the SQL Server to create the field names.
    Example: "SQL1_"

## Example Configuration File

Below is an example configuration file for mssqlbeat:

```
mssqlbeat:
  period: 10s
  threads: 2
  inputs:
    - field: "sql1"
      enabled: true
      mssqlserver_host: "<server_name>.database.windows.net"
      mssqlserver_port: "1433"
      mssqlserver_userid: "<user_id>"
      mssqlserver_password: "<password>"
      mssqlserver_database: "<database_name>"
      mssqlserver_tlsmin: "1.2
      sql_query: "SELECT * FROM TimeCount WHERE CurrentTime > @LastTime"
      sql_time_column: "CurrentTime"
      sql_time_initialize_with_current_time: true 
      field_prefix: "SQL1_"
    - field: "sql2"
      enabled: true
      mssqlserver_host: "<server_name>.database.windows.net"
      mssqlserver_port: "1433"
      mssqlserver_userid: "<user_id>"
      mssqlserver_password: "<password>"
      mssqlserver_database: "<database_name>"
      mssqlserver_tlsmin: "1.2
      sql_query: "SELECT * FROM TimeCount WHERE CurrentTime > @LastTime"
      sql_time_column: "CurrentTime"
      sql_time_initialize_with_current_time: true 
      field_prefix: "SQL2_"
```
