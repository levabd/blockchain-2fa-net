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
    url: 'http://127.0.0.1:8008/state/cd242e0238a9ab0400c59dcc10e5f314acfd1640f058c82514ec49aa8e6f151cd803c6',
    headers: {'Content-Type': 'application/octet-stream'}
}, (err, response) => {
    if (err) return console.log(err)
console.log('response.body', response.body);
    var dataBase64 = JSON.parse(response.body).data
    console.log('length',messages.User.decode(new Buffer(dataBase64, 'base64')));
});
//.




