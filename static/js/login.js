$(document).ready(function() {

    if (getCookie('admin') !== 'true') {
        $('#login').modal({
            keyboard: false,
            backdrop: 'static'
        });
    }

    $('form#loginForm').submit(function(e) {
        e.preventDefault();
        $('#error').addClass('hide')
        var username = $('input#username').val();
        var password = $('input#password').val();
        if (username === 'admin' && password === 'admin') {
            $('#login').modal('hide');
            document.cookie = 'admin=true;'
        } else {
            $('#error').removeClass('hide');
        }
    });

    function getCookie(cname) {
        var name = cname + "=";
        var ca = document.cookie.split(';');
        for(var i=0; i<ca.length; i++) {
            var c = ca[i];
            while (c.charAt(0)==' ') c = c.substring(1);
            if (c.indexOf(name) == 0) return c.substring(name.length,c.length);
        }
        return "";
    }

    $('a#cancel').click(function() {
        document.cookie = "admin=; expires=Thu, 01 Jan 1970 00:00:00 UTC";
        window.location.href="home.php"
    });

});
