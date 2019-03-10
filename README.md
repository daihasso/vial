# Vial

## Description
Vial is a microframework heavily inspired by flask designed to be REST-y and
familiar to help facilitate microservice development in a familiar fashion.

## Basic Usage Example
As a simple use case let's create an ultra simple controller that simply
responds to a path at `hello/<name>` and responds to the provided name.

``` go
package main

import (
    "github.com/daihasso/vial"
    "github.com/daihasso/vial/neterr"
    "github.com/daihasso/vial/responses"
)

func main() {
    server, err := vial.NewServerDefault()
    if err != nil {
        panic(err)
    }

    server.AddController(
        "/hello/<name>",
        vial.FuncHandler(
            "get",
            func(transactor *vial.Transactor) responses.Data {
                name, ok := transactor.Request.PathString("name")
                if !ok {
                    return transactor.Abort(
                        500,
                        neterr.NewCodedError(
                            1,
                            "Could not retrieve name from path",
                        ),
                    )
                }

                return transactor.Respond(
                    200,
                    responses.Body(map[string]string{
                        "hello": name,
                    }),
                )
            },
        ),
    )

    err := server.Start() // Blocking

    if err != nil {
        panic(err)
    }
}
```

Making sense yet? Let's break it down a little to get a better understanding.

First:

``` go
server, err := vial.NewServerDefault()
```

This one is pretty straightforward; this creates our vial server instance and
does some housekeeping tasks like reading an existing config or generating a
default config in its absence, etc.
By default the server will host on `127.0.0.1:8080`.

Now comes the meat:
``` go
server.AddController(
    "/hello/<name>",
    vial.FuncHandler(
        "get",
        func(transactor *vial.Transactor) responses.Data {
            name, ok := transactor.Request.PathString("name")
            if !ok {
                return transactor.Abort(
                    500,
                    neterr.NewCodedError(
                        1,
                        "Could not retrieve name from path",
                    ),
                )
            }

            return transactor.Respond(
                200,
                responses.Body(map[string]string{
                    "hello": name,
                }),
            )
        },
    ),
)
```

Here we declare our actual controller for our route. First we pass in the route
with `"/hello/<name>"` which responds to any `/hello/foo` route and grabs the
value of `foo` and throws it in the path parameter map which we access later
with the `name, ok := transactor.Request.PathString("name")` call.

Next we provide our method then the actual controller function. When using a
pure functional handler we use the `vial.FuncHandler` which just maps a
function to one or more HTTP methods (in our case GET defined by `"get`" as the
first argument to FuncHandler).
\* **Note**: There is also a `FuncControllerMulti` that lets you define
multiple methods for a single functional controller for convenience.

Our function signature in this case looks like
`func(transactor *vial.Transactor) responses.Data`, there are multiple acceptable
signatures which you can read more about [here][add-controller-signatures].

The `Transactor` we receive as an argument here is a kind of a multi-use tool
for handling interactions within a handler. It helps with things like grabbing
path parameters, setting headers, serializing data and more. You can read more
about `Transactor` here [here][transactor]

In this case we first grab our `name` variable from the path with:
```go
name, ok := transactor.Request.PathString("name")
```

The next part is some pretty standard boiler plate "If we don't get the path
param we expect, return an error".
``` go
return transactor.Abort(
    500,
    neterr.NewCodedError(
        1,
        "Could not retrieve name from path",
    ),
)
```

The `Abort` function on transactor aborts the current request with the given
error code and provides one or more `CodedError`s.
`CodedError` is a simple type which is just a code (some integer) and a message
describing the error. The idea is simple, have coded, explanatory errors every
time you abort a request so that your API users have entropy to understand why
things went wrong.

Finally we respond with our message and a response code and a body:
``` go
return transactor.Respond(
    200,
    responses.Body(map[string]string{
        "hello": name,
    }),
)
```

The first argument to `Respond` is always a response code. You can simply
respond with only a response code, but that wouldn't be very interesting so you
probably want to add something in the response. There are a list of available
optional parameters which add body, headers, etc to the response. Read more
about these optional parameters [here][builder-additionals].
Here we use `responses.Body` to add a body (any interface) to the response. If
the body is something other than a `string` or a `[]byte` object then it will be
serialized in the manor appropriate for the specified encoding. Read more about
encoding types [here][encoding-types] (By default a new server uses JSON
as the default encoding type).

Now, if you run the above code and fire up your favorite web requester (we like
[HTTPie][httpie]) and request `127.0.0.1:8080/hello/tester` you will receive
something like:
``` json
{
    "hello": "tester"
}
```

## Struct Controllers
Controllers used by a `Server` can be either funcs as seen above used by
`vial.FuncHandler` but they can also be controllers like so:

``` go
package main

import (
    "github.com/daihasso/vial"
    "github.com/daihasso/vial/responses"
)

type MyController struct {
    AppName string
}

func (self MyController) Get(transactor *vial.Transactor) responses.Data {
    return transactor.Respond(
        200,
        responses.Body(map[string]string{
            "app_name": self.AppName,
        }),
    )
}

func main() {
    server, err := vial.NewServer()
    if err != nil {
        panic(err)
    }

    server.AddController(
        "/test",
        &MyController{
            AppName: "TestServer",
        },
    )

    err = server.Start() // Blocking

    if err != nil {
        panic(err)
    }
}
package main

import (
    "github.com/daihasso/vial"
    "github.com/daihasso/vial/neterr"
)

type MyController struct {
    AppName string
}

func (self MyController) Get(transactor *Transactor) responses.Data {
    return transactor.Respond(
        200,
        responses.Body(map[string]string{
            "app_name": self.AppName,
        }),
    )
}

func main() {
    server, err := vial.NewServerDefault()
    if err != nil {
        panic(err)
    }

    server.AddController(
        "/test",
        &MyController{
            AppName: "TestServer",
        },
    )

    err := server.Start() // Blocking

    if err != nil {
        panic(err)
    }
}
```

Now fire a request to `127.0.0.1:8080/test` and you should get something back
like:

``` json
{
    "app_name": "TestServer"
}
```

Obviously this example is rather contrived but using a struct has other benefits
such as grouping controllers around routes, and more complex logic and variable
tracking. At the end of the day both are first-class citizens for `Vial` so make
your choice whatever way feels comfortable for you.

## The Sequence ID
A `Sequence ID` is an id that is either grabbed from a header value of
`Sequence-Id` passed in with the request or it is generated at the start of the
request. It is then used to identify further actions down the chain of as a
result of the request. It is a first-class citizen in the `vial` framework and
is added to the context upon the start of the request under the
`vial.sequence_id` key. It is also available via the `Transactor.SequenceId()`
method as well as it being set as a
[DefaultExtra][slogging-default-extras] on the `Transactor.Logger` object.

Its primary purpose is to keep track of information on a given request and
should be passed to any thing that produces logs or any subsequent API calls
made to other servers. The `Sequence ID` will also be returned in the headers
for the request to the api so that the caller can utilize it or a developer can
use it to trace a request in an error case.

## Path Parameters
`Vial` has a path parameter parsing library that lets you define specific
matches in the declaration of your url such as:
* float - `<integer:var_name>`
* intger - `<integer:var_name>`
* UUID - `<uuid:var_name>`
* string - `<string:var_name> or <var_name>`

These can be accessed off the request (through the transactor) and auto-coerced
into appropriate formats with:
* `transactor.Request.PathFloat`
* `transactor.Request.PathInteger`
* `transactor.Request.PathUUID`
* `transactor.Request.PathString`

Any url that doesn't match the expected format will be rejected and another
matcher will be attempted if it exists. In other-words you can have two routes:
`/image/<integer:id>`
and
`/image/<uuid:id>`
And both will be matched to the appropriate calls.

[add-controller-signatures]:
https://godoc.org/github.com/daihasso/vial#Server.AddController
"AddController Godocs"
[transactor]:
https://godoc.org/github.com/daihasso/vial#Transactor
"Transactor Godocs"
[builder-additionals]:
https://godoc.org/github.com/daihasso/vial/responses#AdditionalAttribute
"Builder Additional Attributes Godocs"
[encoding-types]:
https://godoc.org/github.com/daihasso/vial/responses#EncodingType
"Builder Encoding Types Godocs"
[httpie]: https://httpie.org
[slogging-default-extras]:
https://godoc.org/github.com/daihasso/slogging#WithDefaultExtras
"Default Extras In Slogging"
