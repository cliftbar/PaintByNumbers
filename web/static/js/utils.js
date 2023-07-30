const reader = new FileReader();
let fileByteArray = [];
let fileByteArrayBuf = new Uint8Array([]);

function readURL(event, target) {
    console.log("hello image load")
    if (event.target.files.length > 0) {
        let src = URL.createObjectURL(event.target.files[0]);
        let preview = document.getElementById(target);
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

let lastCall = 0;
function callbackThrottle(callback, limit) {
    return function() {
        const now = Date.now();
        if (now - lastCall >= limit) {
            lastCall = now;
            return callback.apply(this, arguments);
        } else {
            console.log("throttled " + callback.name)
        }
    };
}
