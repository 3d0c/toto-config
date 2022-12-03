# TOTO Configuration Server

## Disclaimer
In order this is a Testing Task to simplify testing, deployment and delivery initially it uses SQLite as a database backend. It obviously lock the Server into single instance.  
To scale it into multiple instances of the Server, SQLite should be replaced by any RDBMS supported network access. In most cases only DSN should be changed along with Open function dialect parameter.

## Caveats

- logger doesn't create destination directory, so it should be created manually!
