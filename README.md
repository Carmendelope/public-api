# Public API

This component represents the entry point into the Nalej platform for the user requests. Authenticated requests will
be send to this component once the login has been performed, and the component will forward the requests to the
appropriate services.

 ## Server

 To launch the public api execute:

 ```
 $ ./bin/public-api run
 {"level":"info","time":"2018-09-28T15:16:30+02:00","message":"Launching API!"}
 ```

 ## CLI

 A CLI has been added for convenience, use:

 ```
 $ ./bin/public-api-cli
 ```

## Tips for local development

### Embed certificates in minikube

```
$ minikube config set embed-certs true
```

 # Integration tests

 The following table contains the variables that activate the integration tests

 | Variable  | Example Value | Description |
 | ------------- | ------------- |------------- |
 | RUN_INTEGRATION_TEST  | true | Run integration tests |
 | IT_SM_ADDRESS  | localhost:8800 | System Model Address |
 | IT_INFRAMGR_ADDRESS  | localhost:8860 | Infrastructure Manager Address |
 | IT_MONITORING_MANAGER_ADDRESS  | localhost:8423 | Monitoring Manage rAddress |
 | IT_APPMGR_ADDRESS  | localhost:8910 | Applications Manager Address |
 | IT_UL_COORD_ADDRESS | localhost:8323 | Unified Logging Coordinator Address
 | IT_USER_MANAGER_ADDRESS  | localhost:8920 | User Manager Address |
 | IT_DEVICE_MANAGER_ADDRESS | localhost:6010 | Device Manager Address |
