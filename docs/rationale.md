# rationale #

So, why create reel?

There are two main problems we're attempting to solve here.  One is automating
the process of restoring or updating a database for an application, in
development.  The other is automating demonstrations for an app, either for the
general public or a particular prospect.

  I. Developing and maintaining an application.

  This is particularly useful when our workflow is:

    1. vagrant up, on the master branch
    2. git checkout -b <feature> from the master branch.
    3. During our work we modify the database, or create a migration, and need to
       undo our change and see if the migration and feature works.

  II. Providing a demonstration.

  The manual work done involves:

    1. Creating a supervisor or systemd configuration to setup a demo instance
       of an application.
    2. Creating an nginx configuration to proxy requests to the application.
    3. Creating a cron job to reset the database state for an application.

    and:

    4. Repeating 1-3 for a specific instance tied to a specific prospect.

Reel should be able to provide a seamless solution to these problems by
convention in development and configuration as a demonstration.
