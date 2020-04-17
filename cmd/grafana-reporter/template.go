package main

import "html/template"

var homeTempl = template.Must(template.New("home").Parse(`
<!DOCTYPE HTML>
<html>

<head>
    <script type="text/javascript">
        function WebSocketFunction() {
            if ("WebSocket" in window) {
                console.log("WebSocket is supported by your Browser!");

                // Let us open a web socket
                ping_interval = null, missed_pongs = 0;
                var ws;
                var url = window.location.href;
                var query_string_paramter = url.split("?")[1];
                var dashboardName = new URLSearchParams(query_string_paramter).get('dashboard');
                var scheme = window.location.protocol
                if (scheme === "https:") {
                    ws = new WebSocket("wss://" + window.location.hostname + (location.port ? ":" + location.port : "") + "/reporter/ws" + "?" + query_string_paramter);
                }
                if (scheme === "http:") {
                    ws = new WebSocket("ws://" + window.location.hostname + (location.port ? ":" + location.port : "") + "/reporter/ws" + "?" + query_string_paramter);

                }
                ws.onopen = function() {
                    // ...
                    // other code which has to be executed after the client
                    // connected successfully through the websocket
                    // ...
                    if (ping_interval === null) {
                        missed_pongs = 0;
                        ping_interval = setInterval(function() {
                            try {
                                missed_pongs++;
                                if (missed_pongs >= 20)
                                    throw new Error("Too many missed pongs.");
                                ws.send("ping");
                            } catch (e) {
                                clearInterval(ping_interval);
                                ping_interval = null;
                                console.warn("Closing connection. Reason: " + e.message);
                                ws.close();
                            }
                        }, 30000);
                    }
                };

                ws.onmessage = function(evt) {
                    if (evt.data === "pong") {
                        // reset the counter for missed pongs
                        missed_pongs = 0;
                        return;
                    }
                    ws.binaryType = "blob";
                    let blob = new Blob([evt.data], {
                        type: 'application/pdf'
                    });

                    //  saveByteArray("test-pdf",evt.data)
                    downloadFile(blob, dashboardName + ".pdf");
                };

                ws.onclose = function() {
                    console.log("Connection is closed...");
                };
            } else {
                // The browser doesn't support WebSocket
                console.log("WebSocket NOT supported by your Browser!");
            }
        }

        //  function saveByteArray(reportName, byte) {
        //  var blob = new Blob([byte], {type: "application/pdf"});
        //  const fileURL = window.URL.createObjectURL(blob);
        //  window.open(fileURL,"_self");
        //  };

        const downloadFile = (blob, fileName) => {
            const link = document.createElement('a');
            // create a blobURI pointing to our Blob
            link.href = URL.createObjectURL(blob);
            link.download = fileName;
            // some browser needs the anchor to be in the doc
            document.body.append(link);
            link.click();
            link.remove();
            // in case the Blob uses a lot of memory
            window.addEventListener('focus', e => URL.revokeObjectURL(link.href), {
                once: true
            });
            window.open(URL.createObjectURL(blob), "_self");
        };
    </script>
    <style>
        .loading-container {
            background-color: #fff;
            border-radius: 4%;
            padding: 20px;
            width: 100px;
        }
        
        .loader {
            border: 4px solid #EFEFF2;
            border-radius: 50%;
            border-top: 4px solid #0063FF;
            width: 44px;
            height: 44px;
            -webkit-animation: spin 1s cubic-bezier(0.22, 0.61, 0.36, 1) infinite;
            animation: spin 1s cubic-bezier(0.22, 0.61, 0.36, 1) infinite;
        }
        
        .loading-logo {
            margin-top: -45px;
        }
        
        .loading-logo svg {
            height: 44px;
            width: auto;
        }
        
        @-webkit-keyframes spin {
            0% {
                -webkit-transform: rotate(0deg);
            }
            100% {
                -webkit-transform: rotate(360deg);
            }
        }
        
        @keyframes spin {
            0% {
                transform: rotate(0deg);
            }
            100% {
                transform: rotate(360deg);
            }
        }
        
        .loading-overlay {
            position: absolute;
            top: 50vh;
            left: 0;
            right: 0;
            z-index: 13;
        }
    </style>
</head>

<body onload="WebSocketFunction()">
    <div class="loading-overlay">
        <center>
            <div class="loading-container">
                <div class="loader"></div>
            </div>
            Generating report...
        </center>
    </div>
</body>

</html>
`))
