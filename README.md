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