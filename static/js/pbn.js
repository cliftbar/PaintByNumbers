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

function jsDominantColor() {

    let img = new Image();
    img.src = document.getElementById("img_uploaded").src
    let output = new Uint8Array(fileByteArrayBuf.length)
    dominantColor(img.height, img.width, fileByteArrayBuf, output)
    // const content = new Uint8Array(result.instance.exports.memory.buffer.slice(retMap.imgPtr, retMap.imgPtr + retMap.imgPtrLen));
    const content = new Uint8Array(output);

    document.getElementById("img_altered").src = URL.createObjectURL(
        new Blob([content.buffer], { type: "image/png" })
    );
}

// https://stackoverflow.com/questions/37099465/web-workers-terminating-abruptly
function workerJsDominantColor() {
    document.getElementById("img_altered").src = "/static/img/loading.gif"
    let table = document.getElementById("tbl_color_palette")
    table.innerHTML = ""

    let img = new Image();
    img.src = document.getElementById("img_uploaded").src
    let worker = getWorker()
    window.worker = worker.webWorker
    console.log(worker.webWorker)
    useWorker(worker, {"bytes": fileByteArrayBuf, "height": img.height, "width": img.width}, (data) => {
        console.log(data)
    }).then(() => console.log("use worker done"));
    console.log(worker.webWorker)
    console.log("use worker very over")


}

function getWorker() {
    return {
        // webWorker: new window.Worker(`/${lang}_worker.js`),
        webWorker: new Worker(`/static/js/worker.js`)
    };
}

function useWorker(worker, params, callback) {
    const promise = new Promise((resolve, reject) => {
        worker.webWorker.onmessage = (event) => {
            if (event.data.done) {
                console.log("wasm done")
                document.getElementById("img_altered").src = URL.createObjectURL(
                    new Blob([event.data.img.buffer], { type: "image/png" })
                );

                let table = document.getElementById("tbl_color_palette")
                table.innerHTML = ""

                for (const color of event.data.colors) {
                    console.log(color)
                    let rowNode = document.createElement("tr");

                    let pickerCell = document.createElement("td")
                    let picker = document.createElement("input")
                    picker.type = "color"
                    picker.value = color
                    pickerCell.appendChild(picker)
                    rowNode.appendChild(pickerCell)

                    let textCell = document.createElement("td")
                    textCell.textContent = color
                    rowNode.appendChild(textCell)

                    table.appendChild(rowNode)
                }

                resolve();
                console.log("terminate worker")
                // terminateWorker(worker)
                return;
            }
            if (event.data.error) {
                console.log(event.data)
                reject(event.data.error);
                return;
            }
            callback(event.data.message);
        };
        console.log("worker promise over")
    });
    worker.webWorker.postMessage(params);
    return promise;
}

function terminateWorker(worker) {
    worker.webWorker.terminate();
}