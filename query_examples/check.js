const {createContext, CryptoFactory} = require('sawtooth-sdk/signing')

const context = createContext('secp256k1')
const privateKey = context.newRandomPrivateKey()
const signer = new CryptoFactory(context).newSigner(privateKey)
const crypto = require('crypto')

const _hash = (x) => crypto.createHash('sha512').update(x).digest('hex').toLowerCase()
const cbor = require('cbor')
const FAMILY_NAME = 'tfa';
// const FAMILY_NAME = 'kaztel';
const FAMILY_NAMESPACE = _hash(FAMILY_NAME).substring(0, 6)

const atob = require('atob');
const btoa = require('btoa');
const {
    InternalError
} = require('sawtooth-sdk/processor/exceptions')

const request = require('request')
const fs = require('fs')

const protobufLib = require('protocol-buffers')

// pass a proto file as a buffer/string or pass a parsed protobuf-schema object
// const messages = protobufLib(fs.readFileSync('go/src/tfa/service_client/service_client.proto'))
var messages = protobufLib(fs.readFileSync('go/src/tfa/service/service.proto'))

let pn = '77059124921'
// let pn = '77053237001'
// let pn = '77059127941'
const phoneNumberPart = _hash(pn.toString()).slice(-64)

let address = FAMILY_NAMESPACE + phoneNumberPart
// let address = 'cd242e44ef83f7a657e55ca23b438371a5e307ea5756bc2c0c0b572500ad7efec3aef6'

request.get({
    auth: {
        user: 'sawtooth',
        pass: 'z92aGlTdLVYk6mR'
    },                                      // cd242e44ef83f7a657e55ca23b438371a5e307ea5756bc2c0c0b572500ad7efec3aef6
    url: 'http://127.0.0.1:8008/sawtooth/state/'+ address,
    headers: {'Content-Type': 'application/octet-stream'}
}, (err, response) => {
    // console.log('err, response', err, response);
    if (err) return console.log(err)
    // console.log('response.body', response.body);
    const dataBase64 = JSON.parse(response.body).data
    console.log(messages.User.decode(new Buffer(dataBase64, 'base64')));
});
//.




