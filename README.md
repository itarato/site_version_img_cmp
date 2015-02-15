Site version image compare
==========================

A tool that creates image diffs from different versions of a website. It is useful when you would like to see if a CSS change effected (un)intentionally anywhere on the website. The tool created a snaphot of all defined pages when running the script and then it compares agains each previous versions. A diff image is created.

Good example is to use it in Jenkins when building the app.


Plugins
-------

It can handle "virtual" browser script plugins that runs before the actual screenshot. Plugins can be defined in separate files (JavaScript) and expose an executable function to the runner.

There is an example plugin that is able to authenticate to Drupal before the screenshot. The config:

```JSON
"mylist": {
  "url": "http://dev.site/mylist",
  "pre_hooks": [
    {
      "plugin": "drupal.login",
      "params": ["http://dev.site/user/", "myname", "mypassword"]
    }
  ]
}
```

Plugins (hooks) are listed under ```pre_hooks```. They define the file name: ```plugin``` and the arguments for the executable: ```params``` which is a lis of strings. The file is called ```drupal.login.js```. In the file it's expected to have an executable:

```JavaScript
/**
 * @param page PhantomJS page object to act on
 * @param args Arguments from config
 * @param callback Callback to call when all done
 */
exports.execute = function (page, args, callback) {
}
```


Requirements
------------

* Imagemagick (convert)
* PhantomJS (http://phantomjs.org/)


Compile
-------

`./go build main.go`


Setup
-----

* Create your own configuration ```cp default.config.json config.json```
* Define each pages you want to compare with a unique id
* To enable logging set the environment variable: ```export CAPTURE_LOG_LEVEL=0``` where 0=all, 2=only error, 3=nothing


Usage
-----

* Run the script manually (or through a build system or git hook) ```./main```
* Check the result in ```./shots/diff_*.png```
