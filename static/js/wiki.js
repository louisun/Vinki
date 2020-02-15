$(document).ready(function() {

    var setDirectoryBtn = function() {
        $("#directory-btn").click((function() {
            var state = $('#directory-container').css('display');
            if (state == 'none') {
                $("#directory-container").show();
                $("body").css("padding-right", "240px");
            } else {
                $("#directory-container").hide();
                $("body").css("padding-right", "0");
            }
        }));
    }
    var setTitle = function() {
        var title = $("#wiki-content").find("h1").text();
        if (title != "") {
            $("#toolbar-title").text(title);
        } else {
            $("#toolbar-title").text("");
        };
    }
    setDirectoryBtn();
    setTitle();
});

// 调整锚点偏移
(function($, window) {
    var adjustAnchor = function() {
        var $anchor = $(':target'),
            fixedElementHeight = 60;
        if ($anchor.length > 0) {
            window.scrollTo(0, $anchor.offset().top - fixedElementHeight);
        }
    };
    $(window).on('hashchange load', function() {
        adjustAnchor();
    });
})(jQuery, window);