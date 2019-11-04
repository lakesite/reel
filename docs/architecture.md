# architecture #

a high-level overview of what might be

## level 1 ##

```

HTTP requests for /app/ or /app/demos/:

   [1] [2] [3]
    |   |   |
    |   |   |
    |   |   |
+---V---V---V---+
| Server or     |
| container     |          
| nginx       --+------> [1b]
|               |     
+---------------+     


```
[1] - The request is proxied through normally to an always-on instance of the
      application, a demo, which has its state managed by reel.

      Typical configuration would be:
      /                  - Any default web request not defined in [2] or [3].

[2] - The request is proxied through to an instance of the application tied to
      a previously configured unique ID associated with a prospect.

      Typical configuration:

      /demos/<id>        - Any non-default demo request.  If ID is not defined,
                           redirect to /

[3] - The request is a management console for reel.  This allows us to reset a
      configuration state for an application manually, enter an ID
      (and user/pass) for a demo of an app instance and tie it to a new
      database.

      Typical configuration:

      /reel/              - Management requests for reel itself.

## level 2 ##

A little more thought:

* We have two modes of use, in that we're either developing an application
  locally and want to manage the app state in a trivial way, or we're providing
  a demonstration either in general or for a specific prospect.

  If we're developing locally, we might be using Docker, or Vagrant, and in
  either case we've got a database backed application.  We may be working in a
  specific feature branch, where we're bringing the application up for
  development and want a standard convention for database fixture available,
  without even having to specify one, e.g.,

  ``/path/to/app/reel/<feature branch name>_database.sql``

  If we're providing a demo, we have two use cases:

  1. We've got a general demonstration, presumably available to the public, and
     cron is managing a database/app reset of some kind.  We provide this using
     reel and the app's reel configuration by default, e.g.,

     ``/path/to/app/reel/reel.config``

     This configuration specifies the state, assets and database dumpfile
     location to use.  We use a dumpfile from either mysql or postgresql for
     now.

  2. We're providing a specific demonstration to a prospect.  The ID used ties
     us to the prospect, and is used to configure the application such that (for
     example) the graphics/logo and text are tailored to the prospect.  Further,
     we don't actually configure the instance until the request is made.  This
     means that the first time the client visits, we provide them with a screen
     which communicates "Please wait while we prepare your demonstration," which
     finishes by proxying them off to the instance when it's up and running.

     To start, we may simply require the app be up and running.

## level 3 ##

The heavy lifting of this app will be the reverse proxy and the app manager,
since we're essentially providing a complete solution that combines nginx and
supervisor/systemd in a sandbox environment ONLY for development or single-use
demonstrations.
