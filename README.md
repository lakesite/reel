# reel #

🖭 restore app state, for demos and development	🖭

## motivation ##

If you provide demo sites for an application, you might have a cron job setup to
periodically restore the application's state to some default.  This allows
prospects to try out your service by making changes and then reset to a default
state.

* reel aims to:
  - Restore an application's database state to a configurable default.
  - Restore an application to a preset configuration.
  - Provide a web interface for interactive resets for use with vagrant.
  - Define configuration reset points.
  - Provide per-instance configuration proxy,

    e.g. https://app.biz/demos/<ID>/

    Which is unique to a prospect, has a unique database backed instance of
    an application

* Area of focus:
  - Demos and development.
  - Light-weight reverse proxy to app instances tied to prospects.

Further [rationale](docs/rationale.md) provided.

To see how reel works, see the [architecture](docs/architecture.md).
To see what services reel provides, see the [services](docs/services.md).

## usage ##

For use while developing or maintaining a project, please see [usage](docs/usage.md).

## development ##

To run locally and develop, see [development.md](docs/development.md)

## license ##

MIT - See [LICENSE.md](license.md)

## contributing ##

Please review [standards](docs/standards.md) before submitting issues and pull
requests.  Thank you in advance for feedback, criticism, and feature requests.
