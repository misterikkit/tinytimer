let last_build_at = "";

function checkReload() {
    $.get("/built_at.txt", function (data) {
        if (data !== last_build_at) {
            location.reload();
        }
    })
}

$.get("/built_at.txt", function (data) {
    last_build_at = data;
    setInterval(checkReload, 100);
})