function DisplayLEDs(colors) {
    $("li").each(function (i, e) {
        $(e).css({ "background-color": `rgb(${colors[i].R}, ${colors[i].G}, ${colors[i].B})` });
    });
}