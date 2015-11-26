/**
 * Created by PCPrincess on 11/25/2015.
 */
var main = function() {
    $('.icon-menu').click(function() {
        $('.dropdown-menu').toggle();
    });
}

/* jquery will 'run' main (above) once the HTML page is fully loaded */
$(document).ready(main);