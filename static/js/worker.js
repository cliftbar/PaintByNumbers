//https://blog.boot.dev/golang/running-go-in-the-browser-wasm-web-workers/
addEventListener("error", (e) => {
    postMessage({
        message: e
    });
}, false);

addEventListener('message', async (e) => {
    // importScripts("/static/js/wasm_exec.js")
    importScripts("/static/js/wasm_exec_tiny.js")
    // initialize the Go WASM glue
    const go = new self.Go();

    // e.data contains the code from the main thread
    // const result = await WebAssembly.instantiateStreaming(fetch("/static/wasm/pbn.wasm"), go.importObject);
    const result = await WebAssembly.instantiateStreaming(fetch("/static/wasm/tinypbn.wasm"), go.importObject);

    // hijack the console.log function to capture stdout
    let oldLog = console.log;
    // send each line of output to the main thread
    console.log = (line) => { postMessage({
        message: line
    }); };
    // console.log(e.data)

    // run the code

    go.run(result.instance);

    let output = new Uint8Array(e.data.bytes.length)
    console.log("dominantColor start")
    let ret = await dominantColor(e.data.height, e.data.width, e.data.heightFactor, e.data.widthFactor, e.data.numColors, e.data.bytes, output)
    console.log(ret)
    console.log("dominantColor end")
    // let retMap = result.instance.exports.dominantColor(e.data.height, e.data.width, e.data.bytes)
    // console.log(output)
    // const content = new Uint8Array(result.instance.exports.mem.buffer.slice(retMap.imgPtr, retMap.imgPtr + retMap.imgPtrLen));
    console.log("dida thing")
    console.log = oldLog;

    // tell the main thread we are done
    postMessage({
        "done": true,
        "img": output,
        "colors": ret.split(",")
    });
    self.close()
}, false);