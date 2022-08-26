// import wasmWorker from "/static/js/wasm-worker.js";

let fileByteArray = [];
let fileByteArrayBuf = new Uint8Array([]);
const reader = new FileReader();

function readURL(event) {
    console.log("hello image load")
    if (event.target.files.length > 0) {
        let src = URL.createObjectURL(event.target.files[0]);
        let preview = document.getElementById("img_uploaded");
        preview.src = src;
        reader.readAsArrayBuffer(event.target.files[0]);
        reader.onloadend = (evt) => {
            if (evt.target.readyState === FileReader.DONE) {
                const arrayBuffer = evt.target.result;
                    fileByteArrayBuf = new Uint8Array(arrayBuffer);
                for (const a of fileByteArrayBuf) {
                    fileByteArray.push(a);
                }
                // console.log(fileByteArray)
            }
        }
    }
}

function jsPixelizor() {

    let img = new Image();
    img.src = document.getElementById("img_uploaded").src
    let output = new Uint8Array(fileByteArrayBuf.length)
    pixelizor(img.height, img.width, fileByteArrayBuf, output)
    // const content = new Uint8Array(result.instance.exports.memory.buffer.slice(retMap.imgPtr, retMap.imgPtr + retMap.imgPtrLen));
    const content = new Uint8Array(output);

    let imgBlob =  URL.createObjectURL(
        new Blob([content.buffer], { type: "image/png" })
    );

    document.getElementById("img_altered").src = imgBlob;
}

// https://stackoverflow.com/questions/37099465/web-workers-terminating-abruptly
function workerPixelizor() {
    let tableDiv = document.getElementById("div_color_palette")
    tableDiv.innerHTML = ""

    document.getElementById("img_altered").src = "/static/img/loading.gif"

    let heightFactor = parseInt(document.getElementById("inpt_height_factor").value);
    let widthFactor = parseInt(document.getElementById("inpt_width_factor").value);
    let numColors = parseInt(document.getElementById("inpt_num_colors").value);

    let img = new Image();
    img.src = document.getElementById("img_uploaded").src
    let worker = getWorker()

    // keep reference so that we don't garbage collect the worker I guess?
    window.worker = worker.webWorker
    let params = {
        "target": "pixelizor",
        "bytes": fileByteArrayBuf,
        "height": img.height,
        "width": img.width,
        "heightFactor": heightFactor,
        "widthFactor": widthFactor,
        "numColors": numColors
    }
    useWorker(worker, params, (data) => {
        console.log(data)
    }).then((retData) => {
        // Set pixelized img
        document.getElementById("img_altered").src = URL.createObjectURL(
            new Blob([retData.img.buffer], { type: "image/png" })
        );

        setColorPalette(retData.colors);

        console.log("pixelizor post worker done")
    });
    console.log("worker created")
}

function workerDominantColors() {
    let tableDiv = document.getElementById("div_color_palette")
    tableDiv.innerHTML = ""

    let heightFactor = parseInt(document.getElementById("inpt_height_factor").value);
    let widthFactor = parseInt(document.getElementById("inpt_width_factor").value);
    let numColors = parseInt(document.getElementById("inpt_num_colors").value);
    let deltaThreshold = parseInt(document.getElementById("inpt_delta_threshold").value);

    let img = new Image();
    img.src = document.getElementById("img_uploaded").src
    let worker = getWorker()

    // keep reference so that we don't garbage collect the worker I guess?
    window.worker = worker.webWorker

    let params = {
        "target": "dominantColors",
        "bytes": fileByteArrayBuf,
        "height": img.height,
        "width": img.width,
        "heightFactor": heightFactor,
        "widthFactor": widthFactor,
        "numColors": numColors,
        "deltaThreshold": deltaThreshold
    }
    useWorker(worker, params, (data) => {
        console.log(data)
    }).then((retData) => {
        setColorPalette(retData.colors);

        console.log("dominantColors post worker done")
    });
    console.log("worker created")
}

function workerPixelizeFromPalette() {
    document.getElementById("img_altered").src = "/static/img/loading.gif"
    let table = document.getElementById("div_color_palette")
    let colorPalette = []
    let colorInputs = table.querySelectorAll("form input")

    for (const colorInput of colorInputs) {
        colorPalette.push(colorInput.value)
    }

    let heightFactor = parseInt(document.getElementById("inpt_height_factor").value);
    let widthFactor = parseInt(document.getElementById("inpt_width_factor").value);

    let img = new Image();
    img.src = document.getElementById("img_uploaded").src
    let worker = getWorker()

    // keep reference so that we don't garbage collect the worker I guess?
    window.worker = worker.webWorker
    let params = {
        "target": "pixelizeFromPalette",
        "bytes": fileByteArrayBuf,
        "height": img.height,
        "width": img.width,
        "heightFactor": heightFactor,
        "widthFactor": widthFactor,
        "palette": colorPalette
    }
    useWorker(worker, params, (data) => {
        console.log(data)
    }).then((retData) => {
        // Set pixelized img
        document.getElementById("img_altered").src = URL.createObjectURL(
            new Blob([retData.img.buffer], { type: "image/png" })
        );


        console.log("pixelizeFromPalette post worker done")
    });
    console.log("worker created")
}

function getWorker() {
    return {
        // webWorker: new window.Worker(`/${lang}_worker.js`),
        webWorker: new Worker(`/static/js/worker.js`)
    };
}

function useWorker(worker, params, msgCallback) {
    const promise = new Promise((resolve, reject) => {
        worker.webWorker.onmessage = (event) => {
            if (event.data.done) {
                console.log("wasm done")

                resolve(event.data.wasmPayload);
                return;
            }
            if (event.data.error) {
                reject(event.data.error);
                return;
            }
            msgCallback("worker: " + event.data.message);
        };
        console.log("worker message sent")
    });
    worker.webWorker.postMessage(params);
    return promise;
}

function terminateWorker(worker) {
    worker.webWorker.terminate();
}

function setColorPalette(palette){
    let tableDiv = document.getElementById("div_color_palette")
    tableDiv.innerHTML = ""

    let paletteForm = document.createElement("form")

    for (const color of palette) {
        let picker = document.createElement("input")
        picker.id = "palette-" + color
        picker.type = "color"
        picker.value = color

        // div version
        let colorLabel = document.createElement("label");
        colorLabel.innerHTML = color
        colorLabel.htmlFor = picker.id
        paletteForm.appendChild(colorLabel)
        paletteForm.appendChild(picker)
    }

    tableDiv.appendChild(paletteForm)

}