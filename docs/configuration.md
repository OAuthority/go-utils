# Configuration/Family File
The configuration/family file is the main way that you make your wiki known to Gowikibot, and includes the url of the wiki and the login credentials for that specific wiki. It is a JSON file where the key is a unique identifier for the wiki and the value is the credentials.  

> The configuration file must be named **config.json** and should reside in the root directory of Gowikibot (along with the main.go file.)
### Example
```json
{
    "en-wikipedia": {
        "apiUrl": "https://en.wikipedia.org/w/api.php",
        "username": "Username",
        "password": "********"
    },
    "telepedia-meta": {
        "apiUrl": "https://meta.telepedia.net/api.php",
        "username": "Username2",
        "password": "********"
    },
}
```
The key can be anything you like that uniquely identifies a particular wiki. Set this to something rememberable as you will need to pass it to the login script when you start up the bot. 

It is recommended to use a Bot Password, and Gowikibot will handle everything related to that.

Note that for obvious reasons, **each key must be different**. You may have two different keys, however, for one wiki, where you may want to use two different bot accounts depending on the task. 