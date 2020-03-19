$(document).ready(function () {
    $(".table").each(function (index, el) {
        var c = $("<div class='table-tool'><span class='current'>Open Trades</span><span>Margin Requirements</span><span>Counterparty Balances</span></div>");
        var tradesButton = c.children()[0];
        var marginButton = c.children()[1];
        var balanceButton = c.children()[2];

        tradesButton.onclick = function () {
            $(el).children(".wy-table-responsive")[0].style.display = "table";
            $(el).children(".wy-table-responsive")[1].style.display = "none";
            $(el).children(".wy-table-responsive")[2].style.display = "none";

            tradesButton.setAttribute("class", "current");
            marginButton.setAttribute("class", "");
            balanceButton.setAttribute("class", "");
        };
        marginButton.onclick = function () {
            $(el).children(".wy-table-responsive")[0].style.display = "none";
            $(el).children(".wy-table-responsive")[1].style.display = "none";
            $(el).children(".wy-table-responsive")[2].style.display = "table";
            marginButton.setAttribute("class", "current");
            tradesButton.setAttribute("class", "");
            balanceButton.setAttribute("class", "");
        };
        balanceButton.onclick = function () {
            $(el).children(".wy-table-responsive")[0].style.display = "none";
            $(el).children(".wy-table-responsive")[1].style.display = "table";
            $(el).children(".wy-table-responsive")[2].style.display = "none";
            marginButton.setAttribute("class", "");
            tradesButton.setAttribute("class", "");
            balanceButton.setAttribute("class", "current");
        };

        // This places the top bar before the tables
        c.insertBefore(el);

        window.onload = function() {
            $(el).children(".wy-table-responsive")[0].style.display = "table";
            // $(el).children(".wy-table-responsive")[1].style.display = "none";
            // $(el).children(".wy-table-responsive")[2].style.display = "none";

            // tradesButton.setAttribute("class", "current");
            // marginButton.setAttribute("class", "");
            // balanceButton.setAttribute("class", "");
        }
    });
});
