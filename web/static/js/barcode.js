let zebraApiKey = "5rOu2Qj1nBa2eBW1KWVaj6b9TdwHeaGQ";
let scannedDataMap = new Map();

function quaggaInit(restart = true) {
    let readerDecodeType = document.getElementById("slct_decoderReaderType").value
    let viewportDiv = document.getElementById("div_interactive")
    let inputStreamResolution = parseInt(document.getElementById("slct_inputStreamResolution").value);
    let patchSize = document.getElementById("slct_locatorPatchSize").value;
    let halfSample = document.getElementById("chbx_locatorDoHalfSample").checked
    let singleChannel = document.getElementById("chbx_singleChannel").checked

    let decoderConfig = {};
    if (readerDecodeType === "ean_extended") {
        decoderConfig["readers"] = [{
            format: "ean_reader",
            config: {
                supplements: [
                    'ean_5_reader', 'ean_2_reader'
                ]
            }
        }]
    } else if(readerDecodeType === "try_common_long") {
        decoderConfig["readers"] = ["upc_reader", "ean_reader"]
        decoderConfig["multiple"] = false
    } else if(readerDecodeType === "try_common_short") {
        decoderConfig["readers"] = ["upc_e_reader", "ean_8_reader"]
        decoderConfig["multiple"] = false
    } else {
        readerDecodeType = readerDecodeType + "_reader";
        decoderConfig["readers"] = [readerDecodeType];
    }

    if (restart) {
        Quagga.stop();
    }

    Quagga.init({
        locate: true,
        inputStream: {
            name: "Live",
            type: "LiveStream",
            target: viewportDiv,
            constraints: {
                facingMode: 'environment',
                focusMode: 'continuous', // Enable continuous autofocus
                zoom: 'auto' // Enable autozoom if supported
            },
            size: inputStreamResolution,
            singleChannel: singleChannel
        },
        decoder: decoderConfig,
        locator: {
            patchSize: patchSize,
            halfSample: halfSample
        }
    }, function (err) {
        if (err) {
            console.log(err);
            return
        }
        console.log("Initialization finished. Ready to start");
        Quagga.start();
    });

    Quagga.onProcessed(onProcessed);
    Quagga.onDetected(callbackThrottle(detected, 2000));
}

function onProcessed(data) {
    console.log(data)
}


function quaggaFromFileInit(restart = true) {
    let fileUpload = document.getElementById("inpt_fileUpload")
    let readerDecodeType = document.getElementById("slct_decoderReaderType").value + "_reader"
    let inputStreamResolution = parseInt(document.getElementById("slct_inputStreamResolution").value);


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
            "size": inputStreamResolution
        }
    }

    Quagga.decodeSingle(decodeConfig, function (data) {
        console.log("decoded")
        console.log(data)
        if (data.codeResult.format.includes("upc") || data.codeResult.format.includes("ean")) {
            zebraCodeLookup(data.codeResult.code);

            let viewportDiv = document.getElementById("div_interactive")
            viewportDiv.replaceChildren();

            let img = document.createElement("img")
            img.src = fileUrl;
            viewportDiv.prepend(img)
        }
    });
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
        zebraCodeLookup(data.codeResult.code);
    } else if(data.codeResult.format.includes("ean")) {
        zebraCodeLookup(data.codeResult.code, "ean");
    } else {
        displayCodeOnly(data.codeResult.code);
    }
}
function displayCodeOnly(code){
    let li = document.createElement("li")
    li.appendChild(document.createTextNode(`${code}`));
    document.getElementById("ul_thumbnails").prepend(li);
}

function upcCodeLookup(upcCode) {
    // UPC Item DB
    let lookupCodeUrl = `https://api.upcitemdb.com/prod/trial/lookup?upc=${upcCode}`;

    fetch(lookupCodeUrl, {
        method: "GET"
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

function zebraCodeLookup(code, codeType="upc") {
    let lookupUrl = `https://api.zebra.com/v2/tools/barcode/lookup?upc=${code}`

    fetch(lookupUrl, {
        method: "GET",
        headers: {
            "apikey": zebraApiKey
        }
    }).then(res => {

        let apiRemaining = res.headers.get("x-ratelimit-remaining");
        let apiReset = res.headers.get("x-ratelimit-reset");
        let apiResetDate = new Date(parseInt(apiReset) * 1000).toLocaleString()

        document.getElementById("txt_apiCallsRemaining").innerHTML = apiRemaining
        document.getElementById("txt_apiResetsAt").innerHTML = apiResetDate

        if (res.status !== 200) {
            res.text().then( t => {
                let li = document.createElement("li")
                li.appendChild(document.createTextNode(`${codeType.toUpperCase()} ${code}: ${t}`));
                document.getElementById("ul_thumbnails").prepend(li);
            });
        } else {
            res.json().then(data => {
                console.log(data)
                data.items.forEach((item) => {
                    let li = document.createElement("li")
                    li.appendChild(document.createTextNode(`${codeType.toUpperCase()} ${item[codeType]}: ${item["title"]}`));
                    document.getElementById("ul_thumbnails").prepend(li);
                    scannedDataMap.set(item[codeType], item["title"])
                })
            })
        }
    })
}

/**
 * This function copies all the scanned data from the scannedDataMap Map object into a newline-separated string
 * and then writes the string to the user's clipboard.
 * When the write is successful, an alert popup displays the copied scanned list to the user.
 */
function scannedDataCopy() {
    let scannedText = "";
    for (let value of scannedDataMap.values()) {
        scannedText += value + "\n";
    }
    navigator.clipboard.writeText(scannedText).then(() => alert(`Copied scanned list\n${scannedText}`));
}
