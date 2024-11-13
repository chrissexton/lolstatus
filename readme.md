# lolstatus

A quick [webhook](https://github.com/adnanh/webhook/) sink for [status.lol](https://status.lol) posts.

## Options

Usage of lolstatus:

* -date string
    * query for status before this date, format: 2006-01-02 15:04:05
* -db string
   	* path to database (default "status.db")
* -id string
   	* query for id
* -path string
   	* path to status output file (default "status.txt")
* -tpl string
   	* path to template file

## Webhook config

My config is similar but not exact to this:

```json
{
    "id": "statuslol",
    "execute-command": "/path/to/lolstatus",
    "command-working-directory": "/path/to/output/dir",
    "pass-arguments-to-command": [
        {
            "source": "string",
            "name": "-tpl=/path/to/templates/status.tpl"
        },
        {
            "source": "string",
            "name": "-path=/path/to/output/index.html"
        },
    ],
    "pass-file-to-command": [{
        "source": "entire-payload",
        "envname": "FNAME"
    }],
    "trigger-rule": {
        "match": {
            "type": "value",
            "value": "lolhahahasureyeahwhatever",
            "parameter": {
                "source": "header",
                "name": "X-Omg-Lol"
            }
        }
    }
}
```

## Is it good?

Yes.
