      var formUser = document.querySelector('#form-create-user');
      var userName = document.querySelector('#userName');
      var p1 = document.querySelector('#password');
      var p2 = document.querySelector('#password2');
      var btnSubmit = document.querySelector('#btn-create-account');
      var nameErr = document.querySelector('#username-err');
      var pErr = document.querySelector('#password-err');
      //    username must be unique
      userName.addEventListener('input', function(){
          console.log(userName.value);
          var xhr = new XMLHttpRequest();
          xhr.open('POST', '/api/checkUserName');
          xhr.send(userName.value);
          xhr.addEventListener('readystatechange', function(){
              if (xhr.readyState === 4) {
                  var item = xhr.responseText;
                  console.log(item);
                  if (item == 'true') {
                      nameErr.textContent = 'Username taken - Try another name!';
                  } else {
                      nameErr.textContent = '';
                  }
              }
          });
      });
      //    Validate passwords
      //    listen for submit button click
      formUser.addEventListener('submit', function(e){
          var ok = validatePasswords();
          if (!ok) {
              e.preventDefault();
              return;
          }
      });
      function validatePasswords() {
          pErr.textContent = '';
          if (p1.value === '') {
              pErr.textContent = 'Enter a password.';
              return false;
          }
          if (p1.value !== p2.value) {
              pErr.textContent = 'Your passwords did not match. Please re-enter your passwords.';
              p1.value = '';
              p2.value = '';
              return false;
          }
          return true;
      };

    var lsvalue = 'Login';
    (function(){
      $("#defaultRadio").prop("checked", true)
      $("input:radio[name=login-signup]").click(function() {
        lsvalue = $(this).val();
        if(lsvalue == 'Login'){
          $('.loginForm').css('display', 'flex');
          $('.signupForm').css('display', 'none');
        } else {
          $('.loginForm').css('display', 'none');
          $('.signupForm').css('display', 'flex');
        }
      });
      
      $('.accsubmit').click(function(){
        if(lsvalue == 'Login') {
          //login
          var userName = $('#userNamel').val();
          var password = $('#passwordl').val();
          
          $.ajax({
            url: "/api/login",
            type: "POST",
            data: { 'userName': userName, 'password': password },
          });
        } else if(lsvalue == 'Signup') {
          //signup
          var userName = $('#userName').val();
          var email = $('#email').val();
          var pass = $('#password').val();
          $.ajax({
            url: "/api/signup",
            type: "POST",
            data: { 'email': email, 'userName': userName, 'password': pass },
          });
        } else {
          // something wacky happened
          console.log('Error: form not recognized');
        }
      });
    })();
