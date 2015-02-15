/**

 * Logger extension.
 *
 * Uses the CAPTURE_LOG_LEVEL environmental variable to set the logging level.
 *  0 - all
 *  2 - only error
 *  3 - nothing
 */

var system = require('system');

exports.LOG_INFO = 0;
exports.LOG_WARNING = 1;
exports.LOG_ERROR = 2;

var log = function (message, level) {
  // Get log level.
  var system_level_setting = exports.LOG_ERROR + 1; // No logging by default.
  if (system.env.hasOwnProperty('CAPTURE_LOG_LEVEL')) {
    system_level_setting = parseInt(system.env['CAPTURE_LOG_LEVEL']);
  }

  if (level >= system_level_setting) {
    var log_level_map = [
      'info',
      'warning',
      'error'
    ];

    console.log('>>', (new Date()), log_level_map[level], ':', message);
  }
};

exports.info = function () {
  log([].slice.apply(arguments), exports.LOG_INFO);
}

exports.warning = function () {
  log([].slice.apply(arguments), exports.LOG_WARNING);
}

exports.error = function () {
  log([].slice.apply(arguments), exports.LOG_ERROR);
}
