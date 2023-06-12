function fetchGetImageToBlob(url) {
    let retBlob;
    fetch(url, {
        method: "GET",
    }).then(res => {
        console.log("Request complete! response:", res);
        res.blob().then(blb => {
            retBlob = blb
        })
    });

    return retBlob;
}

function fileToBlob(file) {
    const reader = new FileReader();

    URL.createObjectURL(file);

    file.blob()
}