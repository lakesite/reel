# usage #

The development use case for reel is ready for use.  The demo case has not been
completed, although reel can be used with cron to reset an app's database state.

For help on the command line:

 $ ./reel -h

## Development ##

1. If you use Vagrant, expose port 7999 in the guest so web requests can be made.
2. Build and copy /bin/* into your project, under, for example, /utils.
3. Configure your example.toml for your app(s), e.g.;

example.toml
```
[reel]
token="{!#}%p13^JTNG$H@R@$HTBgBNbx/Y$&N21"
default_app="app"

[app]
dbserver = "127.0.0.1"
dbport = "3306"
database = "example"
dbuser = "example"
dbpassword = "example"
dbdriver = "mysql"
dbsource = "./dumpfiles/example.sql"
dbsources = "./dumpfiles/"
proxy_dest = "https://localhost:8101/"

```

4. Make sure /utils/dumpfiles exists with the exports of your database, including
example.sql.

5. Run reel in vagrant dynamically, or access via the web API.

  a. Run interactively;
  $ vagrant ssh
  $ cd /vagrant/utils
  $ ./reel -h

  b. Start the web service and run with web requests.
  $ vagrant ssh
  $ cd /vagrant/utils
  $ ./reel -c example.toml manager

  (from the host) $ curl http://localhost:7999/reel/api/v1/rewind/app

## API ##

Current API endpoints:

/reel/                             - Management requests for reel itself.

/reel/api/v1/sources/              - List of database source files for
                                     default app.
/reel/api/v1/sources/{app}         - List of database source files for
                                     {app}.
/reel/api/v1/rewind/               - Rewind the default app.
/reel/api/v1/rewind/{app}          - Rewind {app}
/reel/api/v1/rewind/{app}/{source} - Rewind {app} with {source}
/reel/api/v1/proxy/{app}           - Proxy requests for {app}
