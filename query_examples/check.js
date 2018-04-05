const {createContext, CryptoFactory} = require('sawtooth-sdk/signing')

const context = createContext('secp256k1')
const privateKey = context.newRandomPrivateKey()
const signer = new CryptoFactory(context).newSigner(privateKey)
const cbor = require('cbor')

const atob = require('atob');
const btoa = require('btoa');
const {
    InternalError
} = require('sawtooth-sdk/processor/exceptions')

const request = require('request')

var protobufLib = require('protocol-buffers')

// pass a proto file as a buffer/string or pass a parsed protobuf-schema object
var messages = protobufLib(fs.readFileSync('go/src/tfa/service_client/service_client.proto'))

console.log('length');

request.get({
    url: 'http://127.0.0.1:8008/state/cd242e09ca183fc3f661681f533667712dee333aa01626c58d49dd067d270e00ba4925',
    headers: {'Content-Type': 'application/octet-stream'}
}, (err, response) => {
    if (err) return console.log(err)
    console.log('response.body', response.body);
    var dataBase64 = JSON.parse(response.body).data
    console.log('length',messages.User.decode(new Buffer(dataBase64, 'base64')));
});
//.




