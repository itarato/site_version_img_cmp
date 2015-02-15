/**
 * @file
 * PhantomJS script to virtually open the page and save a screenshot.
 *
 * Arguments:
 *  1 - filename without extension,
 *  2 - width of the frame,
 *  3 - height of the frame,
 *  4 - page arguments.
 */

'use strict';

var log = require('./log');

// Load system module to handle arguments.
var system = require('system');

// File path to save the screenshot to.
var screenshot_path = system.args[1];

// Virtual browser width.
var width = system.args[2];
var height = system.args[3];

// Configuration for the page. URL + hooks.
var pageConfig = JSON.parse(system.args[4]);

// Load page module for opening the URL.
var page = require('webpage').create();
page.viewportSize = { width: width, height: height };

// Main wrapper callback to execute the screenshots.
var render = function () {
  log.info('Render start.');

  // Open URL.
  page.open(pageConfig.url, function(status) {
    // Save screenshot.
    page.render(screenshot_path);
    // Finish.
    phantom.exit();
  });
}

// Check if there is a plugin argument.
if (pageConfig.hasOwnProperty('pre_hooks')) {
  var callback = render;

  var pre_hook_length = pageConfig.pre_hooks.length;
  for (var idx = 0; idx < pre_hook_length; idx++) {
    var pluginData = pageConfig.pre_hooks[idx];
    var plugin = require('./' + pluginData['plugin']);
    var callback_orig = callback;
    callback = function () {
      // Execute plugin first, then the screenshot.
      plugin.execute(page, pluginData['params'], callback_orig);
    }
  }

  callback();
}
else {
  // No plugin, execute.
  render();
}
