$(document).ready(function() {
    $(".statechannel").each(function(index, el) {
        var c = $("<div class='state-channel-widgets'><span class='current'>Function</span><span>State Channel Tool</span></div>");
        var functionButton = c.children()[0];
        var curlButton = c.children()[1];
        curlButton.onclick = function() {
            $(el).children(".highlight-bash")[0].style.display = "block";
            $(el).children(".highlight-java")[0].style.display = "none";
            functionButton.setAttribute("class", "");
            curlButton.setAttribute("class", "current");
        };
        functionButton.onclick = function() {
            $(el).children(".highlight-bash")[0].style.display = "none";
            $(el).children(".highlight-java")[0].style.display = "block";
            curlButton.setAttribute("class", "");
            functionButton.setAttribute("class", "current");
        };

        if ($(el).children(".highlight-bash").length == 0) {
            // No Java for this example.
            curlButton.style.display = "none";
            // functionButton.style.display = "none";
            // In this case, display sausage by default
            $(el).children(".highlight-java")[0].style.display = "block";
        }
        c.insertBefore(el);
    });
});