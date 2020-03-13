$(document).ready(function () {
    $(".heading").each(function (index, el) {
        var c = $("<div class='header'></div>");
        
        // This places the top bar before the tables
        c.insertBefore(el);
        
    });
});
