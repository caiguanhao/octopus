const net = require('net')
const client = new net.Socket()

function call (method, params) {
  let id = Math.floor(Math.random() * 999999) + 100000

  return new Promise((resolve, reject) => {
    let data = ''
    let result = null
    let done = () => {
      clearTimeout(timeout)
      if (result && result.id === id && result.error === null) {
        resolve(result)
      } else {
        reject(result)
      }
    }
    let timeout = setTimeout(() => {
      done()
    }, 2000)

    client.on('data', (res) => {
      data += res
      try {
        result = JSON.parse(data)
        client.end()
      } catch (e) {}
    })

    client.on('end', done)
    client.on('error', done)

    client.connect(12345, '127.0.0.1', function() {
      client.write(JSON.stringify({
        id,
        method,
        params: [ params ]
      }))
    })
  })
}

function init () {
  return call('Octopus.Init', {
    PortNumber: 0,
    BaudRate: 115200,
    ControllerID: 0
  })
}

function poll (withHistory, maxRetries) {
  maxRetries = maxRetries || 0
  return call('Octopus.Poll', {
    Command: withHistory ? 2 : 1,
    Timeout: 10
  }).catch((err) => {
    if (err && (err.error === '100001' || err.error === '100005')) {
      return init().then(() => {
        return poll(withHistory, maxRetries)
      })
    }
    if (err && err.error === '100032' && maxRetries > 0) {
      return poll(withHistory, maxRetries - 1)
    }
    throw err
  })
}

function deduct (cents, serviceInfo, maxRetries) {
  return poll(false, maxRetries).then(res => {
    return call('Octopus.Deduct', {
      Value: cents,
      ServiceInfo: serviceInfo,
      DeferReleaseFlag: 1
    })
  }).then((res) => {
    return poll(false, maxRetries)
  })
}

deduct(1, '0F0C62D1F7', 10).then(console.log, console.log)
