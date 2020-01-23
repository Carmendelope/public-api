# Public API

This component represents the entry point into the Nalej platform for user requests. Once the user is logged in the platform, their authenticated requests will be sent to this component, and it will adapt them and forward them to the appropriate internal component of the management cluster. It will also receive the response, and transform and prepare it for the user.

Notice that the `public-api` supports REST and gRPC requests by means of the [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway) that is launched when the server starts.

It is important to highlight that the `public-api` must transform gRPC internal entities into web compatible ones for
those clients connecting through the REST endpoint. Therefore, some of the entities are enhanced by the `public-api`
server, so they are able to transform internal types such as enumerations into web compatible string values. Future versions of the
`public-api` will support specific REST and gRPC input/output entities and will provide decorator operations (e.g., sorting)
on the returned results.
​
## Getting Started
​
The `public-api` component is the entry point for user request on the platform. The user needs to log in
the platform before being able to execute any action as the JWT is checked by the `authx-interceptor`.
​
### Prerequisites

* A valid deployment of the whole Nalej platform on the management cluster.
​​
### Build and compile
​
In order to build and compile this repository use the provided Makefile:
​
```
make all
```
​
This operation generates the binaries for this repo, downloads the required dependencies, runs existing tests and generates ready-to-deploy Kubernetes files.
​
### Run tests
​
Tests are executed using Ginkgo. To run all the available tests:
​
```
make test
```
​
### Update dependencies
​
Dependencies are managed using Godep. For an automatic dependencies download use:
​
```
make dep
```
​
In order to have all dependencies up-to-date run:
​
```
dep ensure -update -v
```

### Integration tests

Some integration test are included. To execute those, setup the following environment variables. The execution of
integration tests may have collateral effects on the state of the platform. DO NOT execute those tests in production.

The following table contains the variables that activate the integration tests

| Variable  | Example Value | Description |
| ------------- | ------------- |------------- |
| RUN_INTEGRATION_TEST  | true | Run integration tests |
| IT_SM_ADDRESS  | localhost:8800 | System Model Address |
| IT_INFRAMGR_ADDRESS  | localhost:8860 | Infrastructure Manager Address |
| IT_MONITORING_MANAGER_ADDRESS  | localhost:8423 | Monitoring Manager Address |
| IT_APPMGR_ADDRESS  | localhost:8910 | Applications Manager Address |
| IT_UL_COORD_ADDRESS | localhost:8323 | Unified Logging Coordinator Address
| IT_USER_MANAGER_ADDRESS  | localhost:8920 | User Manager Address |
| IT_DEVICE_MANAGER_ADDRESS | localhost:6010 | Device Manager Address |

## User client interface

To interact with the platform, we offer a Public API CLI that supports all the operations the users can
execute on the platform. Use the following commands to setup the cli, so that less parameters are required.

```
$ ./bin/public-api-cli options update platform <platform_url_without_api_prefix>
$ ./bin/public-api-cli options set --key=output --value=table
```

Next, log into the platform with the user credentials.

```
$ ./bin/public-api-cli login --email <user_email> --password <user_password>
EMAIL          ROLE          ORG_ID                                 EXPIRES
<user_email>   <user_role>   6b735d0c-5987-4f11-bbf5-f133c5efe076   2019-11-07 13:17:36 +0100 CET
```

After this, the user can execute any of the commands. Notice that some commands may fail due to the user
having insufficient priviledges to perform a particular action. Use the CLI help and [platform documentation](https://nalej.gitbook.io) to discover the available commands.

```
$ ./bin/public-api-cli --help
```
​
## Known Issues
​
* The `public-api-cli2` CLI sets the foundation for the refactor of the CLI planned for future releases of the platform.
* AuthX related functionality such as the interceptors will be moved to the specific [authx-interceptors](https://github.com/nalej/authx-interceptors) repository in future releases.
* The REST gateway returns gRPC transformed JSON objects. Therefore some attributes may be omitted when they have the gRPC default value. This behavior
will be modified in future releases to return a fully qualified JSON.
* Some operations such as uploading an application descriptor require passing the internal gRPC entity. A future refactor will ensure that all requests
belong to the `public-api` entities.
* Some of the CLI commands support both parameters and flags for the same attribute. This behavior is being refactored to favor parameters when the value
is required.
​
## Contributing
​
Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.
​
## Versioning
​
We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/public-api/tags). 
​
## Authors
​
See also the list of [contributors](https://github.com/nalej/public-api/contributors) who participated in this project.
​
## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.

