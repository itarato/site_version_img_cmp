/**
 * @file
 * PhantomJS script to virtually open the page and save a screenshot.
 *
 * Arguments:
 *  1 - URL,
 *  2 - filename without extension.
 */

// Load system module to handle arguments.
var system = require('system');
// Load page module for opening the URL.
var page = require('webpage').create();
page.viewportSize = { width: 960, height: 800 };

// Open URL.
page.open(system.args[1], function() {
  // Save screenshot.
  page.render("./shots/" + system.args[2]);
  // Finish.
  phantom.exit();
});
