## /internal
contains logic of all the api calls
internal only has private methods to module

##pkg/*
pkg only has public methods of module
it can also act as proxies between internal methods and users

## pkg/services
contains interface between logic and handler functions

## pkg/transport
contains handler function with API details

## pkg/views
contains structs for defined data structures

##db
contains database related files

##src
contains main function to run the server file

##cmd
contains compiled binary

- [ ] make factory for DB or Cookie caching

    CookieCache struct {
        c *gin.Context
    }

    DB struct {
        conn *pgxpool.Pool
        ctx *context.Context
    }

    create an interface will all the functions
