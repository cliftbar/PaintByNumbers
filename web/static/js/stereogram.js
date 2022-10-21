let depthMapArrayBuf = new Uint8Array([]);
function getDepthMap(){
    document.getElementById("img_depthmap").src = "/static/img/loading.gif"

    let modelSelection = document.getElementById("slct_depthmap_model").value;

    const formData = new FormData();
    const imgInputField = document.getElementById("inpt_img")
    formData.append("fi", imgInputField.files[0])

    let url = new URL("http://localhost:8069/depthmap")
    url.search = new URLSearchParams({name: "fi", model_type: modelSelection}).toString()

    fetch(url, {
        method: "POST",
        body: formData
    }).then(res => {
        console.log("Request complete! response:", res);
        res.blob().then(blb => {
             blb.arrayBuffer().then(buf => {
                 depthMapArrayBuf = new Uint8Array(buf);
            })
            document.getElementById("img_depthmap").src = URL.createObjectURL(blb);
        })
    });
}

function workerStereogram() {
    document.getElementById("img_stereogram").src = "/static/img/loading.gif"

    let heightYFactor = parseInt(document.getElementById("inpt_height_factor").value);
    let widthXFactor = parseInt(document.getElementById("inpt_width_factor").value);
    let patternWidthXFactor = parseInt(document.getElementById("inpt_pattern_factor").value);
    let shiftAmplitude = parseFloat(document.getElementById("inpt_shift_amplitude").value);
    let invert = document.getElementById("inpt_invert_depthmap").checked;

    let img = new Image();
    img.src = document.getElementById("img_depthmap").src
    let worker = getWorker()

    // keep reference so that we don't garbage collect the worker I guess?
    window.worker = worker.webWorker
    let params = {
        "target": "stereogram",
        "bytes": depthMapArrayBuf,
        "heightYFactor": heightYFactor,
        "widthXFactor": widthXFactor,
        "patternWidthXFactor": patternWidthXFactor,
        "shiftAmplitude": shiftAmplitude,
        "invert": invert
    }
    useWorker(worker, params, (data) => {
        console.log(data)
    }).then((retData) => {
        // Set stereogram img
        document.getElementById("img_stereogram").src = URL.createObjectURL(
            new Blob([retData.img.buffer], { type: "image/png" })
        );

        console.log("stereogram post worker done")
    }).catch((err) => {
        console.log("error, terminating worker: " + err)
        document.getElementById("img_stereogram").src = "#"
        terminateWorker(worker);
    });
    console.log("worker created")
}

function citationCopy() {
    let citationText = `@ARTICLE {Ranftl2022,
\tauthor  = "Ren\\'{e} Ranftl and Katrin Lasinger and David Hafner and Konrad Schindler and Vladlen Koltun",
\ttitle   = "Towards Robust Monocular Depth Estimation: Mixing Datasets for Zero-Shot Cross-Dataset Transfer",
\tjournal = "IEEE Transactions on Pattern Analysis and Machine Intelligence",
\tyear    = "2022",
\tvolume  = "44",
\tnumber  = "3"
}
@article{Ranftl2021,
\tauthor    = {Ren\\'{e} Ranftl and Alexey Bochkovskiy and Vladlen Koltun},
\ttitle     = {Vision Transformers for Dense Prediction},
\tjournal   = {ICCV},
\tyear      = {2021},
}
`
    navigator.clipboard.writeText(citationText).then(() => alert("MiDaS citations copied to clipboard"));
}