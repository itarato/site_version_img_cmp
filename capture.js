/**
 * @file
 * PhantomJS script to virtually open the page and save a screenshot.
 *
 * Arguments:
 *  1 - URL,
 *  2 - filename without extension,
 *  3 - width of the frame.
 */

// Load system module to handle arguments.
var system = require('system');

// URL to load.
var url = system.args[1];

// File path to save the screenshot to.
var screenshot_path = system.args[2];

// Virtual browser width.
var width = system.args[3];

var plugin_args = system.args[4];
if (plugin_args) {
  plugin_args_array = plugin_args.split('+');
  var plugin = require('./' + plugin_args_array.shift());
}

// Load page module for opening the URL.
var page = require('webpage').create();
page.viewportSize = { width: width, height: 800 };

var render = function () {
  // Open URL.
  page.open(url, function(status) {
    // Save screenshot.
    page.render(screenshot_path);
    // Finish.
    phantom.exit();
  });
}

// Login.
plugin.execute(page, plugin_args_array, render);
