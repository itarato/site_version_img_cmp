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

// Load page module for opening the URL.
var page = require('webpage').create();
page.viewportSize = { width: width, height: 800 };

// Open URL.
page.open(url, function() {
  // Save screenshot.
  page.render(screenshot_path);
  // Finish.
  phantom.exit();
});
