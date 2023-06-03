function quaggaInit(restart = true) {
    let readerDecodeType = document.getElementById("slct_decoderReaderType").value + "_reader"
    let viewportDiv = document.getElementById("div_interactive")
    if (restart) {
        Quagga.stop();
    }

    Quagga.init({
        inputStream: {
            name: "Live",
            type: "LiveStream",
            target: viewportDiv
        },
        decoder: {
            readers: [readerDecodeType]
        }
    }, function (err) {
        if (err) {
            console.log(err);
            return
        }
        console.log("Initialization finished. Ready to start");
        Quagga.start();
    });

    Quagga.onDetected(detected);
}

function quaggaFromFileInit(restart = true) {
    let fileUpload = document.getElementById("inpt_fileUpload")
    let readerDecodeType = document.getElementById("slct_decoderReaderType").value + "_reader"
    let viewportDiv = document.getElementById("div_interactive")

    if (!fileUpload.files || !fileUpload.files.length) {
        alert("Must Upload a File!")
    }

    let fileUrl = URL.createObjectURL(fileUpload.files[0])

    if (restart) {
        Quagga.stop();
    }

    let decodeConfig = {
        'src': fileUrl,
        "locate": true,
        "decoder": {
            "readers": [readerDecodeType]
        },
        "inputStream": {
            "size": 800
        }
    }

    // Quagga.init({
    //     inputStream: {
    //         "size": 800,
    //         "type": "ImageStream"
    //     },
    //     decoder: {
    //         readers: [readerDecodeType]
    //     }
    // }, function (err) {
    //     if (err) {
    //         console.log(err);
    //         return
    //     }
    //     console.log("Initialization finished. Ready to start");
    //     // Quagga.start();
    //
    // });

    Quagga.decodeSingle(decodeConfig, function (data) {
        console.log("decoded")
        console.log(data)
        if (data.codeResult.format.includes("upc")) {
            upcCodeLookup(data.codeResult.code);

            let viewportDiv = document.getElementById("div_interactive")
            viewportDiv.replaceChildren();

            let img = document.createElement("img")
            img.src = fileUrl;
            viewportDiv.prepend(img)
        }
    });

    // Quagga.onProcessed(detected);
}

function quaggaStop(clearViewport = true) {
    Quagga.stop();
    if (clearViewport) {
        let viewportDiv = document.getElementById("div_interactive")
        viewportDiv.replaceChildren();
    }
}

function detected(data){
    console.log(data)
    quaggaStop(false);

    let drawingCtx = Quagga.canvas.ctx.overlay
    let drawingCanvas = Quagga.canvas.dom.overlay;

    if (data.boxes) {
        drawingCtx.clearRect(0, 0, parseInt(drawingCanvas.getAttribute("width")), parseInt(drawingCanvas.getAttribute("height")));
        data.boxes.filter(function (box) {
            return box !== data.box;
        }).forEach(function (box) {
            Quagga.ImageDebug.drawPath(box, {x: 0, y: 1}, drawingCtx, {color: "green", lineWidth: 2});
        });
    }

    if (data.box) {
        Quagga.ImageDebug.drawPath(data.box, {x: 0, y: 1}, drawingCtx, {color: "#00F", lineWidth: 2});
    }

    if (data.codeResult && data.codeResult.code) {
        Quagga.ImageDebug.drawPath(data.line, {x: 'x', y: 'y'}, drawingCtx, {color: 'red', lineWidth: 3});
    }

    if (data.codeResult.format.includes("upc")) {
        upcCodeLookup(data.codeResult.code);
    }
}

function upcCodeLookup(upcCode) {
    // UPC Item DB
    let lookupCodeUrl = `https://api.upcitemdb.com/prod/trial/lookup?upc=${upcCode}`;

    fetch(lookupCodeUrl, {
        method: "GET",
    }).then(res => {
        if (res.status !== 200) {
            alert(`Error: ${res}`)
        }
        res.json().then(data => {
            console.log(data)
            data.items.forEach((item) => {
                let li = document.createElement("li")
                li.appendChild(document.createTextNode(`${item["upc"]}: ${item["title"]}`));
                document.getElementById("ul_thumbnails").prepend(li);
            })
        })
    })
}
