exports.execute = function (page, login_url, login_name, login_pass, callback) {
  page.open(login_url, function (status) {
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
      phantom.exit();
    });

    // Login.
    var data = "form_id=user_login&op=Log%20in&name=" + login_name + "&pass=" + login_pass + "&form_build_id=" + form_build_id;
    page.open(login_url, 'post', data, function (status) {
      callback(status);
    });
  });
}
