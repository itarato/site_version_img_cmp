/**
 * @file
 * PhantomJS script to virtually open the page and save a screenshot.
 *
 * Arguments:
 *  1 - filename without extension,
 *  2 - width of the frame,
 *  3 - page arguments.
 *
 * Example:
 *
 *  Makes a screenshot:
 *  $ phantomjs http://localhost/drupal ~/Desktop/sample.png 960
 *
 *  Makes a screenshot with as logged in user:
 *  $ phantomjs http://localhost/drupal ~/Desktop/sample.png 960 drupal.login+http://localhost/drupal/user+admin+monkey
 */

'use strict';

// Load system module to handle arguments.
var system = require('system');

// File path to save the screenshot to.
var screenshot_path = system.args[1];

// Virtual browser width.
var width = system.args[2];

// Configuration for the page. URL + hooks.
var pageConfig = JSON.parse(system.args[3]);

// Load page module for opening the URL.
var page = require('webpage').create();
page.viewportSize = { width: width, height: 800 };

// Main wrapper callback to execute the screenshots.
var render = function () {
  // Open URL.
  page.open(pageConfig.url, function(status) {
    // Save screenshot.
    page.render(screenshot_path);
    // Finish.
    phantom.exit();
  });
}

// Check if there is a plugin argument
var plugin_args = system.args[4];
if (plugin_args) {
  plugin_args_array = plugin_args.split('+');
  var plugin = require('./' + plugin_args_array.shift());

  // Execute plugin first, then the screenshot.
  plugin.execute(page, plugin_args_array, render);
}
else {
  // No plugin, execute.
  render();
}
