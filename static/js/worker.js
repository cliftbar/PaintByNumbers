//https://blog.boot.dev/golang/running-go-in-the-browser-wasm-web-workers/
addEventListener("error", (e) => {
    postMessage({
        message: e
    });
}, false);

addEventListener('message', async (e) => {
    // Load golang wasm exec script
    // importScripts("/static/js/wasm_exec.js")
    importScripts("/static/js/wasm_exec_tiny.js")

    // initialize the Go WASM glue
    const go = new self.Go();

    // e.data contains the code from the main thread
    // const result = await WebAssembly.instantiateStreaming(fetch("/static/wasm/pbn.wasm"), go.importObject);
    const result = await WebAssembly.instantiateStreaming(fetch("/static/wasm/tinypbn.wasm"), go.importObject);

    // hijack the console.log function to capture stdout
    // send each line of output to the main thread
    let oldLog = console.log;
    console.log = (line) => { postMessage({
        message: line
    }); };

    // Start wasm process
    console.log("wasm start")
    go.run(result.instance);
    let retData = {
        "target": e.data.target
    }
    if (e.data.target === "pixelizor") {
        // Run method
        let output = new Uint8Array(e.data.bytes.length)

        let ret = await pixelizor(e.data.height, e.data.width, e.data.heightFactor, e.data.widthFactor, e.data.numColors, 0.01, e.data.bytes, output)

        retData["img"] = output
        retData["colors"] = ret.split(",")
    } else if (e.data.target === "dominantColors") {
        let ret = await dominantColors(e.data.height, e.data.width, e.data.numColors, e.data.deltaThreshold, e.data.bytes)
        retData["colors"] = ret.split(",")
    } else if (e.data.target === "pixelizeFromPalette") {
        let output = new Uint8Array(e.data.bytes.length)
        let paletteString = e.data.palette.join(",")

        await pixelizeFromPalette(e.data.height, e.data.width, e.data.heightFactor, e.data.widthFactor, paletteString, e.data.bytes, output)
        retData["img"] = output
    } else {
        retData["success"] = false
        retData["reason"] = "unknown wasm method"
    }
    console.log = oldLog;

    // tell the main thread we are done
    postMessage({
        "done": true,
        "wasmPayload": retData
    });
    self.close()
}, false);