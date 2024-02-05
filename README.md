# kittiz-notify

## Example env
- NT_SYNOLOGY=token|chatid

usage:
```
GET,POST
https://localhost/synology

query:
    ?message=hello world
json:
    {"message":"hello world"}
xml:
    <?xml version="1.0" encoding="UTF-8" ?>
    <root>
        <message>hello world</message>
    </root>
form/x-www-form-urlencoded
    message:hello world

```