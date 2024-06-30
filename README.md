Given that I have used a password manager for many years (thus each password being unique), yet being in many data breaches.... let's SHA1 hash all my passwords and compare it to the haveibeenpwned database to find out who has been irresponsible with my accounts.

# Approach
The HIBP database is a flat file format of sorted hashes, prepended with a count of how many time the password has been "seen". 

The filesize is around 40GB, so not quick to search through by grep alone. 

This tool performs a binary search for all the hashes and returns a response within milliseconds.

# file formats
## haveibeenpwned (hibp) file format
A sorted list of SHA1 hashes
```
01B307ACBA4F54F55AAFC33BB06BBBF6CA803E9A:100
5BAA61E4C9B93F3F0682250B6CF8331B7EE68FD8:20
64A6DA114D17AE8F167F6BE2C4AEBC9E99F7466C:1
B1B3773A05C0ED0176787A4F1574FF0075F7521E:9001
```

[Further reading](https://haveibeenpwned.com/Passwords)

## bitwarden file format
Full format has been omitted for brevity.
```json
{
  "encrypted": false,
  "items": [
    {
      "name": "example.com",
      "login": {
        "username": null,
        "password": "password1234"
      }
    },
    {
      "name": "example2.com",
      "login": {
        "username": null,
        "password": "password"
      }
    }
  ]
}
```

## bitwarden -> whopwnedme (wpm) format
I don't want to have my passwords laying around in plain text. So I'll converted them all to SHA1, and this program will work with the below format:
```json
{
    "passwords": [
        {
            "name": "example.com",
            "username": "adrian",
            "sha1": "64A6DA114D17AE8F167F6BE2C4AEBC9E99F7466C"
        }
    ]
}
```

## answering the question
Q: How do I know which website leaked my password `666mySuperUniquePassword@!"£$%^&*()_`?

A: By comparing the locally generated SHA1 hash of my source password, and comparing it to haveibeenpwned: it was `example.com` - It's unique, and only `example.com` knew it.