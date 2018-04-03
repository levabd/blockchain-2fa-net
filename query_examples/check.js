const {createContext, CryptoFactory} = require('sawtooth-sdk/signing')

const context = createContext('secp256k1')
const privateKey = context.newRandomPrivateKey()
const signer = new CryptoFactory(context).newSigner(privateKey)
const cbor = require('cbor')
const protobuf = require('protobufjs')

const atob = require('atob');
const btoa = require('btoa');
const {
    InternalError
} = require('sawtooth-sdk/processor/exceptions')

const request = require('request')

console.log('length');

request.get({
    url: 'http://127.0.0.1:8008/state/cd242e027ba028a896bdfe9606e3a80dfc7f9a016b2e8e07db2443e9ba97f01a88770e',
    headers: {'Content-Type': 'application/octet-stream'}
}, (err, response) => {
    if (err) return console.log(err)
console.log('response.body', response.body);
    var dataBase64 = JSON.parse(response.body).data
    console.log('length', cbor.decode(new Buffer(dataBase64, 'base64')));
});
//.




