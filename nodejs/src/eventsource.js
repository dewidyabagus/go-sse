var EventSource = require('eventsource')

const source = new EventSource("http://localhost:3000/sse")
source.onmessage = (event) => {
    console.log("OnMessage Called:")
    console.log(event)
    console.log(JSON.parse(event.data))
}