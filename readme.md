# Services Manager

A very basic web application for managing a select group of services on Ubuntu Linux.

---

## Use as a regular user

`$ ./services-manager`

Runs a web server on [http://localhost:8888](http://localhost:8888).  

By default, you can control Apache2, MongoDB, MySQL and PostgreSQL.

Click 'start' to start that service, or 'stop' to stop it.

Click 'start all' or 'stop all' to do just that.

All actions are logged.  Attempting to manage a service that does not exist will result in an error.

**Note:** will require a password for each service changed, every time.

---

## Use with `sudo` or privileged user

E.g.:

```
$ sudo ./services-manager
[sudo] password for your-username:
```
...or:

```
$ sudo -s
[sudo] password for your-username:
# ./services-manager
```

Works exactly the same as above, but no longer requires a password for each service change.