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
 $ ./bin/system-model-cli
 ```

 # Integration tests

 The following table contains the variables that activate the integration tests

 | Variable  | Example Value | Description |
 | ------------- | ------------- |------------- |
 | RUN_INTEGRATION_TEST  | true | Run integration tests |
 | IT_SM_ADDRESS  | localhost:8800 | System Model Address |
 | IT_INFRAMGR_ADDRESS  | localhost:8860 | Infrastructure Manager Address |
 | IT_APPMGR_ADDRESS  | localhost:8910 | Applications Manager Address |
 | IT_ACCESSMGR_ADDRESS  | localhost:8920 | Access Manager Address |
