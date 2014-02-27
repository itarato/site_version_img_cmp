# Site version image compare

A tool that creates image diffs from different versions of a website. It is useful when you would like to see if a CSS change effected (un)intentionally anywhere on the website. The tool created a snaphot of all defined pages when running the script and then it compares agains each previous versions. A diff image is created.

**Drupal specifics:** There is an extra plugin for the PhantomJS script that is able to authenticate to Drupal and make the screenshot as a logged in user.


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
