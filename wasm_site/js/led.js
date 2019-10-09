function DisplayLEDs(colors) {
    $("li").each(function (i, e) {
        console.log(colors[i])
        $(e).css({ "background-color": `rgb(${colors[i].R}, ${colors[i].G}, ${colors[i].B})` });
    });
}