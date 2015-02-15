/**
 * Performs a Drupal login.
 */

var log = require('./log');

/**
 * Plugin callback.
 *
 * Loads the required Drupal login page and log in through the login form.
 */
exports.execute = function (page, plugin_args, callback) {
  log.info('Drupal login attempt');

  var login_url = plugin_args[0];
  var login_name = plugin_args[1];
  var login_pass = plugin_args[2];

  log.info(login_url, login_name, login_pass);

  page.open(login_url, function (status) {
    log.info('login status', status);

    // Extract form build id.
    var form_build_id = page.evaluate(function () {
      var input_elements = document.getElementById('user-login').getElementsByTagName('input');
      var input_element_count = input_elements.length;

      for (var i = 0; i < input_element_count; i++) {
        if (input_elements[i].getAttribute('name') == 'form_build_id') {
          return input_elements[i].getAttribute('value');
        }
      }

      // Cannot find form build id.
      log.error('Cannot find form build id');
      phantom.exit();
    });

    // Login.
    var data = "form_id=user_login&op=Log%20in&name=" + login_name + "&pass=" + login_pass + "&form_build_id=" + form_build_id;
    log.info('Post params', data);

    page.open(login_url, 'post', data, function (status) {
      log.info('Login response status', status);
      callback();
    });
  });
};
